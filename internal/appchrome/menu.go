package appchrome

import "github.com/egoist/mygo"

// MenuAction is emitted when the user picks a custom application menu item.
var MenuAction = mygo.NewEvent[string]("studyMenu")
