package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client wraps the standard http.Client and adds automatic token retrieval for making authenticated requests.
type Client struct {
	httpClient *http.Client
	getToken   func() (string, error) // Function to get or refresh the token
}

// NewClient creates a new Client instance with a specified function for token retrieval.
func NewClient(getTokenFunc func() (string, error)) *Client {
	return &Client{
		httpClient: &http.Client{Timeout: 30 * time.Second},
		getToken:   getTokenFunc,
	}
}

// MakeRequest creates and executes an HTTP request using the given method, URL, and payload.
// It automatically handles token retrieval and will retry the request once if the token is expired.
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
		}
		return "", err // Return original error if not an auth failure

	}

	return responseBody, nil // Return successful response from initial request
}

// prepareRequest creates an *http.Request object with the given method, URL, token, and payload.
func (c *Client) prepareRequest(method, url string, payload interface{}, token string) (*http.Request, error) {
	var body io.Reader

	if payload != nil {
		jsonPayload, err := json.Marshal(payload)
		if err != nil {
			return nil, fmt.Errorf("error marshaling payload: %w", err)
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

// doRequest executes the given *http.Request and returns the response body as a string.
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
