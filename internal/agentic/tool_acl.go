//go:build agentic

package agentic

import "fmt"

var topicalTools = map[string]bool{
	"search":           true,
	"search_bible":     true, // alias
	"read_reference":   true,
	"cross_references": true,
	"list_topics":      true,
	"upsert_topic":     true,
	"run_subagent":     true,
}

var dictionaryTools = map[string]bool{
	"list_sources":         true,
	"list_source_guides":   true,
	"read_source":          true,
	"read_source_bridge":   true,
	"read_entry":           true,
	"list_entries":         true,
	"upsert_entry":         true,
	"lookup_webster":       true,
	"run_subagent":         true,
}

func toolAllowed(job, tool string) bool {
	tool = normalizeToolName(tool)
	if isTopicalJob(job) {
		return topicalTools[tool]
	}
	if job == "distill" {
		return dictionaryTools[tool]
	}
	// ponytail: unknown job — allow bible + both stores
	return topicalTools[tool] || dictionaryTools[tool]
}

func toolAllowedOrError(job, tool string) error {
	if !toolAllowed(job, tool) {
		if isTopicalJob(job) {
			return fmt.Errorf("tool %q is for dictionary jobs only; topical may use read_reference, search, cross_references, list_topics, upsert_topic, run_subagent", tool)
		}
		return fmt.Errorf("tool %q is not allowed for job %s", tool, job)
	}
	return nil
}

func buildSystemPrompt(job, root string) string {
	if isTopicalJob(job) {
		return topicalToolsPrompt + topicalToolRules
	}
	if job == "distill" {
		return dictionaryToolsPrompt + distillToolRules + formatGuidesPromptBlock(root)
	}
	return dictionaryToolsPrompt + topicalToolsPrompt
}
