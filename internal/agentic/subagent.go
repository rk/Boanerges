//go:build agentic

package agentic

import (
	"context"
	"fmt"
)

// runSubagent is a shorter tool loop for one topical unit (own step log role).
func (e *toolEnv) runSubagent(ctx context.Context, thread int, instruction string) (map[string]any, error) {
	role := fmt.Sprintf("subagent:%d", thread)
	_ = appendRawRunLog(e.root, role, fmt.Sprintf("start\ninstruction: %s", instruction))
	userPrompt := instruction
	const maxTurns = 12
	topicWrites := 0
	for i := 0; i < maxTurns; i++ {
		if err := e.runner.waitIfPaused(ctx); err != nil {
			return nil, err
		}
		reply, err := ollamaChat(ctx, e.model, e.systemPrompt(), userPrompt)
		if err != nil {
			_ = appendRawRunLog(e.root, role, fmt.Sprintf("turn %d ollama error: %v", i+1, err))
			return nil, err
		}
		_ = appendRawRunLog(e.root, role, fmt.Sprintf("turn %d assistant\n%s", i+1, reply))
		_ = appendStep(e.root, "subagent", fmt.Sprintf("#%d %s", thread, truncate(reply, 220)))
		line := extractJSONLine(reply)
		tool, args, done, summary, err := parseToolTurn(reply)
		if err != nil {
			_ = appendRawRunLog(e.root, role, fmt.Sprintf("turn %d parse error: %v", i+1, err))
			userPrompt = "Invalid JSON. Reply with one JSON object only.\n" + userPrompt
			continue
		}
		if done {
			if topicWrites == 0 {
				_ = appendRawRunLog(e.root, role, "rejected done — upsert_topic required")
				userPrompt = "Subagent must call upsert_topic for this topic before {\"done\":true}.\nTask: " + instruction
				continue
			}
			out := map[string]any{"thread": thread, "summary": summary, "topicUpserts": topicWrites}
			_ = appendRawRunLog(e.root, role, fmt.Sprintf("end ok topicUpserts=%d\nsummary: %s", topicWrites, summary))
			return out, nil
		}
		if tool == "" {
			userPrompt = "Call a tool or finish with {\"done\":true,\"summary\":\"...\"}.\n" + userPrompt
			continue
		}
		_ = appendRawRunLog(e.root, role, fmt.Sprintf("turn %d tool-call %s\nargs=%s", i+1, tool, string(args)))
		result, err := e.executeToolCall(ctx, tool, args, line)
		if err != nil {
			_ = appendRawRunLog(e.root, role, fmt.Sprintf("turn %d tool-error %s: %v", i+1, tool, err))
			userPrompt = fmt.Sprintf("Tool %s error: %v\nRetry.", tool, err)
			continue
		}
		if tool == "upsert_topic" {
			topicWrites++
		}
		_ = appendRawRunLog(e.root, role, fmt.Sprintf("turn %d tool-result %s → %s", i+1, tool, truncate(result, 4000)))
		_ = appendStep(e.root, "subagent-tool", fmt.Sprintf("#%d %s → %s", thread, tool, truncate(result, 160)))
		userPrompt = fmt.Sprintf("Tool %s returned:\n%s\nContinue.\nTask: %s", tool, truncate(result, 3000), instruction)
	}
	_ = appendRawRunLog(e.root, role, fmt.Sprintf("end fail max turns (%d) topicUpserts=%d", maxTurns, topicWrites))
	return map[string]any{"thread": thread, "topicUpserts": topicWrites, "error": "max turns exceeded"}, fmt.Errorf("subagent max turns exceeded")
}