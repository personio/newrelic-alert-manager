package newrelic_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/fpetkovski/newrelic-alert-manager/pkg/dashboards/domain"
	"github.com/fpetkovski/newrelic-alert-manager/pkg/dashboards/domain/widget"
	"io/ioutil"
	"net/http"
)

func newEmptyDashboard(title string) *domain.Dashboard {
	dashboard := &domain.Dashboard{
		DashboardBody: domain.DashboardBody{
			Id:         nil,
			Title:      title,
			Visibility: "all",
			Editable:   "read_only",
			Metadata: domain.Metadata{
				Version: 1,
			},
		},
	}
	return dashboard
}

func newDashboardWithId(id int, title string) *domain.Dashboard {
	dashboard := newEmptyDashboard(title)
	dashboard.DashboardBody.Id = new(int64)
	*dashboard.DashboardBody.Id = int64(id)

	return dashboard
}

func newEmptyDashboardRequest(title string) []byte {
	request := []byte(fmt.Sprintf(`
		{
			"dashboard": {
				"title": "%s",
				"editable": "read_only",
				"visibility": "all",
				"metadata": {
					"version": 1
				},
				"widgets": null
			}
		}
	`, title))

	buffer := new(bytes.Buffer)
	if err := json.Compact(buffer, request); err != nil {
		panic(err)
	}

	return buffer.Bytes()
}

func newEmptyDashboardResponse(id int64, title string) *http.Response {
	response := []byte(fmt.Sprintf(`
		{
			"dashboard": {
				"id": %d,
				"title": "%s",
				"visibility": "all",
				"editable": "read_only",
				"metadata": {
					"version": 1
				}
			}
		}
	`, id, title))

	buffer := new(bytes.Buffer)
	if err := json.Compact(buffer, response); err != nil {
		panic(err)
	}

	return &http.Response{
		StatusCode: 200,
		Body:       ioutil.NopCloser(bytes.NewReader(buffer.Bytes())),
		Close:      false,
	}
}


func newApmDashboard(title string, metric string, entityId int) *domain.Dashboard {
	var dashboard = &domain.Dashboard{
		DashboardBody: domain.DashboardBody{
			Id:         nil,
			Title:      title,
			Visibility: "all",
			Editable:   "read_only",
			Metadata: domain.Metadata{
				Version: 1,
			},
			Widgets: []widget.Widget{
				{
					Visualization: "line_chart",
					Data: widget.DataList{
						{
							ApmMetric: &widget.ApmMetric{
								Duration: 60,
								EntityIds: []int{
									entityId,
								},
								Metrics: []widget.Metric{
									{
										Name:   metric,
										Values: []string{"value"},
									},
								},
								Facet:   "facet",
								OrderBy: "order",
							},
							Nrql: "",
						},
					},
					Layout: widget.Layout{
						Width:  0,
						Height: 0,
						Row:    0,
						Column: 0,
					},
					Presentation: widget.Presentation{
						Title: "test-title",
					},
				},
			},
		},
	}
	return dashboard
}

func newApmDashboardWithId(id int64, title string, metric string, entityId int) *domain.Dashboard {
	dashboard := newApmDashboard(title, metric, entityId)
	dashboard.DashboardBody.Id = new(int64)
	*dashboard.DashboardBody.Id = id

	return dashboard
}

func newApmDashboardRequest(title string, metric string, entityId int) []byte {
	request := []byte(fmt.Sprintf(`
		{
			"dashboard": {
				"title": "%s",
				"editable": "read_only",
				"visibility": "all",
				"metadata": {
					"version": 1
				},
				"widgets": [{
					"visualization": "line_chart",
					"data": [{
						"duration": 60,
						"entity_ids": [%d],
						"metrics": [{
							"name": "%s",
							"values": ["value"]
						}],
						"facet": "facet",
						"order_by": "order"
					}],
					"layout": {
						"width": 0,
						"height": 0,
						"row": 0,
						"column": 0
					},
					"presentation": {
						"title": "test-title"
					}
				}]
			}
		}
	`, title, entityId, metric))

	buffer := new(bytes.Buffer)
	if err := json.Compact(buffer, request); err != nil {
		panic(err)
	}

	return buffer.Bytes()
}

func newApmDashboardResponse(id int, title string, metric string, entityId int) *http.Response {
	response := []byte(fmt.Sprintf(`
		{
			"dashboard": {
				"id": %d,
				"title": "%s",
				"editable": "read_only",
				"visibility": "all",
				"metadata": {
					"version": 1
				},
				"widgets": [{
					"visualization": "line_chart",
					"data": [{
						"duration": 60,
						"entity_ids": [%d],
						"metrics": [{
							"name": "%s",
							"values": ["value"]
						}],
						"facet": "facet",
						"order_by": "order"
					}],
					"layout": {
						"width":  0,
						"height": 0,
						"row":    0,
						"column": 0
					},
					"presentation": {
						"title": "test-title"
					}
				}]
			}
		}
	`, id, title, entityId, metric))

	buffer := new(bytes.Buffer)
	if err := json.Compact(buffer, response); err != nil {
		panic(err)
	}

	return &http.Response{
		StatusCode: 200,
		Body:       ioutil.NopCloser(bytes.NewReader(buffer.Bytes())),
		Close:      false,
	}
}

func new404Response() *http.Response {
	return &http.Response{
		StatusCode: 404,
	}
}
