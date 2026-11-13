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

	env := &toolEnv{root: root, model: model, job: job, instruction: instruction, bible: r.bible, dict: r.dict, runner: r.get(root)}
	if m, err := loadManifest(root); err == nil && strings.TrimSpace(m.Translation) != "" {
		env.translation = m.Translation
	}
	h := env.runner

	if job == "distill" {
		if err := ensureSourceGuides(ctx, root, model); err != nil {
			_ = appendStep(root, "system", "Format guides: "+err.Error())
		} else {
			_ = appendStep(root, "system", "Format guides ready (see source-guides.json)")
		}
		_ = RunStateChanged.Broadcast(struct{}{})
	}
	system := buildSystemPrompt(job, root)

	const turnsPerBatch = 24
	const toolResultMax = 1_200
	const instructionTail = 280
	const maxBatches = 250

	topicUpserts := 0
	subagentRuns := 0
	rejectedDoneStreak := 0
	var lastRejectedSummary string
	seedVerseRead := false

batchLoop:
	for batch := 0; batch < maxBatches; batch++ {
		var userPrompt string
		if batch == 0 {
			userPrompt = buildJobPrompt(job, instruction, params)
			if job == "distill" {
				userPrompt += "\n" + distillToolRules
			} else if isTopicalJob(job) {
				userPrompt += "\n" + topicalToolRules
				userPrompt += "\n" + seedReferenceHint(job, params)
			}
		} else if job == "distill" {
			var err error
			userPrompt, err = buildDistillContinuePrompt(root, instruction)
			if err != nil {
				r.fail(root, err.Error())
				return
			}
		} else if isTopicalJob(job) {
			userPrompt = buildTopicalContinuePrompt(root, job, instruction, params, topicUpserts, subagentRuns)
		} else {
			userPrompt = buildJobPrompt(job, instruction, params)
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
				if isTopicalJob(job) {
					userPrompt = "Reply must be ONE JSON object with a topical tool, e.g. {\"tool\":\"read_reference\",\"args\":{\"bookId\":\"jhn\",\"chapter\":3,\"verse\":16}} or {\"tool\":\"search\",\"args\":{\"query\":\"faith\"}}.\n" + userPrompt
				} else {
					userPrompt = "Reply must be ONE JSON object: {\"tool\":\"upsert_entry\",\"args\":{...}} — prose does not write files.\n" + userPrompt
				}
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
				if isTopicalJob(job) {
					if reason := rejectTopicalDone(job, turn+1, topicUpserts, subagentRuns, summary); reason != "" {
						if summary == lastRejectedSummary {
							rejectedDoneStreak++
						} else {
							rejectedDoneStreak = 1
							lastRejectedSummary = summary
						}
						if rejectedDoneStreak >= 12 {
							r.fail(root, "topical stuck: model repeats done without upsert_topic — pause and retry after restart, or use a stronger model")
							return
						}
						_ = appendStep(root, "system", "Continuing: "+reason)
						_ = appendRawRunLog(root, "system", "Rejected done: "+reason)
						_ = RunStateChanged.Broadcast(struct{}{})
						userPrompt = topicalForceNextPrompt(job, params, topicUpserts, subagentRuns, reason)
						continue
					}
					rejectedDoneStreak = 0
				}
				_ = appendStep(root, "system", "Done: "+summary)
				state, _ := loadRun(root)
				state.Status = RunIdle
				_ = saveRun(root, state)
				_ = RunStateChanged.Broadcast(struct{}{})
				return
			}
			if tool == "" && !done {
				if note := progressNoteJSON(line); note != "" {
					userPrompt = topicalProgressHint(job) + "\n"
					_ = appendStep(root, "system", "Ignored done:false — "+truncate(note, 120))
					continue
				}
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
				if isTopicalJob(job) {
					userPrompt = "JSON must include \"tool\" or {\"done\":true,...}. Save with {\"tool\":\"upsert_topic\",\"args\":{...}}\n" + userPrompt
				} else {
					userPrompt = "JSON must include \"tool\" or {\"done\":true,...}. To save entries use {\"tool\":\"upsert_entry\",\"args\":{...}}\n" + userPrompt
				}
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
			if tool == "run_subagent" {
				_ = appendRawRunLog(root, "system", "main: run_subagent blocks until subagent thread finishes (see [subagent:N] sections below)")
			}
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
			if tool == "upsert_topic" && strings.Contains(result, `"written"`) {
				topicUpserts++
			}
			if tool == "run_subagent" {
				subagentRuns++
				topicUpserts += subagentTopicUpserts(result)
			}
			if tool == "read_reference" && isTopicalJob(job) {
				seedVerseRead = true
			}
			if tool == "read_source" || tool == "read_source_bridge" {
				chunkReads++
				upsertsSinceChunk = 0
			}
			_ = appendStep(root, "tool", fmt.Sprintf("%s → %s", tool, truncate(result, 200)))
			_ = RunStateChanged.Broadcast(struct{}{})
			if isTopicalJob(job) && tool == "read_reference" && seedVerseRead && topicUpserts == 0 && subagentRuns == 0 {
				userPrompt = topicalForceNextPrompt(job, params, topicUpserts, subagentRuns,
					"Seed verse loaded. Do not read_reference again until after run_subagent.")
			} else if isTopicalJob(job) {
				userPrompt = fmt.Sprintf("Tool %s returned:\n%s\nNext: one tool JSON only. Job instruction: %s",
					tool, truncate(result, toolResultMax), truncate(instruction, instructionTail))
			} else {
				userPrompt = fmt.Sprintf("Tool %s returned:\n%s\nContinue. Job instruction: %s",
					tool, truncate(result, toolResultMax), truncate(instruction, instructionTail))
			}
		}

		if job == "distill" || isTopicalJob(job) {
			if job == "distill" {
				state, _ := loadRun(root)
				ensureDistillCheckpoint(&state)
				state.Distill.Batch++
				_ = saveRun(root, state)
			}
			_ = appendStep(root, "system", "Batch "+strconv.Itoa(batch+1)+": hit turn limit, starting next batch")
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

func topicalProgressHint(job string) string {
	if isTopicalJob(job) {
		return "Status noted. Next: ONE tool JSON. Example: {\"tool\":\"read_reference\",\"args\":{\"bookId\":\"jhn\",\"chapter\":3,\"verse\":16}}"
	}
	return "Status noted. Next: ONE tool JSON only. If reading: {\"tool\":\"read_source\",\"args\":{\"name\":\"<basename>\",\"offset\":0,\"limit\":8000}}"
}

func progressNoteJSON(line string) string {
	line = strings.TrimSpace(line)
	var probe struct {
		Done    *bool  `json:"done"`
		Summary string `json:"summary"`
		Tool    string `json:"tool"`
	}
	if err := json.Unmarshal([]byte(line), &probe); err != nil {
		return ""
	}
	if probe.Tool != "" || probe.Done == nil || *probe.Done {
		return ""
	}
	return probe.Summary
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "…"
}
