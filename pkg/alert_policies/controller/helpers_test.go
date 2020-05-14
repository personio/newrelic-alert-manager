package controller_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/personio/newrelic-alert-manager/pkg/apis/alerts/v1alpha1"
	"io/ioutil"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
)

func newPolicyWithApmCondition(policyName string, entityName string) *v1alpha1.AlertPolicy {
	return &v1alpha1.AlertPolicy{
		ObjectMeta: v1.ObjectMeta{
			Name: policyName,
		},
		Spec: v1alpha1.AlertPolicySpec{
			Name:               policyName,
			IncidentPreference: "per_policy",
			NrqlConditions:     nil,
			ApmConditions: []v1alpha1.ApmCondition{
				{
					Name:     "condition",
					Type:     "apm_metric",
					Entities: []string{entityName},
				},
			},
			InfraConditions: nil,
		},
	}
}

func newEmptyResponse() *http.Response {
	response := []byte(fmt.Sprintf(`
		{
			"applications": [{}]
		}
	`))

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

func newResponse(appId int, appName string) *http.Response {
	response := []byte(fmt.Sprintf(`
		{
			"applications": [{
				"id": %d,
				"name": "%s"
			}]
		}
	`, appId, appName))

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
