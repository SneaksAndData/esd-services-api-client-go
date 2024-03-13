// Package httpclient provides a convenient wrapper around the standard http.Client, introducing functionality for automatic token retrieval and management.
// This package simplifies the process of making authenticated HTTP requests by encapsulating the logic for token acquisition, token refresh, and request retries upon authentication failures.
//
// The primary component of this package is the Client struct, which extends http.Client
// with additional capabilities to automatically handle authentication tokens for requests.
// Clients can specify a custom function for token retrieval, which is invoked as needed
// to obtain or refresh tokens before making requests.
package httpclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Client wraps the standard httpclient.Client and adds automatic token retrieval for making authenticated requests.
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
func (c *Client) MakeRequest(method, url string, payload interface{}) ([]byte, error) {
	token, err := c.getToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %v", err)
	}

	request, err := c.prepareRequest(method, url, payload, strings.TrimSpace(token))
	if err != nil {
		return nil, err
	}

	responseBody, err := c.doRequest(request)
	if err != nil {
		if err.Error() == "authorization failed" {
			refreshedToken, err := c.getToken()
			if err != nil {
				return nil, err
			}

			// Retry the requests with the new token
			retryRequest, retryErr := c.prepareRequest(method, url, payload, refreshedToken)
			if retryErr != nil {
				return nil, retryErr
			}
			responseBody, err = c.doRequest(retryRequest)
			if err != nil {
				return nil, err // Return error if retry also fails
			}
			return responseBody, nil // Return successful response from retry
		}
		return nil, err // Return original error if not an auth failure

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
func (c *Client) doRequest(req *http.Request) ([]byte, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, fmt.Errorf("authorization failed")
	} else if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		return nil, fmt.Errorf("HTTP request failed with status code: %d", resp.StatusCode)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return respBody, nil
}
