package widget

type Data struct {
	Nrql string `json:"nrql"`
}

func (d Data) getHashKey() string {
	return d.Nrql
}

func (d Data) Equals(other Data) bool {
	return d.Nrql == other.Nrql
}