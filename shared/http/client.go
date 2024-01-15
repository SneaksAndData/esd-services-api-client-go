package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	httpClient *http.Client
	getToken   func() (string, error) // Function to get or refresh the token
}

func NewClient(getTokenFunc func() (string, error)) *Client {
	return &Client{
		httpClient: &http.Client{Timeout: 30 * time.Second},
		getToken:   getTokenFunc,
	}
}

// MakeRequest makes an HTTP request with the given method, URL, and payload.
func (c *Client) MakeRequest(method, url string, payload interface{}) (string, error) {
	token, err := c.getToken()
	if err != nil {
		return "", fmt.Errorf("failed to get token: %v", err)
	}

	request, err := c.prepareRequest(method, url, payload, token)
	if err != nil {
		return "", err
	}

	responseBody, err := c.doRequest(request)
	if err != nil {
		if err.Error() == "authorization failed" {
			refreshedToken, err := c.getToken()
			if err != nil {
				return "", err
			}

			// Retry the requests with the new token
			retryRequest, retryErr := c.prepareRequest(method, url, payload, refreshedToken)
			if retryErr != nil {
				return "", retryErr
			}
			responseBody, err = c.doRequest(retryRequest)
			if err != nil {
				return "", err // Return error if retry also fails
			}
			return responseBody, nil // Return successful response from retry
		} else {
			return "", err // Return original error if not an auth failure
		}
	}

	return responseBody, nil // Return successful response from initial request
}

func (c *Client) prepareRequest(method, url string, payload interface{}, token string) (*http.Request, error) {
	var body io.Reader

	if payload != nil {
		jsonPayload, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}
		body = bytes.NewBuffer(jsonPayload)
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Content-Type", "application/json")

	return req, nil
}

func (c *Client) doRequest(req *http.Request) (string, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return "", fmt.Errorf("authorization failed")
	} else if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		return "", fmt.Errorf("HTTP request failed with status code: %d", resp.StatusCode)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(respBody), nil
}
