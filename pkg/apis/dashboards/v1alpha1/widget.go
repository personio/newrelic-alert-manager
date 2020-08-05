package v1alpha1

import "github.com/personio/newrelic-alert-manager/pkg/dashboards/domain/widget"

// Widget defines the widget parameters \
// For more details, refer to the official [New Relic documentation](https://docs.newrelic.com/docs/insights/insights-api/manage-dashboards/insights-dashboard-api#widget-data)
type Widget struct {
	// The title of the widget created in New Relic
	Title string `json:"title"`
	// Visualization type to use for the widget. \
	// Available options are: \
	// - `application_breakdown` \
	// - `attribute_sheet` \
	// - `background_breakdown` \
	// - `billboard` \
	// - `billboard_comparison` \
	// - `comparison_line_chart` \
	// - `event_table` \
	// - `facet_bar_chart` \
	// - `facet_pie_chart` \
	// - `facet_table` \
	// - `faceted_area_chart` \
	// - `faceted_line_chart` \
	// - `funnel` \
	// - `gauge` \
	// - `heatmap` \
	// - `histogram` \
	// - `json` \
	// - `line_chart` \
	// - `list` \
	// - `metric_line_chart` (used for apm metrics) \
	// +kubebuilder:validation:Enum=application_breakdown;attribute_sheet;background_breakdown;billboard;billboard_comparison;comparison_line_chart;event_table;facet_bar_chart;facet_pie_chart;facet_table;faceted_area_chart;faceted_line_chart;funnel;gauge;heatmap;histogram;json;line_chart;list;metric_line_chart
	Visualization string `json:"visualization"`
	// Notes to add to the widget
	Notes string `json:"notes,omitempty"`
	// The data to plot on the widget
	Data Data `json:"data"`
	// Defines the layout of the widget within the dashboard
	Layout widget.Layout `json:"layout"`
}

// Data represents the data to plot inside the widget. \
// Either Nrql or ApmMetric should be specified, but not both. \
// \
// Leave both fields empty if you want to plot the application breakdown data, \
// also present in the main widget that comes with the default application dashboard. \
// For more information refer to the official [New Relic documentation](https://docs.newrelic.com/docs/insights/insights-api/manage-dashboards/insights-dashboard-api#dashboard-data)
type Data struct {
	// The NRQL query used which defines the data to plot in the widget
	// +optional
	Nrql string `json:"nrql,omitempty"`
	// The APM metric parameters which defines the data to plot in the widget. \
	// When using an APM metric for the data, visualization should be set to either `metric_line_chart` or `application_breakdown`. \
	// +optional
	ApmMetric *Apm `json:"apm,omitempty"`
}

// Apm is the set of metric parameters used for defining the data to plot in the widget
type Apm struct {
	// The time frame in seconds
	SinceSeconds int `json:"sinceSeconds,omitempty"`
	// A list of application names for which to get the metric
	Entities []string `json:"entities"`
	// A list of metrics to use
	Metrics []Metric `json:"metrics,omitempty"`
	// +optional
	Facet string `json:"facet,omitempty"`
	// +optional
	OrderBy string `json:"order_by,omitempty"`
}

// Metric is the name of the metric as shown in Data Explorer
type Metric struct {
	// Name of the metric
	Name string `json:"name"`
	// List of metric values to plot. The available values will depend on the metric you choose. \
	// Check the Data Explorer in New Relic to find out which values are available for which metrics.
	Values []string `json:"values"`
}
