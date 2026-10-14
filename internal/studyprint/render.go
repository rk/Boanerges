package studyprint

import (
	"fmt"
	"html"
	"net/url"
	"strings"

	"boanerges/internal/bible/read"
	"boanerges/internal/settings"
)

func renderDocument(study settings.StudySettings, readability settings.ReadabilitySettings, columns []printColumn) string {
	font := settings.ReaderFontStack(readability.FontFamily)
	landscape := study.ColumnCount > 1
	lineStepPx := float64(readability.FontSize) * readability.LineHeight
	ruleSVG := url.QueryEscape(fmt.Sprintf(
		`<svg xmlns="http://www.w3.org/2000/svg" width="4" height="%.2f" viewBox="0 0 4 %.2f"><rect x="0" y="%.2f" width="4" height="1" fill="#444"/></svg>`,
		lineStepPx, lineStepPx, lineStepPx-1,
	))
	align := "left"
	if readability.JustifyText {
		align = "justify"
	}

	var b strings.Builder
	b.WriteString("<!DOCTYPE html><html lang=\"en\"><head><meta charset=\"utf-8\"><title>Study Print</title><style>")
	fmt.Fprintf(&b, `@page{size:A4 %s;margin:12mm;}* {box-sizing:border-box;margin:0;padding:0;}`, pageOrientation(landscape))
	fmt.Fprintf(&b, `body{font-family:%s;font-size:%dpx;line-height:%.2f;text-align:%s;color:#111;-webkit-print-color-adjust:exact;print-color-adjust:exact;}`, font, readability.FontSize, readability.LineHeight, align)
	b.WriteString(`.columns{display:flex;width:100%%;min-height:calc(100vh - 24mm);align-items:stretch;}`)
	b.WriteString(`.column{display:flex;flex-direction:column;flex:1 1 0;min-width:0;min-height:calc(100vh - 24mm);padding:0 4mm;border-right:1px solid #ccc;}`)
	b.WriteString(`.column:last-child{border-right:none;}`)
	b.WriteString(`.column-title{font-weight:600;font-size:0.95em;margin-bottom:0.75em;padding-bottom:0.35em;border-bottom:1px solid #ddd;}`)
	b.WriteString(`.column-subtitle{font-size:0.75em;color:#666;margin-bottom:0.75em;}`)
	b.WriteString(`.column-bible .reader p,.column-scribe-content .reader p{margin-bottom:0.75em;}`)
	b.WriteString(`.reader sup{font-size:0.7em;color:#555;margin-right:0.15em;}`)
	b.WriteString(`.notes-content{white-space:pre-wrap;}`)
	b.WriteString(`.lined-area{position:relative;flex:1 1 auto;min-height:0;}`)
	fmt.Fprintf(&b, `.lined-block{position:absolute;inset:0;background-image:url('data:image/svg+xml,%s');background-size:100%% %.2fpx;background-repeat:repeat-y;}`, ruleSVG, lineStepPx)
	b.WriteString(`.message{color:#666;font-size:0.9em;font-style:italic;}`)
	b.WriteString(`.verse-list-entry{margin-bottom:0.75em;}`)
	b.WriteString(`.verse-list-entry strong{display:block;margin-bottom:0.25em;}`)
	b.WriteString(`</style></head><body><div class="columns">`)

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

	b.WriteString(`</div></body></html>`)
	return b.String()
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
