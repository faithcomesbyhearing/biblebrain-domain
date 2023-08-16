package domain

type Media struct {
	Id      string
	BibleId string
}

func (m Media) GetBibleId() string {
	// FIXME: this is deprecated and should not be used if a validated BibleId is available
	if m.BibleId == "" {
		return m.Id[:6]
	}

	return m.BibleId
}
