//go:build agentic

package agentic

import (
	"encoding/json"
	"fmt"
	"strings"
)

func isTopicalJob(job string) bool {
	switch job {
	case "from_reference", "from_topic":
		return true
	default:
		return false
	}
}

const topicalToolRules = `
Topical index job only (topics/ on disk). This project has NO dictionary sources and NO entries/.
Forbidden tools: list_sources, list_source_guides, read_source, read_source_bridge, read_entry, list_entries, upsert_entry, lookup_webster.
Scripture: read_reference (plain verse text), search (keyword FTS like Study search column — args.query min 2 chars, not bookId/chapter/verse), cross_references (single verse in → array of refs out).
Workflow: read seed verse → list_topics → name every topic in the verse → run_subagent ONCE PER TOPIC TITLE (subagent searches scripture and must upsert_topic) → main job may upsert_topic after subagent.
Never {"done":true} until every identified topic has upsert_topic written this run (topics/*.json).
Group refs with short labels (e.g. "From", "Ex. Of"). Reuse existing topics when titles match.`

const topicalToolsPrompt = `You are building Boanerges TOPICAL scripture index data (not a dictionary).
Your entire reply MUST be one JSON object (no markdown, no explanation).
Tool call shape: {"tool":"<name>","args":{...}}
Allowed tools ONLY: read_reference, search, cross_references, list_topics, upsert_topic, run_subagent.
read_reference: {"tool":"read_reference","args":{"bookId":"jhn","chapter":3,"verse":16}} — optional verseEnd; returns plain text, not JSON.
search: {"tool":"search","args":{"query":"love neighbor","pageSize":20,"offset":0}} — same keyword search as the Study UI (default translation asv).
cross_references: {"tool":"cross_references","args":{"bookId":"jhn","chapter":3,"verse":16}} — returns JSON array of {bookId,chapter,verse}.
run_subagent: {"tool":"run_subagent","args":{"instruction":"Topic: Eternal Life. Seed jhn 3:16. Use search and cross_references; finish with upsert_topic for this topic only."}}
upsert_topic example: {"tool":"upsert_topic","args":{"slug":"eternal-life","title":"Eternal Life","groups":[{"label":"From","refs":[{"bookId":"jhn","chapter":3,"verse":16}]}]}}
When done: {"done":true,"summary":"..."}
Prose without a tool call does not save data.`

func buildTopicalContinuePrompt(root, job, instruction string, params map[string]any, topicUpserts, subagentRuns int) string {
	topics, _ := listTopicSummaries(root)
	var b strings.Builder
	b.WriteString("Continue topical index job (batch). Job: ")
	b.WriteString(job)
	b.WriteString("\nTopics on disk: ")
	b.WriteString(fmt.Sprint(len(topics)))
	b.WriteString(fmt.Sprintf("\nThis run: upsert_topic=%d run_subagent=%d\n", topicUpserts, subagentRuns))
	if topicUpserts == 0 {
		b.WriteString("Still required: run_subagent per topic title, then upsert_topic. {\"done\":true} is rejected until topics/*.json are written.\n")
	}
	if len(params) > 0 {
		data, _ := json.Marshal(params)
		b.WriteString("Seed params: ")
		b.Write(data)
		b.WriteString("\n")
	}
	b.WriteString(topicalToolRules)
	b.WriteString("\nInstruction: ")
	b.WriteString(instruction)
	b.WriteString(seedReferenceHint(job, params))
	b.WriteString("\nNext: one tool JSON (run_subagent, upsert_topic, read_reference, search, cross_references).")
	return b.String()
}
