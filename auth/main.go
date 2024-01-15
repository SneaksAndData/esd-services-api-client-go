package auth

import (
	"fmt"
	"github.com/SneaksAndData/esd-services-api-client-go/shared/http"
)

type Service interface {
	TokenProvider
}

// TokenProvider is an interface for managing authentication tokens.
type TokenProvider interface {
	GetBoxerToken() (string, error)
}

// service is an implementation of the Service interface.
type service struct {
	Config
}

// Config represents the configuration inputs for creating a new auth service.
type Config struct {
	TokenUrl   string // tokenUrl is the URL used to retrieve the Boxer internal token e.g. http://boxer.test.sneaksanddata.com.
	Env        string
	Provider   string
	HTTPClient *http.Client
}

// New creates a new Connector instance with the provided configuration.
// It returns an implementation of the Service interface.
func New(c Config) (Service, error) {
	s := &service{
		Config: c,
	}

	switch c.Provider {
	case "azuread":
		s.HTTPClient = http.NewClient(s.getAzureDefaultToken)
	default:
		return nil, fmt.Errorf("unsupported token provider: %s", c.Provider)
	}

	return s, nil
}
