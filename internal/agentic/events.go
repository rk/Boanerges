//go:build agentic

package agentic

import "github.com/egoist/mygo"

// RunStateChanged notifies the Studio UI to refresh run logs.
var RunStateChanged = mygo.NewEvent[struct{}]("agenticRunStateChanged")
