package client

import (
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	BaseURL *url.URL
	HTTP *http.Client
	API_KEY string
}


func NewClient(rawURL string, timeout time.Duration, apiKey string) (*Client, error) {
	url, err := url.Parse(rawURL)

	return &Client{
		BaseURL: url,
		HTTP: &http.Client{Timeout: timeout},
		API_KEY: apiKey,
	}, err
}
