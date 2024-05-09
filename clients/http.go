package clients

import (
	"context"
	"net/http"
	"net/url"
	"time"
)

func getHTTPClient() customHTTPClient {
	client := http.Client{
		Timeout: 30 * time.Second,
	}
	return customHTTPClient{httpClient: client}
}

type customHTTPClient struct {
	httpClient http.Client
}

func (chc customHTTPClient) MakeGetRequest(
	urlPath string,
	params map[string]string,
	headers map[string]string,
) (*http.Response, error) {
	finalURL, err := addQueryParams(urlPath, params)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodGet,
		finalURL.String(),
		nil,
	)
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	res, err := chc.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func addQueryParams(
	urlPath string,
	params map[string]string,
) (*url.URL, error) {
	parsedURL, err := url.Parse(urlPath)
	if err != nil {
		return nil, err
	}
	queryParams := parsedURL.Query()
	for k, v := range params {
		queryParams.Add(k, v)
	}
	parsedURL.RawQuery = queryParams.Encode()
	return parsedURL, nil
}
