//go:build !agentic

package appchrome

import "github.com/egoist/mygo"

func agenticFileMenuItems(func(string) func(*mygo.MenuItem, *mygo.Window)) []*mygo.MenuItem {
	return nil
}
