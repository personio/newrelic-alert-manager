package internal

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/go-logr/logr"
	"io/ioutil"
	"net/http"
	"strconv"
)

type NewrelicClient interface {
	Get(path string) (*http.Response, error)
	GetJson(path string) (*http.Response, error)
	PostJson(path string, payload []byte) (*http.Response, error)
	PutJson(path string, payload []byte) (*http.Response, error)
	Delete(path string) (*http.Response, error)
}

type newrelicClient struct {
	client *http.Client
	log    logr.InfoLogger

	url      string
	adminKey string
}

func NewNewrelicClient(log logr.Logger, url string, adminKey string) NewrelicClient {
	return newrelicClient{
		client:   &http.Client{},
		log:      log,
		url:      url,
		adminKey: adminKey,
	}
}

func (newrelic newrelicClient) Get(path string) (*http.Response, error) {
	request := newrelic.newRequest("GET", path, nil)

	return newrelic.executeWithStatusCheck(request)
}

func (newrelic newrelicClient) GetJson(path string) (*http.Response, error) {
	request := newrelic.newJsonRequest("Get", path, nil)

	return newrelic.executeWithStatusCheck(request)
}

func (newrelic newrelicClient) PostJson(path string, payload []byte) (*http.Response, error) {
	request := newrelic.newJsonRequest("POST", path, payload)

	return newrelic.executeWithStatusCheck(request)
}

func (newrelic newrelicClient) PutJson(path string, payload []byte) (*http.Response, error) {
	request := newrelic.newJsonRequest("PUT", path, payload)

	return newrelic.executeWithStatusCheck(request)
}

func (newrelic newrelicClient) Delete(path string) (*http.Response, error) {
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

func (newrelic *newrelicClient) newRequest(method string, path string, body []byte) *http.Request {
	var req *http.Request
	if body == nil {
		req = newRequest(method, newrelic.url, path)
	} else {
		req = newRequestWithBody(method, newrelic.url, path, body)
	}

	req.Header.Add("X-Api-Key", newrelic.adminKey)

	return req
}

func (newrelic *newrelicClient) newJsonRequest(method string, path string, body []byte) *http.Request {
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

func (newrelic newrelicClient) execute(request *http.Request) (*http.Response, error) {
	newrelic.log.Info("Executing request", "Method", request.Method, "Endpoint", request.URL, "Payload", request.Body)
	response, err := newrelic.client.Do(request)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (newrelic newrelicClient) executeWithStatusCheck(request *http.Request) (*http.Response, error) {
	response, err := newrelic.execute(request)
	if err != nil {
		return nil, err
	}

	if response != nil {
		newrelic.log.Info(strconv.Itoa(response.StatusCode))
	}

	if response != nil && response.StatusCode >= 300 {
		responseContent, _ := ioutil.ReadAll(response.Body)
		return nil, errors.New(string(responseContent))
	}

	return response, nil
}
