// Package claim provides functionalities to manage user claims (interacts with Boxer).
package claim

import (
	"fmt"
	"github.com/SneaksAndData/esd-services-api-client-go/shared/httpclient"
	"net/http"
	"strings"
)

// Service represents a client for the claim management service.
type Service struct {
	httpClient *httpclient.Client
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

	return s.httpClient.MakeRequest(http.MethodGet, targetURL, nil)
}

// AddClaim adds claims for a user under a specific provider.
func (s Service) AddClaim(user string, provider string, claims []string) (string, error) {
	targetURL := fmt.Sprintf("%s/claim/%s/%s", s.claimURL, provider, user)

	return s.httpClient.MakeRequest(http.MethodPatch, targetURL, preparePayload(claims, "Insert"))
}

// RemoveClaim removes claims for a user under a specific provider.
func (s Service) RemoveClaim(user string, provider string, claims []string) (string, error) {
	targetURL := fmt.Sprintf("%s/claim/%s/%s", s.claimURL, provider, user)

	return s.httpClient.MakeRequest(http.MethodPatch, targetURL, preparePayload(claims, "Delete"))
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

	return s.httpClient.MakeRequest(http.MethodPost, targetURL, nil)
}

// RemoveUser deletes a user under a specific provider.
func (s Service) RemoveUser(user string, provider string) (string, error) {
	targetURL := fmt.Sprintf("%s/claim/%s/%s", s.claimURL, provider, user)

	return s.httpClient.MakeRequest(http.MethodDelete, targetURL, nil)
}

// Config holds the configuration needed to initialize a new Service instance.
type Config struct {
	ClaimURL     string
	GetTokenFunc func() (string, error)
	HTTPClient   *httpclient.Client
}

// New initializes a new instance of the Service using the provided Config.
func New(c Config) (*Service, error) {
	s := &Service{
		httpClient: httpclient.NewClient(c.GetTokenFunc),
		claimURL:   c.ClaimURL,
	}
	return s, nil
}
