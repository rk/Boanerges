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
	root   string
	model  string
	bible  *bibleapp.BibleService
	dict   *bibleapp.DictionaryService
	runner *runHandle
}

func (e *toolEnv) searchBible(ctx context.Context, query, translation string, limit int) (any, error) {
	if limit <= 0 || limit > 50 {
		limit = 20
	}
	if translation == "" {
		translation = "asv"
	}
	return e.bible.Search(query, translation, limit, 0)
}

func (e *toolEnv) getChapter(ctx context.Context, translation, bookID string, chapter int) (any, error) {
	if translation == "" {
		translation = "asv"
	}
	return e.bible.Chapter(translation, bookID, chapter)
}

func (e *toolEnv) crossRefs(ctx context.Context, bookID string, chapter, verse int) (any, error) {
	return e.bible.CrossReferences(bookID, chapter, verse)
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
	Name    string `json:"name"`
	Offset  int    `json:"offset"`
	Limit   int    `json:"limit"`
	Overlap int    `json:"overlap"`
}

func (e *toolEnv) readSourceBridge(chunks []sourceChunkArgs) (any, error) {
	if len(chunks) == 0 || len(chunks) > maxBridgeChunks {
		return nil, fmt.Errorf("provide 1–%d chunks", maxBridgeChunks)
	}
	parts := make([]map[string]any, 0, len(chunks))
	var combined strings.Builder
	for i, c := range chunks {
		part, err := e.readSource(c.Name, c.Offset, c.Limit, c.Overlap)
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
	switch name {
	case "search_bible":
		var a struct {
			Query        string `json:"query"`
			Translation  string `json:"translation"`
			Limit        int    `json:"limit"`
		}
		if err := json.Unmarshal(args, &a); err != nil {
			return "", err
		}
		out, err := e.searchBible(ctx, a.Query, a.Translation, a.Limit)
		return marshalToolResult(out, err)
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
		out, err := e.crossRefs(ctx, a.BookID, a.Chapter, a.Verse)
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
		out, err := e.readSource(a.Name, a.Offset, a.Limit, a.Overlap)
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
		if err := e.upsertTopicTool(args); err != nil {
			return "", err
		}
		return "ok", nil
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
		summary, err := e.runSubagent(ctx, a.Instruction)
		return marshalToolResult(map[string]string{"summary": summary}, err)
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

const toolsPrompt = `You are an assistant building Boanerges study data from public-domain sources.
Your entire reply MUST be one JSON object (no markdown, no explanation).
Tool call shape:
{"tool":"<name>","args":{...}}
To write a dictionary entry you MUST call:
{"tool":"upsert_entry","args":{"headword":"...","modern":"...","senses":[{"gloss":"...","sourceFile":"...","sourceSpan":"..."}]}}
Listing headwords in summary or prose does NOT save files — only upsert_entry does.
Tool names: search_bible, get_chapter, cross_references, lookup_webster, list_sources, list_source_guides, read_source, read_source_bridge, read_entry, list_entries, list_topics, upsert_entry, upsert_topic, run_subagent (heavy tasks only: use sparingly; format guides are precomputed).
When done with a unit of work, respond with:
{"done":true,"summary":"..."}
Reuse existing entries/topics when headwords or titles match. Cite source spans for dictionary senses. Group topical refs with labels like "Ex. Of".
list_sources returns name and sizeBytes per file.
upsert_entry example: {"tool":"upsert_entry","args":{"headword":"Able","modern":"...","senses":[{"gloss":"...","sourceFile":"file.txt","sourceSpan":"100-200"}]}}
Fields must be inside args (headword required). Prose without a tool call does not write files.`

func parseModelLine(line string) (tool string, args json.RawMessage, done bool, summary string, err error) {
	return parseToolTurn(line)
}
