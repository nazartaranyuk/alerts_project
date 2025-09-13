package client

import (
	"encoding/json"
	"io"
	"log"
	"nazartaraniuk/alertsProject/internal/domain"
	"net/http"
	"time"
)

type Client struct {
	BaseURL string
	HTTP    *http.Client
	API_KEY string
}

func (c *Client) SetAuthorization() {

}

func (c *Client) GetCurrentAlerts() ([]domain.RegionAlarmInfo, error) {
	resp, err := c.HTTP.Get(c.BaseURL + "/alerts")

	if err != nil {
		log.Fatalf("Cannot send GET due to error: %v", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Fatalf("Cannot read body: %v", err)
	}

	var response []domain.RegionAlarmInfo
	err = json.Unmarshal(body, &response)

	if err != nil {
		log.Fatalf("Cannot unmarshal response due to error: %v", err)
	}

	return response, nil

}

type AuthRoundTripper struct {
	ApiKey string
	Next   http.RoundTripper
}

func (art *AuthRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	request := req.Clone(req.Context())
	request.Header.Add("Authorization", art.ApiKey)
	return art.Next.RoundTrip(request)
}

func NewClient(rawURL string, timeout time.Duration, apiKey string) *Client {

	return &Client{
		BaseURL: rawURL,
		HTTP: &http.Client{
			Timeout: timeout,
			Transport: &AuthRoundTripper{
				ApiKey: apiKey,
				Next:   http.DefaultTransport,
			},
		},
		API_KEY: apiKey,
	}
}
