package studyprint

import (
	"fmt"
	"html"
	"strings"

	"boanerges/internal/bible/read"
	"boanerges/internal/bible/schema"
	"boanerges/internal/settings"
	"boanerges/internal/studydata"
)

// BuildStudyHTML renders printable study HTML for the primary bible column (and optional notes).
func BuildStudyHTML(dbPath string, study settings.StudySettings, readability settings.ReadabilitySettings, includeUserWork bool) (string, error) {
	mgr, err := schema.Open(dbPath)
	if err != nil {
		return "", err
	}
	defer mgr.Close()

	reader := read.NewReader(mgr.DB(), mgr, nil)
	ch, err := reader.Chapter(study.TranslationID, study.BookID, study.Chapter)
	if err != nil {
		return "", err
	}

	font := "Georgia, serif"
	if readability.FontFamily == "sans-serif" {
		font = "system-ui, sans-serif"
	}
	align := "left"
	if readability.JustifyText {
		align = "justify"
	}

	var b strings.Builder
	fmt.Fprintf(&b, `<!DOCTYPE html><html><head><meta charset="utf-8"><title>%s %d</title>`,
		html.EscapeString(ch.Book), ch.Chapter)
	fmt.Fprintf(&b, `<style>body{font-family:%s;font-size:%dpx;line-height:%.2f;text-align:%s;margin:1in;} h1{font-size:1.25em;} .verse{margin:0.25em 0;}</style></head><body>`,
		font, readability.FontSize, readability.LineHeight, align)
	fmt.Fprintf(&b, `<h1>%s %d</h1>`, html.EscapeString(ch.Book), ch.Chapter)
	for _, v := range ch.Verses {
		fmt.Fprintf(&b, `<p class="verse"><sup>%d</sup> %s</p>`, v.Number, v.Text)
	}

	if includeUserWork {
		svc, err := studydata.NewService()
		if err != nil {
			return b.String(), nil
		}
		text, err := svc.NotesGet(study.BookID, study.Chapter)
		if err == nil && strings.TrimSpace(text) != "" {
			b.WriteString(`<h2>Notes</h2><pre>`)
			b.WriteString(html.EscapeString(text))
			b.WriteString(`</pre>`)
		}
	}

	b.WriteString(`</body></html>`)
	return b.String(), nil
}
