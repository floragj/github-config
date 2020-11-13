package internal

import (
	"fmt"
	"net/http"
	"net/url"
)

//go:generate faux --interface HTTPClient --output fakes/http_client.go
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type APIClient struct {
	ServerURL string
	AuthToken string
	client    HTTPClient
}

func NewAPIClient(serverURL string, httpClient HTTPClient) APIClient {
	return APIClient{ServerURL: serverURL,
		client: httpClient}
}

func (c *APIClient) Get(path string, params ...string) (*http.Response, error) {

	uri, err := url.Parse(c.ServerURL)
	if err != nil {
		return nil, fmt.Errorf("could not parse server URL: %s", err)
	}

	uri.Path = path
	if len(params) > 0 {
		uri.RawQuery = params[0]
		for i := range params {
			if i != 0 {
				uri.RawQuery = fmt.Sprintf("%s&%s", uri.RawQuery, params[i])
			}
		}
	}

	request, _ := http.NewRequest("GET", uri.String(), nil)

	response, err := c.client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("client couldn't make HTTP request: %s", err)
	}
	return response, nil
}
