package claim

import (
	"fmt"
)

func (s service) AddUser(user string, provider string) (string, error) {
	targetURL := fmt.Sprintf("%s/claim/%s/%s", s.ClaimUrl, provider, user)

	return s.HTTPClient.MakeRequest("POST", targetURL, nil)
}

func (s service) RemoveUser(user string, provider string) (string, error) {
	targetURL := fmt.Sprintf("%s/claim/%s/%s", s.ClaimUrl, provider, user)

	return s.HTTPClient.MakeRequest("DELETE", targetURL, nil)
}
