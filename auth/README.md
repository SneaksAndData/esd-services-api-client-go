# Boxer Auth API Connector 

### Generate authentication token

```go
package main

import (
	"fmt"
	"github.com/SneaksAndData/esd-services-api-client-go/auth"
	"log"
)

func main() {
	// Configuration for the auth service
	config := auth.Config{
		TokenURL: "https://example.com",
		Provider: "azuread",
	}

	// Create a new instance of the auth service
	authService, err := auth.New(config)
	if err != nil {
		log.Fatalf("Failed to create auth service: %v", err)
	}
	// Retrieve token
	token, err := authService.GetBoxerToken()
	if err != nil {
		log.Fatalf("Failed to get token: %v", err)
	}

	fmt.Println("Token:", token)
}

```