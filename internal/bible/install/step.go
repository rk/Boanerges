package install

// Step matches legacy App\Enums\TranslationInstallStep values.
type Step string

const (
	StepPending        Step = "pending"
	StepQueued         Step = "queued"
	StepCreatingSchema Step = "creating_schema"
	StepImporting      Step = "importing"
	StepIndexing       Step = "indexing"
	StepReady          Step = "ready"
	StepFailed         Step = "failed"
)

func (s Step) Percent() int {
	switch s {
	case StepReady:
		return 100
	case StepIndexing:
		return 85
	case StepImporting:
		return 50
	case StepCreatingSchema:
		return 30
	default:
		return 0
	}
}
