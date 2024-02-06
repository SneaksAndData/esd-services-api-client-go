// Package algorithm provides functionalities to interact crystal
package algorithm

import (
	"fmt"

	"github.com/SneaksAndData/esd-services-api-client-go/shared/http"
)

// Service encapsulates the HTTP client and URLs needed to interact with the algorithm service.
type Service struct {
	httpClient   *http.Client
	schedulerURL string
	receiverURL  string
	apiVersion   string
}

// Payload defines the structure of the request body for creating algorithm runs.
type Payload struct {
	AlgorithmParameters interface{}
	AlgorithmName       string
	CustomConfiguration interface{}
	Tag                 string
}

// RetrieveRun fetches the results of a specific algorithm run identified by runID.
func (s Service) RetrieveRun(runID string, algorithmName string) (string, error) {
	targetURL := fmt.Sprintf("%s/algorithm/%s/results/%s/requests/%s", s.schedulerURL, s.apiVersion, algorithmName, runID)

	return s.httpClient.MakeRequest("GET", targetURL, nil)
}

// CreateRun initiates a new run of an algorithm with the given name, input parameters, and tag.
func (s Service) CreateRun(algorithmName string, input map[string]interface{}, tag string) (string, error) {
	targetURL := fmt.Sprintf("%s/algorithm/%s/run/%s", s.schedulerURL, s.apiVersion, algorithmName)
	body := Payload{
		AlgorithmParameters: input["algorithm_parameters"],
		AlgorithmName:       algorithmName,
		CustomConfiguration: input["custom_configuration"],
		Tag:                 tag,
	}
	return s.httpClient.MakeRequest("POST", targetURL, body)

}

// Config represents the configuration needed to create a new Service instance.
type Config struct {
	GetTokenFunc func() (string, error) // Function to retrieve authentication token
	HTTPClient   *http.Client           // HTTP client to be used by the Service
	SchedulerURL string                 // Base URL for the scheduler service
	APIVersion   string                 // API version to be used in requests
}

// New creates a new instance of the Service using the provided Config.
func New(c Config) (*Service, error) {
	s := &Service{
		httpClient:   http.NewClient(c.GetTokenFunc),
		schedulerURL: c.SchedulerURL,
		apiVersion:   c.APIVersion,
	}
	return s, nil
}
