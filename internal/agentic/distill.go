//go:build agentic

package agentic

import (
	"fmt"
	"os"
	"strings"
)

const (
	defaultReadLimit   = 12_000
	maxReadLimit       = 32_000
	minUpsertsPerChunk = 4
	maxBridgeChunks    = 2
)

type SourceInfo struct {
	Name string `json:"name"`
	Size int    `json:"sizeBytes"`
}

type DistillCheckpoint struct {
	SourceOffsets map[string]int `json:"sourceOffsets"`
	Batch         int            `json:"batch"`
}

func listSourceInfos(root string) ([]SourceInfo, error) {
	sourcesDir, _, _ := projectDirs(root)
	entries, err := os.ReadDir(sourcesDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []SourceInfo{}, nil
		}
		return nil, err
	}
	out := make([]SourceInfo, 0, len(entries))
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		info, err := e.Info()
		if err != nil {
			continue
		}
		out = append(out, SourceInfo{Name: e.Name(), Size: int(info.Size())})
	}
	return out, nil
}

func ensureDistillCheckpoint(r *RunState) {
	if r.Distill == nil {
		r.Distill = &DistillCheckpoint{SourceOffsets: map[string]int{}}
	}
	if r.Distill.SourceOffsets == nil {
		r.Distill.SourceOffsets = map[string]int{}
	}
}

func recordSourceRead(root, name string, offset, chunkLen int) error {
	r, err := loadRun(root)
	if err != nil {
		return err
	}
	ensureDistillCheckpoint(&r)
	end := offset + chunkLen
	if end > r.Distill.SourceOffsets[name] {
		r.Distill.SourceOffsets[name] = end
	}
	return saveRun(root, r)
}

func distillSourcesRemain(root string) (bool, string) {
	infos, err := listSourceInfos(root)
	if err != nil || len(infos) == 0 {
		return false, ""
	}
	r, err := loadRun(root)
	if err != nil {
		return true, "could not read checkpoint"
	}
	ensureDistillCheckpoint(&r)
	var unread []string
	for _, s := range infos {
		off := r.Distill.SourceOffsets[s.Name]
		if off < s.Size {
			unread = append(unread, fmt.Sprintf("%s (%d/%d bytes)", s.Name, off, s.Size))
		}
	}
	if len(unread) == 0 {
		return false, ""
	}
	return true, strings.Join(unread, "; ")
}

// rejectDistillDone returns a user-facing reason to keep going, or "" if done may be accepted.
func rejectDistillDone(turns, upserts int, summary, instruction string) string {
	if turns < 2 || upserts == 0 {
		return "no tools or upsert_entry yet — read sources in chunks and write entries before done"
	}
	if fullDictionaryJob(instruction) && partialSummary(summary) {
		return "instruction asks for full coverage but summary sounds partial — keep distilling"
	}
	return ""
}

func fullDictionaryJob(instruction string) bool {
	low := strings.ToLower(instruction)
	for _, w := range []string{"every", "all ", "all\n", "entire", "complete", "full ", "whole "} {
		if strings.Contains(low, w) {
			return true
		}
	}
	return strings.Contains(low, "all headwords") || strings.Contains(low, "all words")
}

func partialSummary(summary string) bool {
	low := strings.ToLower(summary)
	if strings.Contains(low, "sample") || strings.Contains(low, "beginning with") || strings.Contains(low, "first ") {
		return true
	}
	return strings.Contains(low, "headword") && (strings.Contains(low, "12 ") || strings.Contains(low, "few "))
}

func buildDistillContinuePrompt(root, instruction string) (string, error) {
	infos, err := listSourceInfos(root)
	if err != nil {
		return "", err
	}
	r, err := loadRun(root)
	if err != nil {
		return "", err
	}
	ensureDistillCheckpoint(&r)
	entries, _ := listEntrySummaries(root)

	var b strings.Builder
	b.WriteString("Continue dictionary distillation (multi-batch job). One batch is NOT the whole job.\n")
	b.WriteString(fmt.Sprintf("Batch %d. Entries on disk: %d.\n", r.Distill.Batch+1, len(entries)))
	b.WriteString("Unread source ranges — use read_source with limit 12000 (max 32000), advance offset each chunk:\n")
	for _, s := range infos {
		off := r.Distill.SourceOffsets[s.Name]
		if off < s.Size {
			b.WriteString(fmt.Sprintf("- %s: next offset %d, %d bytes remaining\n", s.Name, off, s.Size-off))
		}
	}
	b.WriteString(slowestSourceHint(infos, r.Distill.SourceOffsets))
	if block := formatGuidesPromptBlock(root); block != "" {
		b.WriteString(block)
	}
	b.WriteString("\nWorkflow per chunk: list_source_guides + list_sources + list_entries → read_source (or read_source_bridge) → at least ")
	b.WriteString(fmt.Sprint(minUpsertsPerChunk))
	b.WriteString(" upsert_entry calls → next chunk (advance each source at its own pace).\n")
	b.WriteString("If a headword spans a chunk boundary, re-read with overlap or use read_source_bridge to join two slices. Merge draft entries with read_entry + upsert_entry.\n")
	b.WriteString("Only use {\"done\":true} when EVERY source file is fully read (offsets reached file size) AND the user only wanted a slice.\n")
	b.WriteString("Instruction: ")
	b.WriteString(instruction)
	b.WriteString("\nProceed with the next chunk.")
	return b.String(), nil
}

func slowestSourceHint(infos []SourceInfo, offsets map[string]int) string {
	if len(infos) < 2 {
		return ""
	}
	var pick SourceInfo
	ratio := 1.0
	for _, s := range infos {
		if s.Size <= 0 {
			continue
		}
		off := offsets[s.Name]
		r := float64(off) / float64(s.Size)
		if r < ratio {
			ratio = r
			pick = s
		}
	}
	if pick.Name == "" {
		return ""
	}
	return fmt.Sprintf("Multi-source: files advance independently — catch up %s (~%.0f%% read) if others are ahead.\n", pick.Name, ratio*100)
}

const distillToolRules = `
Dictionary distillation rules:
- At job start, per-source format guides were built from the first 10k bytes (subagents). Follow list_source_guides / SOURCE FORMAT GUIDES when detecting headwords.
- Multiple source files: track each file's offset separately; you may read more from one file before catching up another.
- Sources can be huge. NEVER read a whole file at once. read_source limit 12000 (max 32000). Optional overlap (bytes) re-reads tail of previous chunk to finish split headwords.
- read_source_bridge: up to 2 chunks (same or different files) in one call when a sense spans sources or a boundary.
- Draft entries: list_entries and read_entry before upsert_entry — merge senses from both sources into one entry when needed.
- Each source chunk must yield at least 4 headwords via upsert_entry (unless the chunk truly has fewer definitions). Do not read a new chunk until you wrote enough for the last chunk.
- This job spans MANY batches. Do NOT emit {"done":true} on the first reply.
- {"done":true} only when all bytes of all sources are processed OR the user asked for a small explicit subset you finished.`
