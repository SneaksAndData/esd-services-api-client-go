package auth

import (
	"fmt"
	"github.com/SneaksAndData/esd-services-api-client-go/shared/http"
)

type Service struct {
	httpClient *http.Client
	tokenUrl   string
	provider   string
}

func (s *Service) GetBoxerToken() (string, error) {
	targetURL := fmt.Sprintf("%s/token/%s", s.tokenUrl, s.provider)
	return s.httpClient.MakeRequest("GET", targetURL, nil)
}

// Config represents the configuration inputs for creating a new auth service.
type Config struct {
	TokenUrl string // tokenUrl is the URL used to retrieve the Boxer internal token e.g. http://boxer.test.sneaksanddata.com.
	Env      string
	Provider string
}

// New creates a new Connector instance with the provided configuration.
func New(c Config) (*Service, error) {
	s := &Service{httpClient: nil}
	s.tokenUrl = c.TokenUrl
	s.provider = c.Provider

	switch c.Provider {
	case "azuread":
		s.httpClient = http.NewClient(getAzureDefaultToken)
	default:
		return nil, fmt.Errorf("unsupported token provider: %s", c.Provider)
	}

	return s, nil
}
