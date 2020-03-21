package v1alpha1

import "github.com/fpetkovski/newrelic-alert-manager/pkg/dashboards/domain/widget"

type Widget struct {
	Title string `json:"title"`

	// +kubebuilder:validation:Enum=event_table;line_chart;facet_table;facet_bar_chart;facet_pie_chart;billboard;faceted_area_chart;faceted_line_chart;event_table;comparison_line_chart;heatmap;histogram;billboard_comparison;attribute_sheet;funnel;gauge;json;list
	Visualization string        `json:"visualization"`
	Data          Data          `json:"data"`
	Layout        widget.Layout `json:"layout"`
}

type Data struct {
	Nrql string `json:"nrql"`
}
