package applications

type ApplicationList struct {
	Applications []Application `json:"applications"`
}
type Application struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}
