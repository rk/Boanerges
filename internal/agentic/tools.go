//go:build agentic

package agentic

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"boanerges/internal/bibleapp"
)

type toolEnv struct {
	root             string
	model            string
	job              string
	instruction      string
	translation      string
	bible            *bibleapp.BibleService
	dict             *bibleapp.DictionaryService
	runner           *runHandle
	nextSubagentThread int
}

func (e *toolEnv) systemPrompt() string {
	return buildSystemPrompt(e.job, e.root)
}

func (e *toolEnv) getChapter(ctx context.Context, translation, bookID string, chapter int) (any, error) {
	if translation == "" {
		translation = e.defaultTranslation()
	}
	return e.bible.Chapter(translation, bookID, chapter)
}

func (e *toolEnv) lookupWebster(ctx context.Context, word string) (any, error) {
	return e.dict.Lookup(word)
}

func (e *toolEnv) listSourcesTool() (any, error) {
	return listSourceInfos(e.root)
}

func (e *toolEnv) readSource(name string, offset, limit, overlap int) (map[string]any, error) {
	if limit <= 0 {
		limit = defaultReadLimit
	}
	if limit > maxReadLimit {
		limit = maxReadLimit
	}
	start := offset
	if overlap > 0 {
		start = offset - overlap
		if start < 0 {
			start = 0
		}
	}
	readLen := (offset + limit) - start
	if readLen <= 0 {
		readLen = limit
	}
	if readLen > maxReadLimit+overlap {
		readLen = maxReadLimit + overlap
	}
	text, err := readSourceChunk(e.root, name, start, readLen)
	if err != nil {
		return nil, err
	}
	_ = recordSourceRead(e.root, name, start, len(text))
	return map[string]any{
		"name":         name,
		"text":         text,
		"requestOffset": offset,
		"startOffset":  start,
		"overlap":      overlap,
		"limit":        limit,
		"chars":        len(text),
		"nextOffset":   start + len(text),
	}, nil
}

type sourceChunkArgs struct {
	Name       string `json:"name"`
	SourceFile string `json:"sourceFile"`
	File       string `json:"file"`
	Filename   string `json:"filename"`
	Offset     int    `json:"offset"`
	Limit      int    `json:"limit"`
	Overlap    int    `json:"overlap"`
}

func (a *sourceChunkArgs) resolveName() (string, error) {
	name := strings.TrimSpace(firstNonEmpty(a.Name, a.SourceFile, a.File, a.Filename))
	if name == "" {
		return "", fmt.Errorf("read_source needs args.name (basename), e.g. chambers-1-of-4-A-D.txt")
	}
	name = strings.TrimPrefix(name, "sources/")
	if i := strings.LastIndex(name, "/"); i >= 0 {
		name = name[i+1:]
	}
	return name, nil
}

func (e *toolEnv) readSourceBridge(chunks []sourceChunkArgs) (any, error) {
	if len(chunks) == 0 || len(chunks) > maxBridgeChunks {
		return nil, fmt.Errorf("provide 1–%d chunks", maxBridgeChunks)
	}
	parts := make([]map[string]any, 0, len(chunks))
	var combined strings.Builder
	for i, c := range chunks {
		name, err := c.resolveName()
		if err != nil {
			return nil, err
		}
		part, err := e.readSource(name, c.Offset, c.Limit, c.Overlap)
		if err != nil {
			return nil, err
		}
		parts = append(parts, part)
		if i > 0 {
			combined.WriteString("\n\n--- bridge ---\n\n")
		}
		combined.WriteString(fmt.Sprintf("--- %s @ %d ---\n", c.Name, c.Offset))
		combined.WriteString(part["text"].(string))
	}
	return map[string]any{"chunks": parts, "combinedText": combined.String()}, nil
}

func (e *toolEnv) listEntriesTool() (any, error) {
	return listEntrySummaries(e.root)
}

func (e *toolEnv) listTopicsTool() (any, error) {
	return listTopicSummaries(e.root)
}

func (e *toolEnv) upsertEntryTool(raw json.RawMessage, toolLine string) (map[string]any, error) {
	entries, err := decodeDictEntries(raw, toolLine)
	if err != nil {
		return nil, err
	}
	_, entriesDir, _ := projectDirs(e.root)
	written := make([]string, 0, len(entries))
	for _, entry := range entries {
		if err := validateDictEntry(entry, e.instruction); err != nil {
			return nil, err
		}
		slug := slugify(entry.Headword)
		if slug == "" {
			return nil, fmt.Errorf("invalid headword %q", entry.Headword)
		}
		path := filepath.Join(entriesDir, slug+".json")
		if err := upsertEntry(e.root, entry); err != nil {
			return nil, err
		}
		if _, err := os.Stat(path); err != nil {
			return nil, fmt.Errorf("write verify %s: %w", path, err)
		}
		written = append(written, slug+".json")
	}
	return map[string]any{"written": written, "count": len(written)}, nil
}

func (e *toolEnv) upsertTopicTool(raw json.RawMessage) error {
	var doc TopicDoc
	if err := json.Unmarshal(raw, &doc); err != nil {
		return err
	}
	return upsertTopic(e.root, doc)
}

