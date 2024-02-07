# Boxer Claim API Connector 

### Get claims

```go
package main

import (
	"fmt"
	"github.com/SneaksAndData/esd-services-api-client-go/claim"
	"log"
)

func main() {
	// Configuration for the claim service
	config := claim.Config{
		ClaimURL:     "example.com",
		GetTokenFunc: getToken,
	}

	// Create a new instance of the claim service
	claimService, err := claim.New(config)
	if err != nil {
		log.Fatalf("Failed to create claim service: %v", err)
	}

	// Retrieve user claims
	response, err := claimService.GetClaim("user@ecco.com", "provider")
	if err != nil {
		log.Fatalf("Failed to get Boxer token: %v", err)
	}

	fmt.Println("Response:", response)
}

```