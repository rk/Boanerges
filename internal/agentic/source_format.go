//go:build agentic

package agentic

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const (
	formatSampleBytes = 10_000
	sourceGuidesName  = "source-guides.json"
)

type SourceFormatGuide struct {
	File       string `json:"file"`
	SizeBytes  int    `json:"sizeBytes"`
	Rules      string `json:"rules"`
	AnalyzedAt string `json:"analyzedAt"`
}

type SourceGuidesDoc struct {
	Guides []SourceFormatGuide `json:"guides"`
}

func sourceGuidesPath(root string) string {
	return filepath.Join(root, sourceGuidesName)
}

func loadSourceGuides(root string) (SourceGuidesDoc, error) {
	var doc SourceGuidesDoc
	path := sourceGuidesPath(root)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return SourceGuidesDoc{Guides: []SourceFormatGuide{}}, nil
	}
	if err := readJSON(path, &doc); err != nil {
		return doc, err
	}
	if doc.Guides == nil {
		doc.Guides = []SourceFormatGuide{}
	}
	return doc, nil
}

func saveSourceGuides(root string, doc SourceGuidesDoc) error {
	return writeJSON(sourceGuidesPath(root), doc)
}

func guideFresh(g *SourceFormatGuide, info SourceInfo) bool {
	return g != nil && g.File == info.Name && g.SizeBytes == info.Size && strings.TrimSpace(g.Rules) != ""
}

func findGuide(doc SourceGuidesDoc, name string) *SourceFormatGuide {
	for i := range doc.Guides {
		if doc.Guides[i].File == name {
			return &doc.Guides[i]
		}
	}
	return nil
}

const formatAnalysisSystem = `You analyze raw text samples from public-domain dictionary sources.
Return ONE JSON object on its own line only:
{"done":true,"summary":"<rules>"}

The summary must be 5–12 short bullet lines (use "• ") describing ONLY how to spot dictionary entries in THIS file:
• headword line pattern (bold, ALL CAPS, numbered, etc.)
• how definitions/senses continue on following lines
• cross-reference or etymology lines to skip or attach
• paragraph/section breaks
• false positives that look like headwords but are not
• suggested sourceSpan citation style (line prefix, offset, etc.)

Do not list actual headwords from the sample. No tools. No prose outside the JSON.`

func ollamaAnalyzeSourceFormat(ctx context.Context, model, fileName, sample string) (string, error) {
	if len(sample) > formatSampleBytes {
		sample = sample[:formatSampleBytes]
	}
	user := fmt.Sprintf("Source file: %s\nFirst ~%d characters:\n---\n%s\n---\nAnalyze entry structure.", fileName, len(sample), sample)
	reply, err := ollamaChatFormatted(ctx, model, formatAnalysisSystem, user, true)
	if err != nil {
		return "", err
	}
	line := extractJSONLine(reply)
	_, _, done, summary, err := parseModelLine(line)
	if err != nil || !done || strings.TrimSpace(summary) == "" {
		return "", fmt.Errorf("format analysis invalid response: %s", truncate(reply, 120))
	}
	return strings.TrimSpace(summary), nil
}

// ensureSourceGuides runs one-shot subagents (parallel, max 2) per stale source file.
func ensureSourceGuides(ctx context.Context, root, model string) error {
	infos, err := listSourceInfos(root)
	if err != nil {
		return err
	}
	if len(infos) == 0 {
		return nil
	}
	doc, err := loadSourceGuides(root)
	if err != nil {
		return err
	}
	var mu sync.Mutex
	var errs []string
	var wg sync.WaitGroup
	sem := make(chan struct{}, 2)

	for _, info := range infos {
		if g := findGuide(doc, info.Name); guideFresh(g, info) {
			continue
		}
		wg.Add(1)
		go func(info SourceInfo) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			sample, err := readSourceChunk(root, info.Name, 0, formatSampleBytes)
			if err != nil {
				mu.Lock()
				errs = append(errs, fmt.Sprintf("%s: %v", info.Name, err))
				mu.Unlock()
				return
			}
			_ = appendStep(root, "subagent", fmt.Sprintf("Format scan %s (%d byte sample)", info.Name, len(sample)))
			rules, err := ollamaAnalyzeSourceFormat(ctx, model, info.Name, sample)
			if err != nil {
				mu.Lock()
				errs = append(errs, fmt.Sprintf("%s: %v", info.Name, err))
				mu.Unlock()
				_ = appendStep(root, "system", "Format scan failed "+info.Name+": "+err.Error())
				return
			}
			guide := SourceFormatGuide{
				File:       info.Name,
				SizeBytes:  info.Size,
				Rules:      rules,
				AnalyzedAt: time.Now().UTC().Format(time.RFC3339),
			}
			mu.Lock()
			updated := false
			for i := range doc.Guides {
				if doc.Guides[i].File == info.Name {
					doc.Guides[i] = guide
					updated = true
					break
				}
			}
			if !updated {
				doc.Guides = append(doc.Guides, guide)
			}
			mu.Unlock()
			_ = appendStep(root, "subagent", "Format rules "+info.Name+": "+truncate(rules, 200))
			_ = appendRawRunLog(root, "format-guide", fmt.Sprintf("file=%s\n%s", info.Name, rules))
		}(info)
	}
	wg.Wait()
	if err := saveSourceGuides(root, doc); err != nil {
		return err
	}
	if len(errs) > 0 {
		return fmt.Errorf("some format scans failed: %s", strings.Join(errs, "; "))
	}
	return nil
}

func formatGuidesPromptBlock(root string) string {
	doc, err := loadSourceGuides(root)
	if err != nil || len(doc.Guides) == 0 {
		return ""
	}
	var b strings.Builder
	b.WriteString("\nSOURCE FORMAT GUIDES (per-file; use when splitting headwords):\n")
	for _, g := range doc.Guides {
		b.WriteString(fmt.Sprintf("\n## %s\n%s\n", g.File, g.Rules))
	}
	return b.String()
}

func listSourceGuidesTool(root string) (any, error) {
	return loadSourceGuides(root)
}
