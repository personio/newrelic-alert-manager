package widget

type Data struct {
	*ApmMetric
	Nrql   string `json:"nrql,omitempty"`
	Source string `json:"source,omitempty"`
}

func (d Data) Equals(other Data) bool {
	equalSource := d.Source == other.Source
	equalNrql := d.Nrql == other.Nrql
	equalApm := d.compareApm(other)

	return equalNrql && equalApm && equalSource
}

func (d Data) compareApm(other Data) bool {
	if d.ApmMetric == nil && other.ApmMetric == nil {
		return true
	}

	if d.ApmMetric == nil && other.ApmMetric != nil {
		return false
	}

	return d.ApmMetric.Equals(other.ApmMetric)
}
