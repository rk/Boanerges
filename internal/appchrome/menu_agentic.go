//go:build agentic

package appchrome

import "github.com/egoist/mygo"

func agenticFileMenuItems(emit func(string) func(*mygo.MenuItem, *mygo.Window)) []*mygo.MenuItem {
	return []*mygo.MenuItem{
		{Label: "Agentic Studio…", ID: "file.agentic-studio", Click: emit("file.agentic-studio")},
		mygo.Separator(),
	}
}
