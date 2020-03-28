package widget

import (
	"fmt"
)

type Widget struct {
	Visualization string       `json:"visualization"`
	Data          DataList     `json:"data"`
	Layout        Layout       `json:"layout"`
	Presentation  Presentation `json:"presentation"`
}

func (w Widget) getSortKey() string {
	return fmt.Sprintf("%s-%s-%s-%s", w.Data.getSortKey(), w.Layout.getHashKey(), w.Visualization, w.Presentation.getSortKey())
}

func (w Widget) Equals(other Widget) bool {
	return w.Visualization == other.Visualization &&
		w.Layout.Equals(other.Layout) &&
		w.Presentation.Equals(other.Presentation) &&
		w.Data.Equals(other.Data)
}
