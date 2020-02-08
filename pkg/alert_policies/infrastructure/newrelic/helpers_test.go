package newrelic_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/fpetkovski/newrelic-operator/pkg/alert_policies/domain"
	"io/ioutil"
	"net/http"
	"strings"
)

func newEmptyPolicy(name string) *domain.AlertPolicy {
	policy := &domain.AlertPolicy{
		Policy: domain.Policy{
			Id:                 nil,
			Name:               name,
			IncidentPreference: "per_policy",
		},
		NrqlConditions: []*domain.NrqlCondition{},
	}

	return policy
}

func newEmptyPolicyWithId(id int64, name string) *domain.AlertPolicy {
	policy := newEmptyPolicy(name)

	policyId := new(int64)
	*policyId = id
	policy.Policy.Id = policyId

	return policy
}

func newRequest(name string) []byte {
	request := []byte(fmt.Sprintf(`
		{
			"policy": {
				"name": "%s",
				"incident_preference": "per_policy"
			}
		}
	`, name))

	buffer := new(bytes.Buffer)
	if err := json.Compact(buffer, request); err != nil {
		panic(err)
	}

	return buffer.Bytes()
}

func newResponse(id int64, name string) *http.Response {
	response := []byte(fmt.Sprintf(`
		{
			"policy": {
				"id": %d,
				"name": "%s",
				"incident_preference": "per_policy"
			}
		}
	`, id, name))

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

func newArrayResponse(id int64, name string) *http.Response {
	response := []byte(fmt.Sprintf(`
		{
			"policies": [{
				"id": %d,
				"name": "%s",
				"incident_preference": "per_policy"
			}]
		}
	`, id, name))

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

func newStringResponse(response string) *http.Response {
	return &http.Response{
		StatusCode: 200,
		Body:       ioutil.NopCloser(strings.NewReader(response)),
		Close:      false,
	}
}
