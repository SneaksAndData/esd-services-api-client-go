package auth

import (
	"fmt"
)

func (s *service) GetBoxerToken() (string, error) {
	targetURL := fmt.Sprintf("%s/token/%s", s.TokenUrl, s.Provider)
	return s.HTTPClient.MakeRequest("GET", targetURL, nil)
}
