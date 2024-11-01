package dsr

import (
	"fmt"
	"github.com/SneaksAndData/esd-services-api-client-go/shared/httpclient"
	"net/http"
)

// Service encapsulates the HTTP client and URLs needed to interact with the DSR API.
type Service struct {
	httpClient *httpclient.Client
	dsrBaseUrl string
}

func (s Service) GetDSRRequest(email string) (string, error) {
	targetURL := s.dsrBaseUrl + "/dsr/" + email
	response, err := s.httpClient.MakeRequest(http.MethodGet, targetURL, nil)
	if err != nil {
		return "", fmt.Errorf("error making request to %s: %w", targetURL, err)
	}
	return string(response), nil
}

// Config represents the configuration needed to create a new Service instance.
type Config struct {
	GetTokenFunc func() (string, error) // Function to retrieve authentication token
	HTTPClient   *httpclient.Client     // HTTP client to be used by the Service
	DsrBaseUrl   string                 // Base URL for the DSR API service
}

// New creates a new instance of the Service using the provided Config.
func New(c Config) (*Service, error) {
	s := &Service{
		httpClient: httpclient.NewClient(c.GetTokenFunc),
		dsrBaseUrl: c.DsrBaseUrl,
	}
	return s, nil
}
