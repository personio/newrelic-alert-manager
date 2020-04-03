package widget

import "fmt"

type Layout struct {
	Width  int `json:"width"`
	Height int `json:"height"`
	Row    int `json:"row"`
	Column int `json:"column"`
}

func (l Layout) getHashKey() string {
	return fmt.Sprintf("%d-%d-%d-%d", l.Width, l.Height, l.Row, l.Column)
}

func (l Layout) Equals(other Layout) bool {
	return l.Width == other.Width &&
		l.Height == other.Height &&
		l.Row == other.Row &&
		l.Column == other.Column
}
