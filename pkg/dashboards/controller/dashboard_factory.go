package controller

import (
	"github.com/fpetkovski/newrelic-alert-manager/pkg/apis/newrelic/v1alpha1"
	"github.com/fpetkovski/newrelic-alert-manager/pkg/dashboards/domain"
	"github.com/fpetkovski/newrelic-alert-manager/pkg/dashboards/domain/widget"
)

func NewDashboard(cr *v1alpha1.Dashboard) *domain.Dashboard {
	return &domain.Dashboard{
		DashboardBody: domain.DashboardBody{
			Id:         cr.Status.NewrelicDashboardId,
			Title:      cr.Spec.Title,
			Visibility: "all",
			Editable:   "read_only",
			Metadata: domain.Metadata{
				Version: 1,
			},
			Widgets: newWidgets(cr.Spec.Widgets),
		},
	}
}

func newWidgets(widgets []v1alpha1.Widget) widget.WidgetList {
	result := make(widget.WidgetList, len(widgets))
	for i, w := range widgets {
		result[i] = widget.Widget{
			Visualization: w.Visualization,
			Data:          newData(w.Data),
			Layout:        w.Layout,
			Presentation: widget.Presentation{
				Title: w.Title,
			},
		}
	}

	return result
}

func newData(w v1alpha1.Data) widget.DataList {
	var result widget.DataList
	result[0] = widget.Data{
		Nrql: w.Nrql,
	}

	return result
}
