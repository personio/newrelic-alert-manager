package newrelic

import (
	"bytes"
	"fmt"
	"github.com/go-logr/logr"
	"net/http"
)

type Client struct {
	client *http.Client
	log    logr.Logger

	url      string
	adminKey string
}

func NewClient(log logr.Logger, url string, adminKey string) *Client {
	return &Client{
		client:   &http.Client{},
		log:      log,
		url:      url,
		adminKey: adminKey,
	}
}

func (newrelic Client) Post(path string, payload []byte) (*http.Response, error) {
	request := newrelic.newRequest("POST", path, payload)

	return newrelic.execute(request, payload)
}

func (newrelic Client) Put(path string, payload []byte) (*http.Response, error) {
	request := newrelic.newRequest("PUT", path, payload)

	return newrelic.execute(request, payload)
}

func (newrelic Client) Delete(path string) (*http.Response, error) {
	request := newrelic.newRequest("DELETE", path, nil)

	response, err := newrelic.execute(request, nil)
	if response != nil && response.Status == "404" {
		return response, nil
	}
	if err != nil {
		return response, err
	}

	return response, nil
}

func (newrelic *Client) newRequest(method string, path string, body []byte) *http.Request {
	var req *http.Request
	if body == nil {
		req = newRequest(method, newrelic.url, path)
	} else {
		req = newRequestWithBody(method, newrelic.url, path, body)
	}

	req.Header.Add("X-Api-Key", newrelic.adminKey)
	req.Header.Add("Content-Type", "application/json")

	return req
}

func newRequestWithBody(method string, url string, path string, body []byte) *http.Request {
	req, _ := http.NewRequest(
		method,
		fmt.Sprintf("%s/%s", url, path),
		bytes.NewBuffer(body),
	)
	return req
}

func newRequest(method string, url string, path string) *http.Request {
	req, _ := http.NewRequest(
		method,
		fmt.Sprintf("%s/%s", url, path),
		nil,
	)
	return req
}

func (newrelic Client) execute(request *http.Request, payload []byte) (*http.Response, error) {
	newrelic.log.Info("Executing request", "Method", request.Method, "Endpoint", request.URL, "Payload", payload)
	response, err := newrelic.client.Do(request)
	if err != nil {
		return nil, err
	}

	return response, nil
}
