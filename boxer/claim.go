package boxer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/SneaksAndData/esd-services-api-client-go/shared/http"
	"strings"
)

type claimPayload struct {
	// Fields need to be public so that json package can see it
	Operation string            `json:"operation"`
	Claims    map[string]string `json:"claims"`
}

func (c connector) GetClaim(user string, provider string, token string) (string, error) {
	targetURL := fmt.Sprintf("%s/claim/%s/%s", c.claimUrl, provider, user)
	client := http.NewClient(token)

	return client.MakeRequest("GET", targetURL, nil)
}

func (c connector) AddClaim(user string, provider string, claims []string, token string) (string, error) {
	targetURL := fmt.Sprintf("%s/claim/%s/%s", c.claimUrl, provider, user)

	payload, err := json.Marshal(preparePayload(claims, "Insert"))
	if err != nil {
		return "", err
	}
	client := http.NewClient(token)

	return client.MakeRequest("PATCH", targetURL, bytes.NewBuffer(payload))
}

func (c connector) RemoveClaim(user string, provider string, claims []string, token string) (string, error) {
	targetURL := fmt.Sprintf("%s/claim/%s/%s", c.claimUrl, provider, user)

	payload, err := json.Marshal(preparePayload(claims, "Delete"))
	if err != nil {
		return "", err
	}
	client := http.NewClient(token)

	return client.MakeRequest("PATCH", targetURL, bytes.NewBuffer(payload))
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
