package auth

import (
	"os"
)

func getKubernetesToken() (string, error) {
	tokenFile, err := os.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/token")
	if err != nil {
		return "", err
	}
	return string(tokenFile), nil
}
