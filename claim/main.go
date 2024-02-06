// Package claim provides functionalities to manage user claims (interacts with Boxer).
package claim

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/SneaksAndData/esd-services-api-client-go/shared/http"
	"strings"
)

// Service represents a client for the claim management service.
type Service struct {
	httpClient *http.Client
	claimURL   string
}

// claimPayload defines the structure for sending claim operations via HTTP requests.
type claimPayload struct {
	// Fields need to be public so that json package can see it
	Operation string            `json:"operation"`
	Claims    map[string]string `json:"claims"`
}

// GetClaim retrieves the claims for a given user and provider.
func (s Service) GetClaim(user string, provider string) (string, error) {
	targetURL := fmt.Sprintf("%s/claim/%s/%s", s.claimURL, provider, user)

	return s.httpClient.MakeRequest("GET", targetURL, nil)
}

// AddClaim adds claims for a user under a specific provider.
func (s Service) AddClaim(user string, provider string, claims []string) (string, error) {
	targetURL := fmt.Sprintf("%s/claim/%s/%s", s.claimURL, provider, user)

	payload, err := json.Marshal(preparePayload(claims, "Insert"))
	if err != nil {
		return "", fmt.Errorf("error marshaling payload: %w", err)
	}
	return s.httpClient.MakeRequest("PATCH", targetURL, bytes.NewBuffer(payload))
}

// RemoveClaim removes claims for a user under a specific provider.
func (s Service) RemoveClaim(user string, provider string, claims []string) (string, error) {
	targetURL := fmt.Sprintf("%s/claim/%s/%s", s.claimURL, provider, user)

	payload, err := json.Marshal(preparePayload(claims, "Delete"))
	if err != nil {
		return "", fmt.Errorf("error marshaling payload: %w", err)
	}

	return s.httpClient.MakeRequest("PATCH", targetURL, bytes.NewBuffer(payload))
}

// preparePayload prepares the payload for claim operations.
func preparePayload(claims []string, operation string) claimPayload {
	claimsMap := make(map[string]string)
	for _, s := range claims {
		c := strings.Split(s, ":")
		claimsMap[c[0]] = c[1]

	}
	return claimPayload{
		Operation: operation,
		Claims:    claimsMap,
	}
}

// AddUser creates a new user under a specific provider.
func (s Service) AddUser(user string, provider string) (string, error) {
	targetURL := fmt.Sprintf("%s/claim/%s/%s", s.claimURL, provider, user)

	return s.httpClient.MakeRequest("POST", targetURL, nil)
}

// RemoveUser deletes a user under a specific provider.
func (s Service) RemoveUser(user string, provider string) (string, error) {
	targetURL := fmt.Sprintf("%s/claim/%s/%s", s.claimURL, provider, user)

	return s.httpClient.MakeRequest("DELETE", targetURL, nil)
}

// Config holds the configuration needed to initialize a new Service instance.
type Config struct {
	ClaimURL     string
	GetTokenFunc func() (string, error)
	HTTPClient   *http.Client
}

// New initializes a new instance of the Service using the provided Config.
func New(c Config) (*Service, error) {
	s := &Service{
		httpClient: http.NewClient(c.GetTokenFunc),
		claimURL:   c.ClaimURL,
	}
	return s, nil
}
