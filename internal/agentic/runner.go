//go:build agentic

package agentic

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"boanerges/internal/bibleapp"
)

type runHandle struct {
	mu     sync.Mutex
	root   string
	cancel context.CancelFunc
}

type runnerRegistry struct {
	mu     sync.Mutex
	runs   map[string]*runHandle
	active int // ponytail: global cap 2 concurrent Ollama orchestrators
	bible  *bibleapp.BibleService
	dict   *bibleapp.DictionaryService
}

func newRunnerRegistry(bible *bibleapp.BibleService, dict *bibleapp.DictionaryService) *runnerRegistry {
	return &runnerRegistry{
		runs:  map[string]*runHandle{},
		bible: bible,
		dict:  dict,
	}
}

func (r *runnerRegistry) get(root string) *runHandle {
	r.mu.Lock()
	defer r.mu.Unlock()
	if h, ok := r.runs[root]; ok {
		return h
	}
	h := &runHandle{root: root}
	r.runs[root] = h
	return h
}

func (r *runnerRegistry) start(root, job, instruction, model string, params map[string]any) error {
	state, err := loadRun(root)
	if err != nil {
		return err
	}
	if state.Status == RunRunning {
		return fmt.Errorf("run already in progress")
	}
	r.mu.Lock()
	if r.active >= 2 {
		r.mu.Unlock()
		return fmt.Errorf("concurrency limit: at most 2 runs at once")
	}
	r.active++
	r.mu.Unlock()
	state.Status = RunRunning
	state.Job = job
	state.Instruction = instruction
	state.Params = params
	state.Error = ""
	if job == "distill" {
		ensureDistillCheckpoint(&state)
	}
	if err := saveRun(root, state); err != nil {
		return err
	}
	_ = appendRawRunLog(root, "system", fmt.Sprintf("Started job %s\nInstruction: %s", job, instruction))
	_ = appendStep(root, "system", fmt.Sprintf("Started job %s", job))

	ctx, cancel := context.WithCancel(context.Background())
	h := r.get(root)
	h.mu.Lock()
	h.cancel = cancel
	h.mu.Unlock()

	go func() {
		defer func() {
			r.mu.Lock()
			r.active--
			r.mu.Unlock()
		}()
		r.loop(ctx, root, job, instruction, model, params)
	}()
	return nil
}

func (r *runnerRegistry) pause(root string) error {
	state, err := loadRun(root)
	if err != nil {
		return err
	}
	if state.Status != RunRunning {
		return nil
	}
	state.Status = RunPaused
	_ = appendStep(root, "system", "Paused")
	if err := saveRun(root, state); err != nil {
		return err
	}
	_ = RunStateChanged.Broadcast(struct{}{})
	return nil
}

func (r *runnerRegistry) resume(root string) error {
	state, err := loadRun(root)
	if err != nil {
		return err
	}
	if state.Status != RunPaused {
		return nil
	}
	state.Status = RunRunning
	_ = appendStep(root, "system", "Resumed")
	if err := saveRun(root, state); err != nil {
		return err
	}
	_ = RunStateChanged.Broadcast(struct{}{})
	return nil
}

