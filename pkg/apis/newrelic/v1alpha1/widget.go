package v1alpha1

import "github.com/fpetkovski/newrelic-alert-manager/pkg/dashboards/domain/widget"

type Widget struct {
	Title string `json:"title"`

	// +kubebuilder:validation:Enum=application_breakdown;attribute_sheet;background_breakdown;billboard;billboard_comparison;comparison_line_chart;event_table;facet_bar_chart;facet_pie_chart;facet_table;faceted_area_chart;faceted_line_chart;funnel;gauge;heatmap;histogram;json;line_chart;list;metric_line_chart
	Visualization string        `json:"visualization"`
	Data          Data          `json:"data"`
	Layout        widget.Layout `json:"layout"`
}

type Data struct {
	Nrql      string `json:"nrql,omitempty"`
	ApmMetric *Apm   `json:"apm,omitempty"`
}

type Apm struct {
	Duration int      `json:"duration"`
	Entities []string `json:"entities"`
	Metrics  []Metric `json:"metrics"`
	Facet    string   `json:"facet,omitempty"`
	OrderBy  string   `json:"order_by,omitempty"`
}

type Metric struct {
	Name   string   `json:"name"`
	Values []string `json:"values"`
}
