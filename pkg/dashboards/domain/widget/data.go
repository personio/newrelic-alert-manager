package widget

type Data struct {
	*ApmMetric
	Nrql string `json:"nrql,omitempty"`
}

func (d Data) Equals(other Data) bool {
	equalNrql := d.Nrql == other.Nrql
	equalApm := d.compareApm(other)

	return equalNrql && equalApm
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
