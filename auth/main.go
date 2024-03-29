// Package auth provides authentication services, including token retrieval for various providers (interacts with Boxer).
package auth

import (
	"fmt"
	"github.com/SneaksAndData/esd-services-api-client-go/shared/httpclient"
	"net/http"
	"strings"
)

// Service encapsulates the HTTP client, token URL, and provider for retrieving authentication tokens.
type Service struct {
	httpClient *httpclient.Client
	tokenURL   string
	provider   string
}

// GetBoxerToken retrieves an authentication token from the configured provider.
func (s *Service) GetBoxerToken() (string, error) {
	targetURL := fmt.Sprintf("%s/token/%s", s.tokenURL, s.provider)
	response, err := s.httpClient.MakeRequest(http.MethodGet, targetURL, nil)
	if err != nil {
		return "", fmt.Errorf("error making request to %s: %w", targetURL, err)
	}
	return string(response), nil
}

// Config represents the configuration inputs for creating a new auth service.
type Config struct {
	TokenURL string // tokenURL is the URL used to retrieve the Boxer internal token e.g. http://boxer.test.sneaksanddata.com.
	Env      string
	Provider string
}

// New initializes a new Service instance using the provided Config.
// It sets up the Service with an appropriate HTTP client based on the specified provider.
func New(c Config) (*Service, error) {
	s := &Service{httpClient: nil}
	s.tokenURL = c.TokenURL
	s.provider = c.Provider

	switch {
	case c.Provider == "azuread":
		s.httpClient = httpclient.NewClient(getAzureDefaultToken)
	case strings.HasPrefix(c.Provider, "k8s"):
		s.provider = strings.TrimPrefix(c.Provider, "k8s-")
		s.httpClient = httpclient.NewClient(getKubernetesToken)
	default:
		return nil, fmt.Errorf("unsupported token provider: %s", c.Provider)
	}

	return s, nil
}
