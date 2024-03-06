# Beast Connector

### Run Job

```go
package main

import (
	"fmt"
	"github.com/SneaksAndData/esd-services-api-client-go/spark"
	"log"
)

func main() {
	// Configuration for the spark service
	configSpark := spark.Config{
		BaseURL:     "example.com",
		GetTokenFunc: getToken,
	}

	// Create a new instance of the spark service
	sparkService, err := spark.New(configSpark)
	if err != nil {
		log.Fatalf("Failed to create spark service: %v", err)
	}

	var extraArguments = map[string]interface{}{
		"destination": "'dest'",
		// other arguments
	}

	var inputJobSocket = spark.JobSocket{
		Alias:      "input",
		DataPath:   "abfss://...",
		DataFormat: "delta",
	}
	
	var outputJobSocket = spark.JobSocket{
		Alias:      "target",
		DataPath:   "abfss://...",
		DataFormat: "delta",
	}
	
	parameters := spark.JobParams{
		ClientTag:           "",
		ExtraArguments:      extraArguments,
		ProjectInputs:       inputJobSocket,
		ProjectOutputs:      outputJobSocket,
		ExpectedParallelism: 1,
    }
	// Run job
	response, err := sparkService.RunJob(parameters, "spark-job-name")
	if err != nil {fmt.Print(err)}
	
	fmt.Println("Response:", response)
}
```

### Get Lifecycle Stage
```go
package main

import (
	"fmt"
	"github.com/SneaksAndData/esd-services-api-client-go/spark"
	"log"
)

func main() {
	// Configuration for the spark service
	configSpark := spark.Config{
		BaseURL:      "example.com",
		GetTokenFunc: getToken,
	}

	// Create a new instance of the spark service
	sparkService, err := spark.New(configSpark)
	if err != nil {
		log.Fatalf("Failed to create spark service: %v", err)
	}

	stage, err := sparkService.GetLifecycleStage("job-id")
	if err != nil {
		fmt.Print(err)
	}

	fmt.Println(stage)
}
```

### Get Runtime Info
```go
package main

import (
	"fmt"
	"github.com/SneaksAndData/esd-services-api-client-go/spark"
	"log"
)

func main() {
	// Configuration for the spark service
	configSpark := spark.Config{
		BaseURL:      "example.com",
		GetTokenFunc: getToken,
	}

	// Create a new instance of the spark service
	sparkService, err := spark.New(configSpark)
	if err != nil {
		log.Fatalf("Failed to create spark service: %v", err)
	}

	info, err := sparkService.GetRuntimeInfo("job-id")
	if err != nil {
		fmt.Print(err)
	}

	fmt.Println(info)
}
```

### Get Configuration
```go
package main

import (
	"fmt"
	"github.com/SneaksAndData/esd-services-api-client-go/spark"
	"log"
)

func main() {
	// Configuration for the spark service
	configSpark := spark.Config{
		BaseURL:      "example.com",
		GetTokenFunc: getToken,
	}

	// Create a new instance of the spark service
	sparkService, err := spark.New(configSpark)
	if err != nil {
		log.Fatalf("Failed to create spark service: %v", err)
	}

	conf, err := sparkService.GetConfiguration("some-configuration-name")
	if err != nil {
		fmt.Print(err)
	}

	fmt.Println(conf)
}

```

### Get Logs
```go
package main

import (
	"fmt"
	"github.com/SneaksAndData/esd-services-api-client-go/spark"
	"log"
)

func main() {
	// Configuration for the spark service
	configSpark := spark.Config{
		BaseURL:      "example.com",
		GetTokenFunc: getToken,
	}

	// Create a new instance of the spark service
	sparkService, err := spark.New(configSpark)
	if err != nil {
		log.Fatalf("Failed to create spark service: %v", err)
	}

	logs, err := sparkService.GetLogs("some-id")
	if err != nil {
		fmt.Print(err)
	}

	fmt.Println(logs)
}

```