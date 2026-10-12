package settings

// Service is bound for the web frontend (`mygo.Bind`).
type Service struct {
	readability *ReadabilityStore
	study       *StudyStore
}

func NewService(repo *Repository) *Service {
	return &Service{
		readability: NewReadabilityStore(repo),
		study:       NewStudyStore(repo),
	}
}

// GetReadability returns current readability settings.
func (s *Service) GetReadability() ReadabilitySettings {
	return s.readability.Get()
}

// UpdateReadability validates and saves readability settings.
func (s *Service) UpdateReadability(settings ReadabilitySettings) (ReadabilitySettings, error) {
	return s.readability.Update(settings)
}

// GetStudy returns current study layout settings.
func (s *Service) GetStudy() StudySettings {
	return s.study.Get()
}

// UpdateStudy validates and saves study settings.
func (s *Service) UpdateStudy(settings StudySettings) (StudySettings, error) {
	if err := ValidateStudy(settings); err != nil {
		return StudySettings{}, err
	}

	return s.study.Update(settings), nil
}
