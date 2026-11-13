//go:build agentic

package agentic

import (
	"fmt"
	"strings"
)

func validateDictEntry(entry DictEntry, jobInstruction string) error {
	hw := strings.TrimSpace(entry.Headword)
	if len(hw) < 2 {
		return fmt.Errorf("headword too short %q", hw)
	}
	low := strings.ToLower(hw)
	if low == "is" || low == "a" || low == "an" || low == "the" || low == "distill" {
		return fmt.Errorf("reject function word or meta headword %q", hw)
	}
	inst := strings.TrimSpace(jobInstruction)
	if inst != "" {
		if strings.EqualFold(hw, "distill") || strings.Contains(inst, hw) && len(hw) < 20 && len(entry.Modern) > 120 {
			return fmt.Errorf("headword looks like job instruction echo")
		}
		if len(entry.Modern) > 80 && strings.Contains(entry.Modern, "upsert_entry") {
			return fmt.Errorf("modern text looks like job instructions, not a definition")
		}
	}
	if len(entry.Senses) == 0 && entry.Modern == "" {
		return fmt.Errorf("need modern gloss or senses")
	}
	for _, s := range entry.Senses {
		if span := s.SourceSpan; span != "" && strings.Count(span, "0-") > 0 && len(span) > 12 {
			return fmt.Errorf("sourceSpan too broad %q — cite chunk offsets only", span)
		}
	}
	return nil
}
