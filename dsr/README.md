# Data Subject Request (DSR) Service

###  Search for an email
```go
package main

import (
	"fmt"
	"github.com/SneaksAndData/esd-services-api-client-go/dsr"
	"log"
)

func main() {
	// Configuration for the DSR service
	configDsr := dsr.Config{
		GetTokenFunc: getToken,
		DsrBaseUrl:   "https://dsr.example.com",
	}

	// Create a new instance of the DSR service
	dsrService, err := dsr.New(configDsr)
	if err != nil {
		log.Fatalf("Failed to create DSR service: %v", err)
	}

	response, err := dsrService.GetDSRRequest("some-email")
	if err != nil {
		fmt.Print(err)
	}

	fmt.Println(response)
}
```