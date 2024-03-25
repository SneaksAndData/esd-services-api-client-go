// Package spark provides functionalities to interact with Beast
package spark

import (
	"encoding/json"
	"fmt"
	"github.com/SneaksAndData/esd-services-api-client-go/shared/httpclient"
	"golang.org/x/exp/slices"
	"log"
	"net/http"
	"strings"
)

var failedStages = []string{
	"FAILED",
	"SCHEDULING_FAILED",
	"RETRIES_EXCEEDED",
	"SUBMISSION_FAILED",
	"STALE",
}
var successStages = []string{"COMPLETED"}

// Service encapsulates the HTTP client and URL needed to interact with the Spark service.
type Service struct {
	httpClient *httpclient.Client
	baseURL    string
}

// JobParams defines the parameters for a Beast job
type JobParams struct {
	ClientTag           string                 `json:"clientTag"`
	ExtraArguments      map[string]interface{} `json:"extraArguments"`
	ProjectInputs       []JobSocket            `json:"projectInputs"`
	ProjectOutputs      []JobSocket            `json:"projectOutputs"`
	ExpectedParallelism *int                   `json:"expectedParallelism"`
}

// JobSocket defines the input/output data map
type JobSocket struct {
	// Alias: mapping key to be used by a consumer
	Alias string `json:"alias"`
	// DataPath: fully qualified path to actual data, i.e. abfss://..., s3a://... etc.
	DataPath string `json:"dataPath"`
	// DataFormat: data format, i.e. csv, json, delta etc.
	DataFormat string `json:"dataFormat"`
}

type JobRequest struct {
	Inputs              []JobSocket            `json:"inputs"`
	Outputs             []JobSocket            `json:"outputs"`
	ExtraArgs           map[string]interface{} `json:"extraArgs"`
	ClientTag           string                 `json:"clientTag"`
	ExpectedParallelism *int                   `json:"expectedParallelism"`
}

// SubmissionConfiguration defines the CRD used by Beast to run Spark jobs
type SubmissionConfiguration struct {
	RootPath          string            `json:"rootPath"`
	ProjectName       string            `json:"projectName"`
	Runnable          string            `json:"runnable"`
	SubmissionDetails SubmissionDetails `json:"submissionDetails"`
}

// SubmissionDetails defines job runtime details
type SubmissionDetails struct {
	Version                         string            `json:"version"`
	ExecutionGroup                  string            `json:"executionGroup"`
	ExpectedParallelism             int               `json:"expectedParallelism"`
	FlexibleDriver                  bool              `json:"flexibleDriver"`
	AdditionalDriverNodeTolerations map[string]string `json:"additionalDriverNodeTolerations"`
	MaxRuntimeHours                 int               `json:"maxRuntimeHours"`
	DebugMode                       RequestDebugMode  `json:"debugMode"`
	SubmissionMode                  string            `json:"submissionMode"`
	ExtendedCodeMount               bool              `json:"extendedCodeMount"`
	SubmissionJobTemplate           string            `json:"submissionJobTemplate"`
	ExecutorSpecTemplate            string            `json:"executorSpecTemplate"`
	DriverJobRetries                int               `json:"driverJobRetries"`
	DefaultArguments                map[string]string `json:"defaultArguments"`
	Inputs                          []JobSocket       `json:"inputs"`
	Outputs                         []JobSocket       `json:"outputs"`
	Overwrite                       bool              `json:"overwrite"`
}

// RequestDebugMode defines debug mode configuration
type RequestDebugMode struct {
	EventLogLocation string `json:"eventLogLocation"`
	MaxSizePerFile   string `json:"maxSizePerFile"`
}

type submission struct {
	ID    string
	Stage string
}

// RunJob runs a job through Beast
//
// Parameters:
//
// - request: Parameters for Beast Job body
//
// - sparkJobName: Name of the SparkJob to invoke
func (s Service) RunJob(request JobParams, sparkJobName string) (string, error) {
	submissionID, err := s.checkExistingSubmission(request.ClientTag)
	if err != nil {
		return "", fmt.Errorf("failed to check if submission exists: %w", err)
	}

	if submissionID != "" {
		return submissionID, nil
	}
	payload := JobRequest{
		Inputs:              request.ProjectInputs,
		Outputs:             request.ProjectOutputs,
		ExtraArgs:           request.ExtraArguments,
		ClientTag:           request.ClientTag,
		ExpectedParallelism: request.ExpectedParallelism,
	}

	r, err := s.submitJob(payload, sparkJobName)
	if err != nil {
		return "", fmt.Errorf("submit job failed with error: %w", err)
	}
	return r.ID, nil
}

func (s Service) submitJob(request JobRequest, sparkJobName string) (submission, error) {
	log.Printf("Submitting request: %+v", request)
	targetURL := fmt.Sprintf("%s/job/submit/%s", s.baseURL, sparkJobName)
	result, err := s.httpClient.MakeRequest(http.MethodPost, targetURL, request)
	if err != nil {
		return submission{}, fmt.Errorf("error making request to %s: %w", targetURL, err)
	}
	var sub submission
	if err := json.Unmarshal([]byte(result), &sub); err != nil {
		return submission{
			ID:    "",
			Stage: "",
		}, fmt.Errorf("error unmarshaling response: %w", err)
	}
	log.Printf("Beast has accepted the request, stage: %s, id: %s", sub.Stage, sub.ID)
	return sub, nil
}

