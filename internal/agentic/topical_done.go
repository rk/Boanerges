//go:build agentic

package agentic

import (
	"encoding/json"
	"fmt"
	"strings"
)

func subagentTopicUpserts(result string) int {
	var probe struct {
		TopicUpserts int `json:"topicUpserts"`
	}
	if err := json.Unmarshal([]byte(result), &probe); err != nil {
		return 0
	}
	return probe.TopicUpserts
}

// rejectTopicalDone returns a reason to keep going, or "" if done may be accepted.
func rejectTopicalDone(job string, turns, topicUpserts, subagentRuns int, summary string) string {
	if topicUpserts == 0 && summaryClaimsTopicWrite(summary) {
		return "summary claims upsert/save but upsert_topic never ran — call upsert_topic JSON or run_subagent that upserts; prose does not write files"
	}
	if topicUpserts == 0 {
		return "no upsert_topic yet — run_subagent to research each topic, then upsert_topic writes topics/*.json; done forbidden until then"
	}
	if job == "from_reference" && subagentRuns == 0 {
		return "from_reference: dispatch run_subagent once per topic title (search/cross_references/read_reference inside), then upsert_topic for each"
	}
	if job == "from_reference" && subagentRuns > 0 && topicUpserts < subagentRuns {
		return "each run_subagent must upsert_topic once — upsert count still below subagent count"
	}
	if topicalPrematureSummary(summary) {
		return "planning or identifying topics is not enough — finish run_subagent + upsert_topic for every topic from the seed verse"
	}
	if turns < 3 && topicUpserts == 0 {
		return "too early — read seed verse, research topics, upsert_topic before done"
	}
	return ""
}

func topicalPrematureSummary(summary string) bool {
	low := strings.ToLower(summary)
	markers := []string{
		"identified", "no new topics", "reused existing", "no topics were needed",
		"already covers", "grouped under", "conceptually",
	}
	hits := 0
	for _, m := range markers {
		if strings.Contains(low, m) {
			hits++
		}
	}
	if hits >= 2 {
		return true
	}
	if strings.Contains(low, "identified") && !summaryClaimsTopicWrite(summary) {
		return true
	}
	return false
}

func summaryClaimsTopicWrite(summary string) bool {
	low := strings.ToLower(summary)
	for _, m := range []string{"upsert", "wrote topic", "saved topic", "written topic"} {
		if strings.Contains(low, m) {
			return true
		}
	}
	return false
}

// topicalForceNextPrompt tells a stuck model the single valid next JSON (short; no huge rule replay).
func topicalForceNextPrompt(job string, params map[string]any, topicUpserts, subagentRuns int, reason string) string {
	book := paramString(params, "bookId")
	ch := paramInt(params, "chapter")
	v := paramInt(params, "verse")
	var b strings.Builder
	b.WriteString(reason)
	b.WriteString("\n{\"done\":true} is rejected. ")
	if topicUpserts == 0 && subagentRuns == 0 {
		b.WriteString("Do not call read_reference again. Reply with ONLY one tool JSON:\n")
		if job == "from_reference" && book != "" && ch > 0 && v > 0 {
			b.WriteString(fmt.Sprintf(
				`{"tool":"run_subagent","args":{"instruction":"Topic: Eternal Life. Seed %s %d:%d. Use search and cross_references; finish with upsert_topic (groups like From with the seed ref)."}}`+"\n",
				book, ch, v,
			))
			b.WriteString(`Pick topic titles from the seed verse (not from this prompt). run_subagent once per topic.`)
		} else {
			b.WriteString(`{"tool":"run_subagent","args":{"instruction":"Topic: <title>. Use search; finish with upsert_topic."}}`)
		}
		return b.String()
	}
	if topicUpserts < subagentRuns {
		b.WriteString("Reply with upsert_topic for the topic your last subagent researched.")
		return b.String()
	}
	b.WriteString("Reply with run_subagent for the next topic, or upsert_topic if merging results.")
	return b.String()
}
