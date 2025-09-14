package client

import (
	"encoding/json"
	"io"
	"nazartaraniuk/alertsProject/internal/app/client/trippers"
	"nazartaraniuk/alertsProject/internal/domain"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

type Client struct {
	BaseURL string
	HTTP    *http.Client
	APIKey  string
}

func NewClient(rawURL string, timeout time.Duration, apiKey string) *Client {
	return &Client{
		BaseURL: rawURL,
		HTTP: &http.Client{
			Timeout: timeout,
			Transport: &trippers.AuthRoundTripper{
				APIKey: apiKey,
				Next:   http.DefaultTransport,
			},
		},
		APIKey: apiKey,
	}
}

func (c *Client) GetCurrentAlerts() ([]domain.RegionAlarmInfo, error) {
	resp, err := c.HTTP.Get(c.BaseURL + "/alerts")
	if err != nil {
		logrus.Fatalf("Cannot send GET due to error: %v", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.Fatalf("Cannot read body: %v", err)
	}

	var response []domain.RegionAlarmInfo
	err = json.Unmarshal(body, &response)
	if err != nil {
		logrus.Printf("Cannot unmarshal response due to error: %v", err)
		return nil, err
	}

	return response, nil
}