func (s Service) checkExistingSubmission(tag string) (string, error) {
	log.Printf("Looking for existing submission of %s", tag)
	targetURL := fmt.Sprintf("%s/job/requests/tags/%s", s.baseURL, tag)
	response, err := s.httpClient.MakeRequest(http.MethodGet, targetURL, nil)
	if err != nil {
		return "", fmt.Errorf("error making request to %s: %w", targetURL, err)
	}
	if len(response) == 0 {
		log.Printf("No previous submissions found for %s", tag)
		return "", nil
	}

	var ids []string
	if err := json.Unmarshal([]byte(response), &ids); err != nil {
		return "", fmt.Errorf("error unmarshaling response: %w", err)
	}

	var runningSubmissions []submission
	for _, id := range ids {
		stage, err := s.GetLifecycleStage(id)
		if err != nil {
			return "", fmt.Errorf("error getting lifecycle stage for %s: %w", id, err)
		}
		if !slices.Contains(successStages, stage.(string)) && !slices.Contains(failedStages, stage.(string)) {
			log.Printf("Found a running submission of %s: %s", tag, id)
			runningSubmissions = append(runningSubmissions, submission{ID: id, Stage: stage.(string)})
		}
	}

	if len(runningSubmissions) == 0 {
		log.Println("None of found submissions are active")
		return "", nil
	}

	if len(runningSubmissions) > 1 {
		return "", fmt.Errorf("fatal: more than one submission of %s is running: %+v. Please review their status and restart/terminate the task accordingly", tag, runningSubmissions)
	}
	run, err := json.Marshal(runningSubmissions[0])
	if err != nil {
		return "", fmt.Errorf("error marshaling running submission: %w", err)
	}

	return string(run), err
}

// GetLifecycleStage returns the lifecycle stage for a given request
//
// Parameters:
//
// - id: A request identifier to read lifecycle stage info for
func (s Service) GetLifecycleStage(id string) (interface{}, error) {
	targetURL := fmt.Sprintf("%s/job/requests/%s", s.baseURL, id)
	response, err := s.httpClient.MakeRequest(http.MethodGet, targetURL, nil)
	if err != nil {
		return "", fmt.Errorf("error making request to %s: %w", targetURL, err)
	}
	var jsonMap map[string]interface{}
	if err := json.Unmarshal([]byte(response), &jsonMap); err != nil {
		return "", fmt.Errorf("error unmarshaling response: %w", err)
	}
	return jsonMap["lifeCycleStage"], nil
}

// GetRuntimeInfo returns runtime information for the given request
//
// Parameters:
//
// - id: A request identifier to read runtime info for
func (s Service) GetRuntimeInfo(id string) (string, error) {
	targetURL := fmt.Sprintf("%s/job/requests/%s", s.baseURL, id)
	response, err := s.httpClient.MakeRequest(http.MethodGet, targetURL, nil)
	if err != nil {
		return "", fmt.Errorf("error making request to %s: %w", targetURL, err)
	}
	return string(response), nil
}

// GetConfiguration returns a deployed SparkJob configuration
//
// Parameters:
//
// - name: Name of the configuration to find
func (s Service) GetConfiguration(name string) (SubmissionConfiguration, error) {
	targetURL := fmt.Sprintf("%s/job/deployed/%s", s.baseURL, name)
	response, err := s.httpClient.MakeRequest(http.MethodGet, targetURL, nil)
	if err != nil {
		fmt.Println(err)
	}
	var jsonMap SubmissionConfiguration
	if err := json.Unmarshal(response, &jsonMap); err != nil {
		return SubmissionConfiguration{}, fmt.Errorf("error unmarshalling response %w", err)
	}

	return jsonMap, nil
}

// GetLogs returns logs for a running or a completed submission
//
// Parameters:
//
// - id: Submission request identifier
func (s Service) GetLogs(id string) (string, error) {
	targetURL := fmt.Sprintf("%s/job/logs/%s", s.baseURL, id)
	response, err := s.httpClient.MakeRequest(http.MethodGet, targetURL, nil)
	if err != nil {
		return "", fmt.Errorf("error making request to %s: %w", targetURL, err)
	}
	var logsArray []string
	err = json.Unmarshal(response, &logsArray)
	if err != nil {
		return "", fmt.Errorf("error parsing API response: %v", err)
	}
	return strings.Join(logsArray, "\n"), nil
}

// Config represents the configuration needed to create a new spark Service instance.
type Config struct {
	BaseURL      string
	GetTokenFunc func() (string, error)
	HTTPClient   *httpclient.Client
}

// New creates a new instance of the spark Service using the provided Config.
func New(c Config) (*Service, error) {
	s := &Service{
		httpClient: httpclient.NewClient(c.GetTokenFunc),
		baseURL:    c.BaseURL,
	}
	return s, nil
}