// executeToolCall runs one tool invocation from model JSON (simplified contract).
func (e *toolEnv) executeToolCall(ctx context.Context, name string, args json.RawMessage, toolLine string) (string, error) {
	name = normalizeToolName(name)
	if err := toolAllowedOrError(e.job, name); err != nil {
		return "", err
	}
	switch name {
	case "search", "search_bible":
		var a searchArgs
		if err := json.Unmarshal(args, &a); err != nil {
			return "", err
		}
		if err := a.validate(); err != nil {
			return "", err
		}
		out, err := e.keywordSearch(a.Query, a.Translation, a.pageSize(), a.Offset)
		return marshalToolResult(out, err)
	case "read_reference":
		var a readReferenceArgs
		if err := json.Unmarshal(args, &a); err != nil {
			return "", err
		}
		start, end, err := a.verseRange()
		if err != nil {
			return "", err
		}
		text, err := e.readReferencePlain(a.Translation, a.BookID, a.Chapter, start, end)
		if err != nil {
			return "", err
		}
		return text, nil
	case "get_chapter":
		var a struct {
			Translation string `json:"translation"`
			BookID      string `json:"bookId"`
			Chapter     int    `json:"chapter"`
		}
		if err := json.Unmarshal(args, &a); err != nil {
			return "", err
		}
		out, err := e.getChapter(ctx, a.Translation, a.BookID, a.Chapter)
		return marshalToolResult(out, err)
	case "cross_references":
		var a struct {
			BookID  string `json:"bookId"`
			Chapter int    `json:"chapter"`
			Verse   int    `json:"verse"`
		}
		if err := json.Unmarshal(args, &a); err != nil {
			return "", err
		}
		out, err := e.crossRefsArray(a.BookID, a.Chapter, a.Verse)
		return marshalToolResult(out, err)
	case "lookup_webster":
		var a struct {
			Word string `json:"word"`
		}
		if err := json.Unmarshal(args, &a); err != nil {
			return "", err
		}
		out, err := e.lookupWebster(ctx, a.Word)
		return marshalToolResult(out, err)
	case "list_sources":
		out, err := e.listSourcesTool()
		return marshalToolResult(out, err)
	case "list_source_guides":
		out, err := listSourceGuidesTool(e.root)
		return marshalToolResult(out, err)
	case "read_source":
		var a sourceChunkArgs
		if err := json.Unmarshal(args, &a); err != nil {
			return "", err
		}
		name, err := a.resolveName()
		if err != nil {
			return "", err
		}
		out, err := e.readSource(name, a.Offset, a.Limit, a.Overlap)
		result, err := marshalToolResult(out, err)
		if err == nil {
			_ = RunStateChanged.Broadcast(struct{}{})
		}
		return result, err
	case "read_source_bridge":
		var a struct {
			Chunks []sourceChunkArgs `json:"chunks"`
		}
		if err := json.Unmarshal(args, &a); err != nil {
			return "", err
		}
		out, err := e.readSourceBridge(a.Chunks)
		result, err := marshalToolResult(out, err)
		if err == nil {
			_ = RunStateChanged.Broadcast(struct{}{})
		}
		return result, err
	case "read_entry":
		var a struct {
			Headword string `json:"headword"`
		}
		if err := json.Unmarshal(args, &a); err != nil {
			return "", err
		}
		out, err := readEntry(e.root, a.Headword)
		return marshalToolResult(out, err)
	case "list_entries":
		out, err := e.listEntriesTool()
		return marshalToolResult(out, err)
	case "list_topics":
		out, err := e.listTopicsTool()
		return marshalToolResult(out, err)
	case "upsert_entry":
		out, err := e.upsertEntryTool(args, toolLine)
		if err != nil {
			return "", err
		}
		_ = RunStateChanged.Broadcast(struct{}{})
		return marshalToolResult(out, nil)
	case "upsert_topic":
		var doc TopicDoc
		if err := json.Unmarshal(args, &doc); err != nil {
			return "", err
		}
		if err := upsertTopic(e.root, doc); err != nil {
			return "", err
		}
		slug := strings.TrimSpace(doc.Slug)
		if slug == "" {
			slug = slugify(doc.Title)
		}
		_ = RunStateChanged.Broadcast(struct{}{})
		return marshalToolResult(map[string]any{"written": slug + ".json", "title": doc.Title}, nil)
	case "run_subagent":
		var a struct {
			Instruction string `json:"instruction"`
		}
		if err := json.Unmarshal(args, &a); err != nil {
			return "", err
		}
		if a.Instruction == "" {
			return "", fmt.Errorf("instruction required")
		}
		e.nextSubagentThread++
		thread := e.nextSubagentThread
		out, err := e.runSubagent(ctx, thread, a.Instruction)
		if err != nil {
			_ = appendRawRunLog(e.root, fmt.Sprintf("subagent:%d", thread), "error: "+err.Error())
		}
		return marshalToolResult(out, err)
	default:
		return "", fmt.Errorf("unknown tool %q", name)
	}
}

func marshalToolResult(out any, err error) (string, error) {
	if err != nil {
		return "", err
	}
	data, err := json.Marshal(out)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

const dictionaryToolsPrompt = `You are building Boanerges DICTIONARY entries from public-domain sources (sources/ and entries/ only).
Your entire reply MUST be one JSON object (no markdown, no explanation).
Tool call shape: {"tool":"<name>","args":{...}}
Dictionary tools ONLY: list_sources, list_source_guides, read_source, read_source_bridge, read_entry, list_entries, upsert_entry, lookup_webster, run_subagent.
upsert_entry: {"tool":"upsert_entry","args":{"headword":"Aaron","modern":"...","senses":[{"gloss":"...","sourceFile":"file.txt","sourceSpan":"1200-1280"}]}}
read_source: {"tool":"read_source","args":{"name":"basename.txt","offset":0,"limit":8000}} — use args.name, not sourceFile.
Never put job instructions into headword/modern. {"done":false,...} is ignored.
When done: {"done":true,"summary":"..."}
list_sources returns name and sizeBytes. Prose does not write files.`

func parseModelLine(line string) (tool string, args json.RawMessage, done bool, summary string, err error) {
	return parseToolTurn(line)
}
