package studyprint

import (
	_ "embed"
	"fmt"
	"strings"

	"boanerges/internal/settings"
)

//go:embed study-print.css
var studyPrintCSS string

func printStyles(study settings.StudySettings, readability settings.ReadabilitySettings) string {
	lineStepPx := float64(readability.FontSize) * readability.LineHeight
	align := "left"
	if readability.JustifyText {
		align = "justify"
	}

	var b strings.Builder
	fmt.Fprintf(&b, `@page{size:A4 %s;margin:0.5in;}`, pageOrientation(study.ColumnCount > 1))
	fmt.Fprintf(&b, `.study-print-body{--study-print-font:%s;--study-print-font-size:%dpx;--study-print-line-height:%.2f;--study-print-text-align:%s;--study-print-lined-step:%.2fpx;}`,
		settings.ReaderFontStack(readability.FontFamily),
		readability.FontSize,
		readability.LineHeight,
		align,
		lineStepPx,
	)
	b.WriteString(studyPrintCSS)
	return b.String()
}
