package domain

import "github.com/fpetkovski/newrelic-alert-manager/pkg/dashboards/domain/widget"

type Dashboard struct {
	DashboardBody DashboardBody `json:"dashboard"`
}

func (d Dashboard) Equals(other Dashboard) bool {
	return d.DashboardBody.Equals(other.DashboardBody)
}

type DashboardList struct {
	Dashboards []DashboardBody `json:"dashboards"`
}

type DashboardBody struct {
	Id         *int64            `json:"id,omitempty"`
	Title      string            `json:"title"`
	Editable   string            `json:"editable"`
	Visibility string            `json:"visibility"`
	Metadata   Metadata          `json:"metadata"`
	Widgets    widget.WidgetList `json:"widgets"`
}

func (d DashboardBody) Equals(other DashboardBody) bool {
	return d.Title == other.Title && d.Widgets.Equals(other.Widgets)
}

type Metadata struct {
	Version int `json:"version"`
}
