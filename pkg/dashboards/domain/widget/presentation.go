package widget

type Presentation struct {
	Title string `json:"title"`
	Notes string `json:"notes,omitempty"`
}

func (p Presentation) Equals(other Presentation) bool {
	return p.Title == other.Title && p.Notes == other.Notes
}
