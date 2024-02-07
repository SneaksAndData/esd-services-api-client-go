// Package algorithm provides functionalities to interact crystal
package algorithm

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/SneaksAndData/esd-services-api-client-go/shared/httpclient"
)

// Service encapsulates the HTTP client and URLs needed to interact with the algorithm service.
type Service struct {
	httpClient   *httpclient.Client
	schedulerURL string
	receiverURL  string
	apiVersion   string
}

// Payload defines the structure of the request body for creating algorithm runs.
type Payload struct {
	AlgorithmParameters map[string]interface{}
	AlgorithmName       string
	CustomConfiguration CustomConfiguration
	Tag                 string
}

type CustomConfiguration struct {
	ImageRepository      string
	ImageTag             string
	DeadlineSeconds      int
	MaximumRetries       int
	Env                  []ConfigurationEntry
	Secrets              []string
	Args                 []ConfigurationEntry
	CpuLimit             string
	MemoryLimit          string
	Workgroup            string
	AdditionalWorkgroups map[string]string
	Version              string
	MonitoringParameters []string
	CustomResources      map[string]string
	SpeculativeAttempts  int
}

type ConfigurationValueType string

const (
	PLAIN              ConfigurationValueType = "PLAIN"
	RELATIVE_REFERENCE ConfigurationValueType = "RELATIVE_REFERENCE"
)

type ConfigurationEntry struct {
	name      string
	value     string
	valueType *ConfigurationValueType
}

// RetrieveRun fetches the results of a specific algorithm run identified by runID.
func (s Service) RetrieveRun(runID string, algorithmName string) (string, error) {
	targetURL := fmt.Sprintf("%s/algorithm/%s/results/%s/requests/%s", s.schedulerURL, s.apiVersion, algorithmName, runID)
	return s.httpClient.MakeRequest(http.MethodGet, targetURL, nil)
}

// CreateRun initiates a new run of an algorithm with the given name, input parameters, and tag.
func (s Service) CreateRun(algorithmName string, input map[string]interface{}, tag string) (string, error) {
	targetURL := fmt.Sprintf("%s/algorithm/%s/run/%s", s.schedulerURL, s.apiVersion, algorithmName)
	// Handle CustomConfiguration
	conf, exists := input["custom_configuration"]
	if !exists {
		return "", fmt.Errorf("custom_configuration not provided")
	}
	customConfigMap, ok := conf.(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("custom_configuration is not of the expected type")
	}
	customConfigJSON, err := json.Marshal(customConfigMap)
	if err != nil {
		return "", fmt.Errorf("failed to marshal custom_configuration: %w", err)
	}
	var customConfig CustomConfiguration
	if err := json.Unmarshal(customConfigJSON, &customConfig); err != nil {
		return "", fmt.Errorf("failed to unmarshal custom_configuration into CustomConfiguration: %w", err)
	}

	body := Payload{
		AlgorithmParameters: input["algorithm_parameters"].(map[string]interface{}),
		AlgorithmName:       algorithmName,
		CustomConfiguration: customConfig,
		Tag:                 tag,
	}
	return s.httpClient.MakeRequest(http.MethodPost, targetURL, body)

}

// Config represents the configuration needed to create a new Service instance.
type Config struct {
	GetTokenFunc func() (string, error) // Function to retrieve authentication token
	HTTPClient   *httpclient.Client     // HTTP client to be used by the Service
	SchedulerURL string                 // Base URL for the scheduler service
	APIVersion   string                 // API version to be used in requests
}

// New creates a new instance of the Service using the provided Config.
func New(c Config) (*Service, error) {
	s := &Service{
		httpClient:   httpclient.NewClient(c.GetTokenFunc),
		schedulerURL: c.SchedulerURL,
		apiVersion:   c.APIVersion,
	}
	return s, nil
}
