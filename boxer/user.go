package boxer

import (
	"fmt"
	"github.com/SneaksAndData/esd-services-api-client-go/shared/http"
)

func (c connector) AddUser(user string, provider string, token string) (string, error) {
	targetURL := fmt.Sprintf("%s/claim/%s/%s", c.claimUrl, provider, user)
	client := http.NewClient(token)

	return client.MakeRequest("POST", targetURL, nil)
}

func (c connector) RemoveUser(user string, provider string, token string) (string, error) {
	targetURL := fmt.Sprintf("%s/claim/%s/%s", c.claimUrl, provider, user)
	client := http.NewClient(token)

	return client.MakeRequest("DELETE", targetURL, nil)
}
