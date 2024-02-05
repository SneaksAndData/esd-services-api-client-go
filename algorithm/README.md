# Crystal API Connector

### Retrieve run

```go
package main

import (
	"fmt"
	"github.com/SneaksAndData/esd-services-api-client-go/algorithm"
	"log"
)

func main() {
	// Configuration for the algorithm service
	var config = algorithm.Config{
		GetTokenFunc: getCachedBoxerToken,
		SchedulerUrl: "https://example.com",
		ApiVersion:   "v1.2",
	}
	// Create a new instance of the algorithm service
	algorithmService, err := algorithm.New(config)
	if err != nil {
		log.Fatalf("Failed to create algorithm service: %v", err)
	}

	// Retrieve run info
	run, err := algorithmService.RetrieveRun("run-id", "algorithm-name")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Run:", run)
}

```


### Submit run

```go
package main

import (
	"encoding/json"
	"fmt"
	"github.com/SneaksAndData/esd-services-api-client-go/algorithm"
	"io"
	"log"
	"os"
)

func main() {
	// Configuration for the algorithm service
	var config = algorithm.Config{
		GetTokenFunc: getCachedBoxerToken,
		SchedulerUrl: "https://example.com",
		ApiVersion:   "v1.2",
	}
	// Create a new instance of the algorithm service
	algorithmService, err := algorithm.New(config)
	if err != nil {
		log.Fatalf("Failed to create algorithm service: %v", err)
	}

	// Open payload json file
	jsonFile, err := os.Open("crystal-payload.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)

	var input map[string]interface{}
	if err := json.Unmarshal([]byte(byteValue), &input); err != nil {
		fmt.Errorf("error unmarshaling response: %w", err)
	}

	// Run algorithm
	response, err := algorithmService.CreateRun("algorithm-name", input, "tag")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Response:", response)
}

```