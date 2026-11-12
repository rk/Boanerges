//go:build !agentic

package appchrome

import "errors"

// OpenAgenticStudio is unavailable in release builds.
func (s *AppChrome) OpenAgenticStudio() error {
	return errors.New("agentic studio is not included in this build")
}
