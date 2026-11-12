//go:build agentic

package agentic

import (
	"encoding/json"
	"fmt"
	"strings"
)

// decodeDictEntries accepts several model JSON shapes for upsert_entry.
func decodeDictEntries(args json.RawMessage, toolLine string) ([]DictEntry, error) {
	if len(args) > 0 && string(args) != "null" {
		var batch struct {
			Entries []DictEntry `json:"entries"`
		}
		if err := json.Unmarshal(args, &batch); err == nil && len(batch.Entries) > 0 {
			return normalizeEntries(batch.Entries)
		}
		if e, err := decodeDictEntryJSON(args); err == nil {
			return []DictEntry{e}, nil
		}
	}
	if strings.TrimSpace(toolLine) != "" {
		if list, err := decodeDictEntryFromToolLine(toolLine); err == nil {
			return list, nil
		}
	}
	return nil, fmt.Errorf("could not parse entry: need headword in args or on the tool line")
}

func decodeDictEntryJSON(raw json.RawMessage) (DictEntry, error) {
	var e DictEntry
	if err := json.Unmarshal(raw, &e); err == nil {
		e = normalizeEntry(e)
		if e.Headword != "" {
			return e, nil
		}
	}
	var wrap struct {
		Entry DictEntry `json:"entry"`
	}
	if err := json.Unmarshal(raw, &wrap); err == nil {
		wrap.Entry = normalizeEntry(wrap.Entry)
		if wrap.Entry.Headword != "" {
			return wrap.Entry, nil
		}
	}
	var loose struct {
		Headword   string      `json:"headword"`
		Word       string      `json:"word"`
		Title      string      `json:"title"`
		Modern     string      `json:"modern"`
		Definition string      `json:"definition"`
		Gloss      string      `json:"gloss"`
		Senses     []DictSense `json:"senses"`
		SourceFile string      `json:"sourceFile"`
		SourceSpan string      `json:"sourceSpan"`
	}
	if err := json.Unmarshal(raw, &loose); err != nil {
		return DictEntry{}, err
	}
	hw := strings.TrimSpace(loose.Headword)
	if hw == "" {
		hw = strings.TrimSpace(loose.Word)
	}
	if hw == "" {
		hw = strings.TrimSpace(loose.Title)
	}
	if hw == "" {
		return DictEntry{}, fmt.Errorf("headword required")
	}
	e = DictEntry{Headword: hw, Modern: loose.Modern, Senses: loose.Senses}
	if len(e.Senses) == 0 {
		gloss := loose.Gloss
		if gloss == "" {
			gloss = loose.Definition
		}
		if gloss == "" && loose.Modern != "" {
			gloss = loose.Modern
		}
		if gloss != "" {
			e.Senses = []DictSense{{
				Gloss:      gloss,
				SourceFile: loose.SourceFile,
				SourceSpan: loose.SourceSpan,
			}}
		}
	}
	return normalizeEntry(e), nil
}

func decodeDictEntryFromToolLine(line string) ([]DictEntry, error) {
	var probe struct {
		Tool    string          `json:"tool"`
		Args    json.RawMessage `json:"args"`
		Entries []DictEntry     `json:"entries"`
		DictEntry
	}
	if err := json.Unmarshal([]byte(line), &probe); err != nil {
		return nil, err
	}
	if len(probe.Entries) > 0 {
		return normalizeEntries(probe.Entries)
	}
	if len(probe.Args) > 0 {
		return decodeDictEntries(probe.Args, "")
	}
	probe.DictEntry = normalizeEntry(probe.DictEntry)
	if probe.DictEntry.Headword != "" {
		return []DictEntry{probe.DictEntry}, nil
	}
	return nil, fmt.Errorf("headword required")
}

func normalizeEntries(in []DictEntry) ([]DictEntry, error) {
	out := make([]DictEntry, 0, len(in))
	for _, e := range in {
		e = normalizeEntry(e)
		if e.Headword == "" {
			return nil, fmt.Errorf("headword required")
		}
		out = append(out, e)
	}
	return out, nil
}

func normalizeEntry(e DictEntry) DictEntry {
	e.Headword = strings.TrimSpace(e.Headword)
	e.Modern = strings.TrimSpace(e.Modern)
	if e.Headword == "" {
		return e
	}
	if e.Modern == "" && len(e.Senses) > 0 {
		e.Modern = e.Senses[0].Gloss
	}
	return e
}
