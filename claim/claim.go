package claim

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

type claimPayload struct {
	// Fields need to be public so that json package can see it
	Operation string            `json:"operation"`
	Claims    map[string]string `json:"claims"`
}

func (s service) GetClaim(user string, provider string) (string, error) {
	targetURL := fmt.Sprintf("%s/claim/%s/%s", s.ClaimUrl, provider, user)

	return s.HTTPClient.MakeRequest("GET", targetURL, nil)
}

func (s service) AddClaim(user string, provider string, claims []string) (string, error) {
	targetURL := fmt.Sprintf("%s/claim/%s/%s", s.ClaimUrl, provider, user)

	payload, err := json.Marshal(preparePayload(claims, "Insert"))
	if err != nil {
		return "", err
	}
	return s.HTTPClient.MakeRequest("PATCH", targetURL, bytes.NewBuffer(payload))
}

func (s service) RemoveClaim(user string, provider string, claims []string) (string, error) {
	targetURL := fmt.Sprintf("%s/claim/%s/%s", s.ClaimUrl, provider, user)

	payload, err := json.Marshal(preparePayload(claims, "Delete"))
	if err != nil {
		return "", err
	}

	return s.HTTPClient.MakeRequest("PATCH", targetURL, bytes.NewBuffer(payload))
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
