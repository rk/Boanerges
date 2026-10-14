package settings

const readabilityKey = "readability"

// ReadabilitySettings match legacy readability store shape.
type ReadabilitySettings struct {
	FontSize    int     `json:"fontSize"`
	LineHeight  float64 `json:"lineHeight"`
	Theme       string  `json:"theme"`
	FontFamily  string  `json:"fontFamily"`
	JustifyText bool    `json:"justifyText"`
}

type ReadabilityStore struct {
	repo *Repository
}

func NewReadabilityStore(repo *Repository) *ReadabilityStore {
	return &ReadabilityStore{repo: repo}
}

func (s *ReadabilityStore) Defaults() ReadabilitySettings {
	return ReadabilitySettings{
		FontSize:    16,
		LineHeight:  1.2,
		Theme:       "auto",
		FontFamily:  "serif",
		JustifyText: true,
	}
}

func (s *ReadabilityStore) Get() ReadabilitySettings {
	raw := s.repo.Get(readabilityKey, readabilityDefaultsMap())
	out := mapToReadability(raw)
	if out.Theme == "sepia" {
		out.Theme = "caramellatte"
	}

	return out
}

func (s *ReadabilityStore) Update(patch ReadabilitySettings) (ReadabilitySettings, error) {
	if err := ValidateReadability(patch); err != nil {
		return ReadabilitySettings{}, err
	}

	merged := s.repo.Update(readabilityKey, readabilityToMap(patch), readabilityDefaultsMap())

	return mapToReadability(merged), nil
}

func readabilityDefaultsMap() map[string]any {
	d := NewReadabilityStore(&Repository{}).Defaults()

	return readabilityToMap(d)
}

func readabilityToMap(s ReadabilitySettings) map[string]any {
	return map[string]any{
		"fontSize":    s.FontSize,
		"lineHeight":  s.LineHeight,
		"theme":       s.Theme,
		"fontFamily":  s.FontFamily,
		"justifyText": s.JustifyText,
	}
}

func mapToReadability(m map[string]any) ReadabilitySettings {
	d := ReadabilitySettings{
		FontSize:    16,
		LineHeight:  1.2,
		Theme:       "auto",
		FontFamily:  "serif",
		JustifyText: true,
	}

	d.FontSize = asInt(m, "fontSize", d.FontSize)
	if v, ok := m["lineHeight"].(float64); ok {
		d.LineHeight = v
	}
	d.Theme = asString(m, "theme", d.Theme)
	d.FontFamily = asString(m, "fontFamily", d.FontFamily)
	d.JustifyText = asBool(m, "justifyText", d.JustifyText)

	return d
}
