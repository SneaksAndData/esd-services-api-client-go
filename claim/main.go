package claim

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/SneaksAndData/esd-services-api-client-go/shared/http"
	"strings"
)

type Service struct {
	httpClient *http.Client
	claimUrl   string
}

type claimPayload struct {
	// Fields need to be public so that json package can see it
	Operation string            `json:"operation"`
	Claims    map[string]string `json:"claims"`
}

func (s Service) GetClaim(user string, provider string) (string, error) {
	targetURL := fmt.Sprintf("%s/claim/%s/%s", s.claimUrl, provider, user)

	return s.httpClient.MakeRequest("GET", targetURL, nil)
}

func (s Service) AddClaim(user string, provider string, claims []string) (string, error) {
	targetURL := fmt.Sprintf("%s/claim/%s/%s", s.claimUrl, provider, user)

	payload, err := json.Marshal(preparePayload(claims, "Insert"))
	if err != nil {
		return "", err
	}
	return s.httpClient.MakeRequest("PATCH", targetURL, bytes.NewBuffer(payload))
}

func (s Service) RemoveClaim(user string, provider string, claims []string) (string, error) {
	targetURL := fmt.Sprintf("%s/claim/%s/%s", s.claimUrl, provider, user)

	payload, err := json.Marshal(preparePayload(claims, "Delete"))
	if err != nil {
		return "", err
	}

	return s.httpClient.MakeRequest("PATCH", targetURL, bytes.NewBuffer(payload))
}

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

func (s Service) AddUser(user string, provider string) (string, error)
{
	targetURL := fmt.Sprintf("%s/claim/%s/%s", s.claimUrl, provider, user)

	return s.httpClient.MakeRequest("POST", targetURL, nil)
}

func (s Service) RemoveUser(user string, provider string) (string, error) {
	targetURL := fmt.Sprintf("%s/claim/%s/%s", s.claimUrl, provider, user)

	return s.httpClient.MakeRequest("DELETE", targetURL, nil)
}

type Config struct {
	ClaimUrl     string
	GetTokenFunc func() (string, error)
	HTTPClient   *http.Client
}

func New(c Config) (*Service, error) {
	s := &Service{
		httpClient: http.NewClient(c.GetTokenFunc),
		claimUrl:   c.ClaimUrl,
	}
	return s, nil
}
