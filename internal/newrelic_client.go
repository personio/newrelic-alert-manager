package internal

import (
	"bytes"
	"fmt"
	"github.com/go-logr/logr"
	"net/http"
)

type NewrelicClient struct {
	client *http.Client
	log    logr.Logger

	url      string
	adminKey string
}

func NewNewrelicClient(log logr.Logger, url string, adminKey string) *NewrelicClient {
	return &NewrelicClient{
		client:   &http.Client{},
		log:      log,
		url:      url,
		adminKey: adminKey,
	}
}

func (newrelic NewrelicClient) Get(path string) (*http.Response, error) {
	request := newrelic.newRequest("Get", path, nil)

	return newrelic.execute(request)
}

func (newrelic NewrelicClient) PostJson(path string, payload []byte) (*http.Response, error) {
	request := newrelic.newJsonRequest("POST", path, payload)

	return newrelic.execute(request)
}

func (newrelic NewrelicClient) PutJson(path string, payload []byte) (*http.Response, error) {
	request := newrelic.newJsonRequest("PUT", path, payload)

	return newrelic.execute(request)
}

func (newrelic NewrelicClient) Delete(path string) (*http.Response, error) {
	request := newrelic.newJsonRequest("DELETE", path, nil)

	response, err := newrelic.execute(request)
	if response != nil && response.Status == "404" {
		return response, nil
	}
	if err != nil {
		return response, err
	}

	return response, nil
}

func (newrelic *NewrelicClient) newRequest(method string, path string, body []byte) *http.Request {
	var req *http.Request
	if body == nil {
		req = newJsonRequest(method, newrelic.url, path)
	} else {
		req = newJsonRequestWithBody(method, newrelic.url, path, body)
	}

	req.Header.Add("X-Api-Key", newrelic.adminKey)

	return req
}

func (newrelic *NewrelicClient) newJsonRequest(method string, path string, body []byte) *http.Request {
	var req *http.Request
	if body == nil {
		req = newJsonRequest(method, newrelic.url, path)
	} else {
		req = newJsonRequestWithBody(method, newrelic.url, path, body)
	}

	req.Header.Add("X-Api-Key", newrelic.adminKey)
	req.Header.Add("Content-Type", "application/json")

	return req
}

func newJsonRequestWithBody(method string, url string, path string, body []byte) *http.Request {
	req, _ := http.NewRequest(
		method,
		fmt.Sprintf("%s/%s", url, path),
		bytes.NewBuffer(body),
	)
	return req
}

func newJsonRequest(method string, url string, path string) *http.Request {
	req, _ := http.NewRequest(
		method,
		fmt.Sprintf("%s/%s", url, path),
		nil,
	)
	return req
}

func (newrelic NewrelicClient) execute(request *http.Request) (*http.Response, error) {
	newrelic.log.Info("Executing request", "Method", request.Method, "Endpoint", request.URL, "Payload", request.Body)
	response, err := newrelic.client.Do(request)
	if err != nil {
		return nil, err
	}

	return response, nil
}
