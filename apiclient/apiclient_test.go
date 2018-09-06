package apiclient

import (
	"errors"
	"net/http"
	"testing"
)

func TestPublish(t *testing.T) {
	client := &Client{
		HTTPClient:  &ClientMock{},
		EndpointURL: "https://task.appranix.net/issues.json",
	}
	client.Publish([]byte("sample data"))
	//	return response, nil
}

func TestPublishError(t *testing.T) {

	client := &Client{
		HTTPClient:  &ClientMock{},
		EndpointURL: "",
	}
	client.Publish([]byte("sample data"))
	//	return response, nil
}

func TestNewHTTPClient(t *testing.T) {

	client := NewHTTPClient()
	if client == nil {
		t.Errorf("Expected http client but received nothing")
	}
}

type ClientMock struct {
}

func (c *ClientMock) Do(req *http.Request) (*http.Response, error) {
	if req.Host == "" {
		return nil, errors.New("URL cannot be empty")
	}
	return &http.Response{}, nil
}
