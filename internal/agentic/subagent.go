//go:build agentic

package agentic

import (
	"context"
	"fmt"
)

// runSubagent is a shorter tool loop for one topical unit (own step log role).
func (e *toolEnv) runSubagent(ctx context.Context, instruction string) (string, error) {
	userPrompt := instruction
	const maxTurns = 12
	for i := 0; i < maxTurns; i++ {
		if err := e.runner.waitIfPaused(ctx); err != nil {
			return "", err
		}
		reply, err := ollamaChat(ctx, e.model, toolsPrompt, userPrompt)
		if err != nil {
			return "", err
		}
		_ = appendStep(e.root, "subagent", truncate(reply, 240))
		line := extractJSONLine(reply)
		tool, args, done, summary, err := parseToolTurn(reply)
		if err != nil {
			userPrompt = "Invalid JSON. Reply with one JSON object only.\n" + userPrompt
			continue
		}
		if done {
			return summary, nil
		}
		if tool == "" {
			userPrompt = "Call a tool or finish with {\"done\":true,\"summary\":\"...\"}.\n" + userPrompt
			continue
		}
		result, err := e.executeToolCall(ctx, tool, args, line)
		if err != nil {
			userPrompt = fmt.Sprintf("Tool %s error: %v\nRetry.", tool, err)
			continue
		}
		_ = appendStep(e.root, "subagent-tool", fmt.Sprintf("%s → %s", tool, truncate(result, 160)))
		userPrompt = fmt.Sprintf("Tool %s returned:\n%s\nContinue.\nTask: %s", tool, truncate(result, 3000), instruction)
	}
	return "", fmt.Errorf("subagent max turns exceeded")
}