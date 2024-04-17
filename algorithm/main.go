// Package algorithm provides functionalities to interact with Crystal
package algorithm

import (
	"encoding/json"
	"fmt"
	"github.com/SneaksAndData/esd-services-api-client-go/shared/httpclient"
	"github.com/go-playground/validator/v10"
	"log"
	"net/http"
)

// Service encapsulates the HTTP client and URLs needed to interact with the algorithm service.
type Service struct {
	httpClient   *httpclient.Client
	schedulerURL string
	apiVersion   string
}

// Payload defines the structure of the request body for creating algorithm runs.
type Payload struct {
	AlgorithmParameters map[string]interface{} `validate:"required"`
	AlgorithmName       string
	CustomConfiguration CustomConfiguration
	Tag                 string
}

type CustomConfiguration struct {
	ImageRepository      *string              `json:"imageRepository"`
	ImageTag             *string              `json:"imageTag"`
	DeadlineSeconds      *int                 `json:"deadlineSeconds"`
	MaximumRetries       *int                 `json:"maximumRetries"`
	Env                  []ConfigurationEntry `json:"env"`
	Secrets              []string             `json:"secrets"`
	Args                 []ConfigurationEntry `json:"args"`
	CpuLimit             *string              `json:"cpuLimit"`
	MemoryLimit          *string              `json:"memoryLimit"`
	Workgroup            *string              `json:"workgroup"`
	AdditionalWorkgroups map[string]string    `json:"additionalWorkgroups"`
	Version              *string              `json:"version"`
	MonitoringParameters []string             `json:"monitoringParameters"`
	CustomResources      map[string]string    `json:"customResources"`
	SpeculativeAttempts  *int                 `json:"speculativeAttempts"`
}

type ConfigurationValueType string

const (
	PLAIN              ConfigurationValueType = "PLAIN"
	RELATIVE_REFERENCE ConfigurationValueType = "RELATIVE_REFERENCE"
)

type ConfigurationEntry struct {
	Name      string                  `json:"name"`
	Value     string                  `json:"value"`
	ValueType *ConfigurationValueType `json:"valueFrom"`
}

type PayloadResponse struct {
	RequestID  string `json:"requestId"`
	PayloadUri string `json:"payloadUri"`
}

// RetrieveRun fetches the results of a specific algorithm run identified by runID.
func (s Service) RetrieveRun(runID string, algorithmName string) (string, error) {
	targetURL := fmt.Sprintf("%s/algorithm/%s/results/%s/requests/%s", s.schedulerURL, s.apiVersion, algorithmName, runID)
	response, err := s.httpClient.MakeRequest(http.MethodGet, targetURL, nil)
	if err != nil {
		return "", fmt.Errorf("error making request to %s: %w", targetURL, err)
	}
	return string(response), nil
}

// RetrievePayloadUri fetches the payload URI of a specific algorithm run identified by runID.
func (s Service) RetrievePayloadUri(runID string, algorithmName string) (*PayloadResponse, error) {
	targetURL := fmt.Sprintf("%s/algorithm/%s/payload/%s/requests/%s", s.schedulerURL, s.apiVersion, algorithmName, runID)
	response, err := s.httpClient.MakeRequest(http.MethodGet, targetURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error making request to %s: %w", targetURL, err)
	}

	var payloadResponse PayloadResponse
	err = json.Unmarshal(response, &payloadResponse)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %w", err)
	}
	return &payloadResponse, nil
}

func (s Service) CreateRun(algorithmName string, input Payload, tag string) (string, error) {
	if err := validator.New().Struct(input); err != nil {
		log.Fatalf("Validation failed: %v\n", err)
	}

	targetURL := fmt.Sprintf("%s/algorithm/%s/run/%s", s.schedulerURL, s.apiVersion, algorithmName)

	input.AlgorithmName = algorithmName
	input.Tag = tag
	response, err := s.httpClient.MakeRequest(http.MethodPost, targetURL, input)
	if err != nil {
		return "", fmt.Errorf("error making request to %s: %w", targetURL, err)
	}
	return string(response), nil
}

// CancelRun cancels an ongoing algorithm run
func (s Service) CancelRun(algorithmName string, requestId string, initiator string, reason string) (string, error) {
	targetURL := fmt.Sprintf("%s/algorithm/%s/cancel/%s/requests/%s", s.schedulerURL, s.apiVersion, algorithmName, requestId)
	payload := make(map[string]string)
	payload["initiator"] = initiator
	payload["reason"] = reason
	response, err := s.httpClient.MakeRequest(http.MethodPost, targetURL, payload)
	if err != nil {
		return "", fmt.Errorf("error making request to %s: %w", targetURL, err)
	}
	return string(response), nil
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
