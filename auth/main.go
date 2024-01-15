package auth

import (
	"fmt"
	"github.com/SneaksAndData/esd-services-api-client-go/shared/http"
)

type TokenProvider struct {
	httpClient *http.Client
	TokenUrl   string
	Provider   string
	//GetBoxerToken func() (string, error)
}

// Config represents the configuration inputs for creating a new auth service.
type Config struct {
	TokenUrl string // tokenUrl is the URL used to retrieve the Boxer internal token e.g. http://boxer.test.sneaksanddata.com.
	Env      string
	Provider string
}

// New creates a new Connector instance with the provided configuration.
// Returns a TokenProvider
func New(c Config) (*TokenProvider, error) {
	t := &TokenProvider{httpClient: nil}
	t.TokenUrl = c.TokenUrl
	t.Provider = c.Provider

	switch c.Provider {
	case "azuread":
		t.httpClient = http.NewClient(getAzureDefaultToken)
	default:
		return nil, fmt.Errorf("unsupported token provider: %s", c.Provider)
	}

	return t, nil
}

func (t *TokenProvider) GetBoxerToken() (string, error) {
	targetURL := fmt.Sprintf("%s/token/%s", t.TokenUrl, t.Provider)
	return t.httpClient.MakeRequest("GET", targetURL, nil)
}
