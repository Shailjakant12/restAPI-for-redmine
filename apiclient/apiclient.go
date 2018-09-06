package apiclient

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"net/http"
)

type Client struct {
	HTTPClient  ClientInterface
	EndpointURL string
}

type ClientInterface interface {
	Do(req *http.Request) (*http.Response, error)
}

func NewHTTPClient() *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	return &http.Client{Transport: tr}
}

func (apiclient *Client) Publish(data []byte) (*http.Response, error) {
	request, _ := http.NewRequest("POST", apiclient.EndpointURL, bytes.NewBuffer(data))
	request.Header.Set("authorization", "Basic YzBmNGU2NjNlNjU5M2Q1NTdhOTU2MDc2N2ZjMGI5ODBkMDhmOGM4NA==")
	request.Header.Set("content-type", "application/json")
	request.Header.Set("accept", "application/json")
	response, err := apiclient.HTTPClient.Do(request)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	}
	return response, err
}
