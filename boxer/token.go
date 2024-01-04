package boxer

import (
	"fmt"
	"github.com/SneaksAndData/esd-services-api-client-go/shared/http"
)

type ExternalToken struct {
	GetToken func() (string, error)
	Provider string
	Retry    bool
}

func (c connector) GetToken() (string, error) {
	targetURL := fmt.Sprintf("%s/token/%s", c.tokenUrl, c.auth.Provider)

	token, err := c.auth.GetToken()
	if err != nil {
		return "", err
	}
	client := http.NewClient(token)

	return client.MakeRequest("GET", targetURL, nil)
}
