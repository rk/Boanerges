package studyprint

import (
	"boanerges/internal/settings"
	"boanerges/internal/studydata"
)

// BuildStudyHTML renders printable multi-column study HTML (legacy StudyPrintHtmlBuilder parity).
func BuildStudyHTML(dbPath string, study settings.StudySettings, readability settings.ReadabilitySettings, includeUserWork bool) (string, error) {
	reader, trans, mgr, err := openStudyServices(dbPath)
	if err != nil {
		return "", err
	}
	defer mgr.Close()

	primary, err := reader.Chapter(study.TranslationID, study.BookID, study.Chapter)
	if err != nil {
		return "", err
	}

	var studySvc *studydata.Service
	if includeUserWork {
		studySvc, _ = studydata.NewService()
	}

	columns, err := buildColumns(reader, trans, studySvc, study, primary, includeUserWork)
	if err != nil {
		return "", err
	}

	return renderDocument(study, readability, columns), nil
}

// PrintDocument is printable study content split for embedding in a preview shell.
type PrintDocument struct {
	Styles string `json:"styles"`
	Body   string `json:"body"`
}

// BuildStudyPrintDocument renders study print CSS and HTML body (no preview chrome).
func BuildStudyPrintDocument(dbPath string, study settings.StudySettings, readability settings.ReadabilitySettings, includeUserWork bool) (PrintDocument, error) {
	reader, trans, mgr, err := openStudyServices(dbPath)
	if err != nil {
		return PrintDocument{}, err
	}
	defer mgr.Close()

	primary, err := reader.Chapter(study.TranslationID, study.BookID, study.Chapter)
	if err != nil {
		return PrintDocument{}, err
	}

	var studySvc *studydata.Service
	if includeUserWork {
		studySvc, _ = studydata.NewService()
	}

	columns, err := buildColumns(reader, trans, studySvc, study, primary, includeUserWork)
	if err != nil {
		return PrintDocument{}, err
	}

	return PrintDocument{
		Styles: printStyles(study, readability),
		Body:   printBodyHTML(columns),
	}, nil
}
