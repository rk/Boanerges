//go:build agentic

package agentic

import (
	"encoding/json"
	"fmt"
	"strings"
)

var toolAliases = map[string]string{
	"upsert":       "upsert_entry",
	"write_entry":  "upsert_entry",
	"add_entry":    "upsert_entry",
	"save_entry":   "upsert_entry",
	"read":         "read_source",
	"list_sources": "list_sources",
}

func normalizeToolName(name string) string {
	name = strings.TrimSpace(name)
	if a, ok := toolAliases[strings.ToLower(name)]; ok {
		return a
	}
	return name
}

// parseToolTurn extracts a tool call or done marker from a model reply (markdown, aliases, multiple objects).
func parseToolTurn(reply string) (tool string, args json.RawMessage, done bool, summary string, err error) {
	objects := extractJSONObjectStrings(reply)
	if len(objects) == 0 {
		return "", nil, false, "", fmt.Errorf("expected JSON")
	}
	var firstErr error
	for _, obj := range objects {
		t, a, d, s, e := parseModelObject(obj)
		if e != nil {
			if firstErr == nil {
				firstErr = e
			}
			continue
		}
		if d {
			return "", nil, true, s, nil
		}
		if t != "" {
			return t, a, false, "", nil
		}
	}
	if implicit, e := parseImplicitUpsert(objects); e == nil && implicit != "" {
		return "upsert_entry", json.RawMessage(implicit), false, "", nil
	}
	if firstErr != nil {
		return "", nil, false, "", firstErr
	}
	return "", nil, false, "", fmt.Errorf("no tool in reply")
}

func parseModelObject(line string) (tool string, args json.RawMessage, done bool, summary string, err error) {
	line = strings.TrimSpace(line)
	if !strings.HasPrefix(line, "{") {
		return "", nil, false, "", fmt.Errorf("expected JSON object")
	}
	var probe struct {
		Tool      string          `json:"tool"`
		Name      string          `json:"name"`
		Action    string          `json:"action"`
		Function  string          `json:"function"`
		Args      json.RawMessage `json:"args"`
		Arguments json.RawMessage `json:"arguments"`
		Done      bool            `json:"done"`
		Summary   string          `json:"summary"`
		Type      string          `json:"type"`
	}
	if err := json.Unmarshal([]byte(line), &probe); err != nil {
		return "", nil, false, "", err
	}
	if probe.Done || strings.EqualFold(probe.Type, "done") {
		return "", nil, true, probe.Summary, nil
	}
	tool = normalizeToolName(firstNonEmpty(probe.Tool, probe.Name, probe.Action, probe.Function))
	args = probe.Args
	if len(args) == 0 || string(args) == "null" {
		args = probe.Arguments
	}
	if tool != "" {
		args = coalesceToolArgs(line, tool, args)
		return tool, args, false, "", nil
	}
	return "", nil, false, "", nil
}

func coalesceToolArgs(line, tool string, args json.RawMessage) json.RawMessage {
	if len(args) > 0 && string(args) != "null" && string(args) != "{}" {
		return args
	}
	if tool != "upsert_entry" && tool != "upsert_topic" {
		return args
	}
	var m map[string]json.RawMessage
	if err := json.Unmarshal([]byte(line), &m); err != nil {
		return args
	}
	delete(m, "tool")
	delete(m, "name")
	delete(m, "action")
	delete(m, "function")
	delete(m, "args")
	delete(m, "arguments")
	delete(m, "done")
	delete(m, "summary")
	delete(m, "type")
	if len(m) == 0 {
		return args
	}
	data, err := json.Marshal(m)
	if err != nil {
		return args
	}
	return data
}

func parseImplicitUpsert(objects []string) (string, error) {
	for _, obj := range objects {
		if _, err := decodeDictEntryJSON(json.RawMessage(obj)); err == nil {
			return obj, nil
		}
		var batch struct {
			Entries []DictEntry `json:"entries"`
		}
		if err := json.Unmarshal([]byte(obj), &batch); err == nil && len(batch.Entries) > 0 {
			return obj, nil
		}
	}
	return "", fmt.Errorf("no entry-shaped object")
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if strings.TrimSpace(v) != "" {
			return strings.TrimSpace(v)
		}
	}
	return ""
}

func extractJSONObjectStrings(reply string) []string {
	reply = stripCodeFences(reply)
	var out []string
	for i := 0; i < len(reply); i++ {
		if reply[i] != '{' {
			continue
		}
		obj, next := scanJSONObject(reply, i)
		if obj != "" {
			out = append(out, obj)
		}
		if next <= i {
			i++
			continue
		}
		i = next - 1
	}
	trim := strings.TrimSpace(reply)
	if len(out) == 0 && strings.HasPrefix(trim, "{") {
		if obj, _ := scanJSONObject(trim, 0); obj != "" {
			out = append(out, obj)
		}
	}
	return out
}

func stripCodeFences(s string) string {
	s = strings.ReplaceAll(s, "```json", "")
	s = strings.ReplaceAll(s, "```JSON", "")
	s = strings.ReplaceAll(s, "```", "")
	return s
}

func scanJSONObject(s string, start int) (string, int) {
	if start >= len(s) || s[start] != '{' {
		return "", start + 1
	}
	depth := 0
	inStr := false
	esc := false
	for i := start; i < len(s); i++ {
		c := s[i]
		if inStr {
			if esc {
				esc = false
				continue
			}
			if c == '\\' {
				esc = true
				continue
			}
			if c == '"' {
				inStr = false
			}
			continue
		}
		switch c {
		case '"':
			inStr = true
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 {
				return strings.TrimSpace(s[start : i+1]), i + 1
			}
		}
	}
	return "", len(s)
}

// salvageUpsertsFromReply writes entries found as bare JSON in the model reply (no tool wrapper).
func salvageUpsertsFromReply(env *toolEnv, reply string) (int, string, error) {
	objects := extractJSONObjectStrings(reply)
	written := 0
	var files []string
	for _, obj := range objects {
		entries, err := decodeDictEntries(json.RawMessage(obj), obj)
		if err != nil {
			continue
		}
		var probe struct {
			Tool string `json:"tool"`
			Done bool   `json:"done"`
		}
		_ = json.Unmarshal([]byte(obj), &probe)
		if probe.Tool != "" && normalizeToolName(probe.Tool) != "upsert_entry" {
			continue
		}
		if probe.Done {
			continue
		}
		for _, entry := range entries {
			out, err := env.upsertEntryTool(json.RawMessage(mustMarshal(entry)), "")
			if err != nil {
				continue
			}
			written++
			if ws, ok := out["written"].([]string); ok {
				files = append(files, ws...)
			}
		}
	}
	if written == 0 {
		return 0, "", fmt.Errorf("no salvageable entries")
	}
	data, _ := json.Marshal(map[string]any{"written": files, "count": written, "salvaged": true})
	return written, string(data), nil
}

func mustMarshal(v any) string {
	b, _ := json.Marshal(v)
	return string(b)
}
