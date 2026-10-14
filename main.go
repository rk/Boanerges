package main

import (
	"log"
	"runtime"

	"boanerges/internal/bible/bootstrap"
	"boanerges/internal/appchrome"
	"boanerges/internal/bibleapp"
	"boanerges/internal/settings"
	"boanerges/internal/studydata"

	"github.com/egoist/mygo"
)

func main() {
	boot := bootstrap.NewService()
	dbPath, err := bibleapp.ResolveUserDBPath()
	if err != nil {
		log.Fatal(err)
	}

	chrome := appchrome.NewAppChrome(dbPath)

	translations := bibleapp.NewTranslationService(dbPath, boot)
	mygo.Bind(
		settings.NewService(&settings.Repository{}),
		bibleapp.NewBibleService(dbPath, boot),
		translations,
		bibleapp.NewDictionaryService(dbPath),
		&studydata.StudyData{},
		chrome,
	)

	mygo.App.WhenReady(func() {
		win := mygo.NewWindow(mygo.WindowOptions{
			Title:           "Boanerges",
			Width:           960,
			Height:          640,
			MinWidth:        480,
			MinHeight:       360,
			BackgroundColor: "light-dark(#f6f7f9, #0f1115)",
			StateKey:        "main",
			URL:             "/",
			Page: mygo.PageOptions{
				DevTools: mygo.DevToolsEnabled,
			},
		})

		boot.SetWindow(win)
		chrome.SetWindow(win)
		appchrome.SetupApplicationMenu(win)
		go func() {
			resDir, err := mygo.App.Path(mygo.PathResources)
			if err != nil {
				log.Printf("bundled ASV import: %v", err)
				return
			}
			if err := boot.RunASV(dbPath, resDir); err != nil {
				log.Printf("bundled ASV import: %v", err)
			}
			if err := translations.ResumePendingInstalls(); err != nil {
				log.Printf("resume translation installs: %v", err)
			}
		}()
	})

	if err := mygo.App.Run(); err != nil {
		log.Fatal(err)
	}
}

// SystemInfo kept for mygo generate compatibility during transition.
type SystemInfo struct {
	OS        string `json:"os"`
	Arch      string `json:"arch"`
	GoVersion string `json:"goVersion"`
}

type Greeter struct{}

func (Greeter) Greet(name string) string {
	if name == "" {
		name = "stranger"
	}

	return "Hello, " + name + "! This message comes from Go."
}

func (Greeter) Info() SystemInfo {
	return SystemInfo{OS: runtime.GOOS, Arch: runtime.GOARCH, GoVersion: runtime.Version()}
}
