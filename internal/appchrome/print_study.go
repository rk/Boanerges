package appchrome

import (
	"sync"
	"sync/atomic"

	"boanerges/internal/settings"

	"github.com/egoist/mygo"
)

// ponytail: MyGo Page.Print() does not block until the sheet dismisses; wait on afterprint before closing the window.
const printAndWaitJS = `return new Promise((resolve) => {
	window.addEventListener('afterprint', () => resolve(0), { once: true });
	window.print();
});`

// PrintStudy opens the system print dialog for the current study layout (legacy NativePHP System::print parity).
// Known broken — see TODO.md (Study print). Prefer a MyGo Page().Print() / PrintToPDF() design per official docs.
func (s *AppChrome) PrintStudy(study settings.StudySettings, readability settings.ReadabilitySettings, includeUserWork bool) error {
	path, err := s.WriteStudyPrintFile(study, readability, includeUserWork, "study-print-preview.html")
	if err != nil {
		return err
	}

	done := make(chan error, 1)
	printWin := mygo.NewWindow(mygo.WindowOptions{
		Title:  "Study Print",
		Hidden: true,
		Width:  800,
		Height: 600,
	})

	var once sync.Once
	var armed atomic.Bool
	printWin.Page().OnDOMReady(func() {
		if !armed.Load() {
			return
		}
		once.Do(func() {
			printWin.Show()
			_, evalErr := printWin.Page().Eval(printAndWaitJS)
			printWin.Close()
			done <- evalErr
		})
	})

	armed.Store(true)
	if err := printWin.Page().LoadFile(path); err != nil {
		printWin.Close()
		return err
	}

	return <-done
}