func (r *runnerRegistry) loop(ctx context.Context, root, job, instruction, model string, params map[string]any) {
	defer func() {
		state, _ := loadRun(root)
		if state.Status == RunRunning {
			state.Status = RunIdle
			_ = saveRun(root, state)
		}
	}()

	env := &toolEnv{root: root, model: model, bible: r.bible, dict: r.dict, runner: r.get(root)}
	h := env.runner

	system := toolsPrompt
	if job == "distill" {
		if err := ensureSourceGuides(ctx, root, model); err != nil {
			_ = appendStep(root, "system", "Format guides: "+err.Error())
		} else {
			_ = appendStep(root, "system", "Format guides ready (see source-guides.json)")
		}
		_ = RunStateChanged.Broadcast(struct{}{})
		system += distillToolRules
		system += formatGuidesPromptBlock(root)
	}

	const turnsPerBatch = 40
	const maxBatches = 250

batchLoop:
	for batch := 0; batch < maxBatches; batch++ {
		var userPrompt string
		if batch == 0 {
			userPrompt = buildJobPrompt(job, instruction, params)
			if job == "distill" {
				userPrompt += "\n" + distillToolRules
			}
		} else {
			var err error
			userPrompt, err = buildDistillContinuePrompt(root, instruction)
			if err != nil {
				r.fail(root, err.Error())
				return
			}
		}

		upserts := 0
		chunkReads := 0
		upsertsSinceChunk := 0
		for turn := 0; turn < turnsPerBatch; turn++ {
			if err := h.waitIfPaused(ctx); err != nil {
				r.fail(root, err.Error())
				return
			}

			reply, err := ollamaChat(ctx, model, system, userPrompt)
			if err != nil {
				r.fail(root, err.Error())
				return
			}
			_ = appendRawRunLog(root, "assistant", reply)
			_ = appendStep(root, "assistant", truncate(reply, 240))
			line := extractJSONLine(reply)
			tool, args, done, summary, err := parseToolTurn(reply)
			if err != nil {
				if job == "distill" {
					if n, result, serr := salvageUpsertsFromReply(env, reply); serr == nil && n > 0 {
						_ = appendStep(root, "system", fmt.Sprintf("Salvaged %d entries from malformed reply", n))
						_ = appendRawRunLog(root, "salvage", result)
						upserts += n
						upsertsSinceChunk += n
						_ = appendStep(root, "tool", "salvage → "+truncate(result, 200))
						_ = RunStateChanged.Broadcast(struct{}{})
						userPrompt = fmt.Sprintf("Salvaged %d entries. Continue with strict JSON tool calls only.\nOriginal: %s", n, instruction)
						continue
					}
				}
				userPrompt = "Reply must be ONE JSON object: {\"tool\":\"upsert_entry\",\"args\":{...}} — prose does not write files.\n" + userPrompt
				continue
			}
			if done {
				if job == "distill" {
					if reason := rejectDistillDone(turn+1, upserts, summary, instruction); reason != "" {
						_ = appendStep(root, "system", "Continuing: "+reason)
						_ = RunStateChanged.Broadcast(struct{}{})
						continue batchLoop
					}
					if more, detail := distillSourcesRemain(root); more {
						state, _ := loadRun(root)
						ensureDistillCheckpoint(&state)
						state.Distill.Batch++
						_ = saveRun(root, state)
						_ = appendStep(root, "system", fmt.Sprintf("Batch %d: unread sources — %s", state.Distill.Batch, detail))
						_ = RunStateChanged.Broadcast(struct{}{})
						continue batchLoop
					}
				}
				_ = appendStep(root, "system", "Done: "+summary)
				state, _ := loadRun(root)
				state.Status = RunIdle
				_ = saveRun(root, state)
				_ = RunStateChanged.Broadcast(struct{}{})
				return
			}
			if tool == "" {
				if job == "distill" {
					if n, result, serr := salvageUpsertsFromReply(env, reply); serr == nil && n > 0 {
						upserts += n
						upsertsSinceChunk += n
						_ = appendStep(root, "tool", "salvage → "+truncate(result, 200))
						_ = RunStateChanged.Broadcast(struct{}{})
						userPrompt = fmt.Sprintf("Salvaged %d entries. Next reply must be {\"tool\":\"upsert_entry\",\"args\":{...}} JSON only.\n", n)
						continue
					}
				}
				userPrompt = "JSON must include \"tool\" or {\"done\":true,...}. To save entries use {\"tool\":\"upsert_entry\",\"args\":{...}}\n" + userPrompt
				continue
			}
			if job == "distill" && (tool == "read_source" || tool == "read_source_bridge") {
				if chunkReads > 0 && upsertsSinceChunk < minUpsertsPerChunk {
					userPrompt = fmt.Sprintf(
						"Need at least %d upsert_entry calls for the previous chunk (have %d). Use list_entries, read_entry to merge drafts, then upsert before reading more source.",
						minUpsertsPerChunk, upsertsSinceChunk,
					)
					continue
				}
			}
			_ = appendRawRunLog(root, "tool-call", fmt.Sprintf("tool=%s\nline=%s\nargs=%s", tool, line, string(args)))
			result, err := env.executeToolCall(ctx, tool, args, line)
			if err != nil {
				_ = appendRawRunLog(root, "tool-error", err.Error())
				userPrompt = fmt.Sprintf("Tool %s error: %v\nRetry.", tool, err)
				continue
			}
			_ = appendRawRunLog(root, "tool-result", fmt.Sprintf("%s → %s", tool, result))
			if tool == "upsert_entry" && strings.Contains(result, `"written"`) {
				upserts++
				upsertsSinceChunk++
			}
			if tool == "read_source" || tool == "read_source_bridge" {
				chunkReads++
				upsertsSinceChunk = 0
			}
			_ = appendStep(root, "tool", fmt.Sprintf("%s → %s", tool, truncate(result, 200)))
			_ = RunStateChanged.Broadcast(struct{}{})
			userPrompt = fmt.Sprintf("Tool %s returned:\n%s\nContinue the job.\nOriginal: %s", tool, truncate(result, 4000), instruction)
		}

		if job == "distill" {
			state, _ := loadRun(root)
			ensureDistillCheckpoint(&state)
			state.Distill.Batch++
			_ = saveRun(root, state)
			_ = appendStep(root, "system", "Batch "+strconv.Itoa(state.Distill.Batch)+": hit turn limit, starting next batch")
			_ = RunStateChanged.Broadcast(struct{}{})
			continue batchLoop
		}
		r.fail(root, "max turns exceeded")
		return
	}
	r.fail(root, "max batches exceeded")
}

func (r *runnerRegistry) fail(root, msg string) {
	state, _ := loadRun(root)
	state.Status = RunFailed
	state.Error = msg
	_ = saveRun(root, state)
	_ = appendStep(root, "system", "Failed: "+msg)
	_ = RunStateChanged.Broadcast(struct{}{})
}

func (h *runHandle) waitIfPaused(ctx context.Context) error {
	for {
		state, err := loadRun(h.root)
		if err != nil {
			return err
		}
		if state.Status != RunPaused {
			return nil
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(400 * time.Millisecond):
		}
	}
}

func buildJobPrompt(job, instruction string, params map[string]any) string {
	var b strings.Builder
	b.WriteString("Job: ")
	b.WriteString(job)
	b.WriteString("\nInstruction: ")
	b.WriteString(instruction)
	if len(params) > 0 {
		data, _ := json.Marshal(params)
		b.WriteString("\nParams: ")
		b.Write(data)
	}
	b.WriteString("\nBegin.")
	return b.String()
}

func extractJSONLine(reply string) string {
	for _, line := range strings.Split(reply, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "{") {
			return line
		}
	}
	return strings.TrimSpace(reply)
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "…"
}
