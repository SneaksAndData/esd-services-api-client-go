// Package spark provides functionalities to manage and interact with Spark job submissions using Beast API.
package spark

import (
	"encoding/json"
	"fmt"
	"github.com/SneaksAndData/esd-services-api-client-go/shared/httpclient"
	"golang.org/x/exp/slices"
	"log"
	"net/http"
)

// Predefined lists of stages indicating job failure.
var failedStages = []string{
	"FAILED",
	"SCHEDULING_FAILED",
	"RETRIES_EXCEEDED",
	"SUBMISSION_FAILED",
	"STALE",
}

// Predefined lists of stages indicating job success.
var successStages = []string{"COMPLETED"}

// Service struct encapsulates the HTTP client and base URL for interacting with Beast.
type Service struct {
	httpClient *httpclient.Client
	baseURL    string
}

// JobParams defines the parameters for a Spark job submission.
type JobParams struct {
	ClientTag           string
	ExtraArguments      map[string]interface{}
	ProjectInputs       []JobSocket
	ProjectOutputs      []JobSocket
	ExpectedParallelism int
}

// JobSocket defines a data source or target for a Spark job.
type JobSocket struct {
	Alias      string
	DataPath   string
	DataFormat string
}

// SubmissionConfiguration represents the configuration of a Spark job submission.
type SubmissionConfiguration struct {
	RootPath          string      `json:"rootPath"`
	ProjectName       string      `json:"projectName"`
	Runnable          string      `json:"runnable"`
	SubmissionDetails interface{} `json:"submissionDetails"`
}

// submission represents the state of a Spark job submission.
type submission struct {
	ID    string
	Stage string
}

// RunJob submits a new Spark job or returns the ID of an existing job if one matches the ClientTag.
func (s Service) RunJob(request JobParams, sparkJobName string) (string, error) {
	submissionID, err := s.checkExistingSubmission(request.ClientTag)
	if err != nil {
		return "", fmt.Errorf("failed to check if submission exists: %w", err)
	}

	if submissionID != "" {
		return submissionID, nil
	}

	r, err := s.submitJob(request, sparkJobName)
	if err != nil {
		return "", fmt.Errorf("submit job failed with error: %w", err)
	}
	return r.ID, nil
}

// submitJob handles the actual submission of a Spark job.
func (s Service) submitJob(request JobParams, sparkJobName string) (submission, error) {
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

// checkExistingSubmission checks if there is an existing submission for the given tag.
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

// GetLifecycleStage retrieves the current lifecycle stage of a submission.
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

// GetRuntimeInfo retrieves the runtime information of a submission.
func (s Service) GetRuntimeInfo(id string) (string, error) {
	targetURL := fmt.Sprintf("%s/job/requests/%s", s.baseURL, id)
	return s.httpClient.MakeRequest(http.MethodGet, targetURL, nil)
}

// GetConfiguration checks if a configuration is deployed with the specified name.
func (s Service) GetConfiguration(name string) (SubmissionConfiguration, error) {
	targetURL := fmt.Sprintf("%s/job/deployed/%s", s.baseURL, name)
	response, err := s.httpClient.MakeRequest(http.MethodGet, targetURL, nil)
	if err != nil {
		fmt.Println(err)
	}
	var jsonMap SubmissionConfiguration
	if err := json.Unmarshal([]byte(response), &jsonMap); err != nil {
		return SubmissionConfiguration{}, fmt.Errorf("error unmarshaing response %w", err)
	}

	return jsonMap, nil
}

// GetLogs retrieves the logs of a submission.
func (s Service) GetLogs(id string) (string, error) {
	targetURL := fmt.Sprintf("%s/job/logs/%s", s.baseURL, id)
	return s.httpClient.MakeRequest(http.MethodGet, targetURL, nil)
}

// Config represents the configuration needed to create a new Service instance.
type Config struct {
	BaseURL      string
	GetTokenFunc func() (string, error)
	HTTPClient   *httpclient.Client
}

// New creates a new instance of the Service using the provided Config.
func New(c Config) (*Service, error) {
	s := &Service{
		httpClient: httpclient.NewClient(c.GetTokenFunc),
		baseURL:    c.BaseURL,
	}
	return s, nil
}
