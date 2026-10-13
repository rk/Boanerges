package appchrome

import (
	"strconv"

	"github.com/egoist/mygo"
)

const feedbackMailto = "mailto:feedback@boanerges.app?subject=Boanerges%20Feedback"

func SetupApplicationMenu(win *mygo.Window) {
	emit := func(id string) func(*mygo.MenuItem, *mygo.Window) {
		return func(_ *mygo.MenuItem, w *mygo.Window) {
			target := win
			if w != nil {
				target = w
			}
			if target != nil {
				_ = MenuAction.Emit(target, id)
			}
		}
	}

	menu := mygo.NewMenu([]*mygo.MenuItem{
		{Role: mygo.RoleAppMenu},
		{
			Label: "File",
			Submenu: []*mygo.MenuItem{
				{Label: "Manage Translations…", ID: "file.manage-translations", Click: emit("file.manage-translations")},
				mygo.Separator(),
				{Role: mygo.RoleQuit},
			},
		},
		{Role: mygo.RoleEditMenu},
		{
			Label: "Study",
			Submenu: []*mygo.MenuItem{
				{Label: "Notes", ID: "study.notes", Click: emit("study.notes")},
				{Label: "Scribe", ID: "study.scribe", Click: emit("study.scribe")},
				mygo.Separator(),
				{Label: "Search", ID: "study.search", Accelerator: "CmdOrCtrl+F", Click: emit("study.search")},
				{Label: "Cross-References", ID: "study.cross-references", Accelerator: "CmdOrCtrl+Shift+R", Click: emit("study.cross-references")},
				{Label: "Dictionary", ID: "study.dictionary", Accelerator: "CmdOrCtrl+Shift+D", Click: emit("study.dictionary")},
				{Label: "Comparison", ID: "study.comparison", Click: emit("study.comparison")},
				{Label: "Verse List", ID: "study.verse-list", Click: emit("study.verse-list")},
				mygo.Separator(),
				{Label: "Feedback", ID: "study.feedback", Click: emit("study.feedback")},
				mygo.Separator(),
				{Label: "Print", ID: "study.print", Accelerator: "CmdOrCtrl+P", Click: emit("study.print")},
			},
		},
		{
			Label: "View",
			Submenu: []*mygo.MenuItem{
				{Label: "1 Column", ID: "view.columns.1", Type: mygo.MenuItemRadio, Click: emit("view.columns.1")},
				{Label: "2 Columns", ID: "view.columns.2", Type: mygo.MenuItemRadio, Click: emit("view.columns.2")},
				{Label: "3 Columns", ID: "view.columns.3", Type: mygo.MenuItemRadio, Click: emit("view.columns.3")},
				mygo.Separator(),
				{Label: "Link scroll", ID: "view.scroll-sync", Type: mygo.MenuItemCheckbox, Click: emit("view.scroll-sync")},
				mygo.Separator(),
				{Label: "Edit Settings", ID: "view.settings", Accelerator: "CmdOrCtrl+,", Click: emit("view.settings")},
				mygo.Separator(),
				{Label: "Developer Tools", ID: "view.devtools", Click: emit("view.devtools")},
				mygo.Separator(),
				{Role: mygo.RoleToggleFullScreen},
			},
		},
		{Role: mygo.RoleWindowMenu},
	})

	mygo.App.SetMenu(menu)
}

// SyncMenuState updates radio/checkbox items to match persisted study UI state.
func SyncMenuState(columnCount int, scrollSync bool) {
	menu := mygo.App.Menu()
	if menu == nil {
		return
	}
	for n := 1; n <= 3; n++ {
		if item := menu.ItemByID("view.columns." + strconv.Itoa(n)); item != nil {
			item.SetChecked(columnCount == n)
		}
	}
	if item := menu.ItemByID("view.scroll-sync"); item != nil {
		item.SetChecked(scrollSync)
	}
}
