package studyprint

import (
	"fmt"
	"html"
	"strings"

	"boanerges/internal/bible/read"
	"boanerges/internal/settings"
)

func printColumnsHTML(columns []printColumn) string {
	var b strings.Builder
	b.WriteString(`<div class="columns">`)

	for _, col := range columns {
		fmt.Fprintf(&b, `<div class="column column-%s">`, html.EscapeString(col.Kind))
		fmt.Fprintf(&b, `<div class="column-title">%s</div>`, html.EscapeString(col.Label))
		switch col.Kind {
		case "bible", "scribe-content":
			b.WriteString(`<div class="reader">`)
			writeChapterParagraphs(&b, col.Verses, col.Kind == "scribe-content")
			b.WriteString(`</div>`)
		case "notes":
			fmt.Fprintf(&b, `<div class="notes-content">%s</div>`, html.EscapeString(col.Content))
		case "lined-notes":
			b.WriteString(`<div class="column-subtitle">Notes</div><div class="lined-area"><div class="lined-block" aria-hidden="true"></div></div>`)
		case "scribe":
			b.WriteString(`<div class="column-subtitle">Scribe</div><div class="lined-area"><div class="lined-block" aria-hidden="true"></div></div>`)
		case "message":
			fmt.Fprintf(&b, `<p class="message">%s</p>`, html.EscapeString(col.Message))
		case "verse-list":
			for _, e := range col.Entries {
				b.WriteString(`<div class="verse-list-entry"><strong>`)
				b.WriteString(html.EscapeString(e.Label))
				b.WriteString(`</strong><p>`)
				b.WriteString(e.Text)
				b.WriteString(`</p></div>`)
			}
		}
		b.WriteString(`</div>`)
	}

	b.WriteString(`</div>`)
	return b.String()
}

func printBodyHTML(columns []printColumn) string {
	return `<div class="study-print-body">` + printColumnsHTML(columns) + `</div>`
}

func renderDocument(study settings.StudySettings, readability settings.ReadabilitySettings, columns []printColumn) string {
	styles := printStyles(study, readability)
	body := printBodyHTML(columns)
	return "<!DOCTYPE html><html lang=\"en\"><head><meta charset=\"utf-8\"><title>Study Print</title><style>" +
		styles + `body{margin:0;}` +
		`</style></head><body>` + body + `</body></html>`
}

func pageOrientation(landscape bool) string {
	if landscape {
		return "landscape"
	}

	return "portrait"
}

func writeChapterParagraphs(b *strings.Builder, verses []read.Verse, escapeText bool) {
	var paragraph []read.Verse
	flush := func() {
		if len(paragraph) == 0 {
			return
		}
		b.WriteString("<p>")
		for i, part := range paragraph {
			if i > 0 {
				b.WriteString("&nbsp;")
			}
			fmt.Fprintf(b, `<sup>%d</sup>`, part.Number)
			if escapeText {
				b.WriteString(html.EscapeString(part.Text))
			} else {
				b.WriteString(part.Text)
			}
		}
		b.WriteString("</p>")
		paragraph = nil
	}
	for _, verse := range verses {
		if verse.ParagraphStart && len(paragraph) > 0 {
			flush()
		}
		paragraph = append(paragraph, verse)
	}
	flush()
}
