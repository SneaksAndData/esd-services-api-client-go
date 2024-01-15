package auth

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

func (s *service) getAzureDefaultToken() (string, error) {
	cred, err := getAzureCredentials()
	if err != nil {
		return "", err
	}
	token, err := cred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{"https://management.core.windows.net/.default"}})
	if err != nil {
		return "", err
	}
	return token.Token, nil
}

func getAzureCredentials() (*azidentity.DefaultAzureCredential, error) {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, err
	}
	return cred, nil
}
