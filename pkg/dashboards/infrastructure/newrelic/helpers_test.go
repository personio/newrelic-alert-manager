package newrelic_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/fpetkovski/newrelic-alert-manager/pkg/dashboards/domain"
	"io/ioutil"
	"net/http"
)

func newDashboard(title string) *domain.Dashboard {
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
	dashboard := newDashboard(title)
	dashboard.DashboardBody.Id = new(int64)
	*dashboard.DashboardBody.Id = int64(id)

	return dashboard
}

func newRequest(title string) []byte {
	request := []byte(fmt.Sprintf(	`
		{
			"dashboard": {
				"title": "%s",
				"editable": "read_only",
				"visibility": "all",
				"metadata": {
					"version": 1
				}
			}
		}
	`, title))

	buffer := new(bytes.Buffer)
	if err := json.Compact(buffer, request); err != nil {
		panic(err)
	}

	return buffer.Bytes()
}

func newResponse(id int64, title string) *http.Response {
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

func new404Response() *http.Response {
	return &http.Response{
		StatusCode: 404,
	}
}