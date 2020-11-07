package widget

type Widget struct {
	Visualization string       `json:"visualization"`
	Data          DataList     `json:"data"`
	Layout        Layout       `json:"layout"`
	Presentation  Presentation `json:"presentation"`
}

func (w Widget) Equals(other Widget) bool {
	return w.Visualization == other.Visualization &&
		w.Layout.Equals(other.Layout) &&
		w.Presentation.Equals(other.Presentation) &&
		w.Data.Equals(other.Data)
}

func (w Widget) getComparisonKey() int {
	return w.Layout.Column * 1000 + w.Layout.Row * 100 + w.Layout.Width * 10 * w.Layout.Height
}
