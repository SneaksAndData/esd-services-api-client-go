package claim

import (
	"github.com/SneaksAndData/esd-services-api-client-go/shared/http"
)

type Service interface {
	Claim
	User
}

// Claim is an interface for managing claims.
type Claim interface {
	GetClaim(user string, provider string) (string, error)
	AddClaim(user string, provider string, claims []string) (string, error)
	RemoveClaim(user string, provider string, claims []string) (string, error)
}

type User interface {
	AddUser(user string, provider string) (string, error)
	RemoveUser(user string, provider string) (string, error)
}

type service struct {
	Config
}

type Config struct {
	ClaimUrl     string
	GetTokenFunc func() (string, error)
	HTTPClient   *http.Client
}

func New(c Config) (Service, error) {
	s := &service{
		Config: c,
	}
	httpClient := http.NewClient(c.GetTokenFunc)
	s.HTTPClient = httpClient

	return s, nil
}
