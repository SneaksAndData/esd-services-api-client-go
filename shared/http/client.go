package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client wraps the http.Client and provides methods for making HTTP requests.
type Client struct {
	HTTPClient *http.Client
	Token      string
}

// NewClient creates and returns a new Client instance.
func NewClient(token string) *Client {
	return &Client{
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		Token: token,
	}
}

// MakeRequest makes an HTTP request with the given method, URL, and payload.
// It automatically adds necessary headers including the authorization token.
func (c *Client) MakeRequest(method, url string, payload interface{}) (string, error) {
	var body io.Reader

	if payload != nil {
		jsonPayload, err := json.Marshal(payload)
		if err != nil {
			return "", err
		}
		body = bytes.NewBuffer(jsonPayload)
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return "", err
	}

	req.Header.Add("Authorization", "Bearer "+c.Token)
	req.Header.Add("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", err
	}
	defer func() {
		errClose := resp.Body.Close()
		if err == nil && errClose != nil {
			err = errClose
		}
	}()

	if resp.StatusCode == http.StatusUnauthorized {
		return "", fmt.Errorf("authorization failed, please run 'snd login' to refresh your token")
	} else if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		return "", fmt.Errorf("HTTP request failed with status code: %d", resp.StatusCode)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(respBody), nil
}
