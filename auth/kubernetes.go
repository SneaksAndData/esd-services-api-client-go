package auth

import (
	"fmt"
	"github.com/SneaksAndData/esd-services-api-client-go/shared/file"
	"os"
)

func getKubernetesToken() (string, error) {
	tokenFilePath := "/var/run/secrets/kubernetes.io/serviceaccount/token"
	if !file.FileExists(tokenFilePath) {
		return "", fmt.Errorf("could not find token file at %s", tokenFilePath)
	}
	token, err := os.ReadFile(tokenFilePath)
	if err != nil {
		return "", fmt.Errorf("error reading token file: %w", err)
	}
	return string(token), nil
}
