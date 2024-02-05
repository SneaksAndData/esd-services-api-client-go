package spark

import (
	"encoding/json"
	"fmt"
	"github.com/SneaksAndData/esd-services-api-client-go/shared/http"
	"golang.org/x/exp/slices"
	"log"
)

var failedStages = []string{
	"FAILED",
	"SCHEDULING_FAILED",
	"RETRIES_EXCEEDED",
	"SUBMISSION_FAILED",
	"STALE",
}
var successStages = []string{"COMPLETED"}

type Service struct {
	httpClient *http.Client
	baseUrl    string
}

type JobParams struct {
	ClientTag           string
	ExtraArguments      map[string]interface{}
	ProjectInputs       []JobSocket
	ProjectOutputs      []JobSocket
	ExpectedParallelism int
}

type JobSocket struct {
	Alias      string
	DataPath   string
	DataFormat string
}

type SubmissionConfiguration struct {
	RootPath          string      `json:"rootPath"`
	ProjectName       string      `json:"projectName"`
	Runnable          string      `json:"runnable"`
	SubmissionDetails interface{} `json:"submissionDetails"`
}

type submission struct {
	Id    string
	Stage string
}

func (s Service) RunJob(request JobParams, sparkJobName string) (string, error) {
	submissionId, err := s.checkExistingSubmission(request.ClientTag)
	if err != nil {
		return "", fmt.Errorf("failed to check if submission exists: %w", err)
	}

	if submissionId != "" {
		return submissionId, nil
	}

	r, err := s.submitJob(request, sparkJobName)
	if err != nil {
		return "", fmt.Errorf("submit job failed with error: %w", err)
	}
	return r.Id, nil
}

func (s Service) submitJob(request JobParams, sparkJobName string) (submission, error) {
	log.Printf("Submitting request: %+v", request)
	targetURL := fmt.Sprintf("%s/job/submit/%s", s.baseUrl, sparkJobName)
	result, err := s.httpClient.MakeRequest("POST", targetURL, request)
	if err != nil {
		return submission{}, fmt.Errorf("error making request to %s: %w", targetURL, err)
	}
	var sub submission
	if err := json.Unmarshal([]byte(result), &sub); err != nil {
		return submission{
			Id:    "",
			Stage: "",
		}, fmt.Errorf("error unmarshaling response: %w", err)
	}
	log.Printf("Beast has accepted the request, stage: %s, id: %s", sub.Stage, sub.Id)
	return sub, nil
}

func (s Service) checkExistingSubmission(tag string) (string, error) {
	log.Printf("Looking for existing submission of %s", tag)
	targetURL := fmt.Sprintf("%s/job/requests/tags/%s", s.baseUrl, tag)
	response, err := s.httpClient.MakeRequest("GET", targetURL, nil)
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
			runningSubmissions = append(runningSubmissions, submission{Id: id, Stage: stage.(string)})
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

func (s Service) GetLifecycleStage(id string) (interface{}, error) {
	targetURL := fmt.Sprintf("%s/job/requests/%s", s.baseUrl, id)
	response, err := s.httpClient.MakeRequest("GET", targetURL, nil)
	if err != nil {
		return "", fmt.Errorf("error making request to %s: %w", targetURL, err)
	}
	var jsonMap map[string]interface{}
	if err := json.Unmarshal([]byte(response), &jsonMap); err != nil {
		return "", fmt.Errorf("error unmarshaling response: %w", err)
	}
	return jsonMap["lifeCycleStage"], nil
}

func (s Service) GetRuntimeInfo(id string) (string, error) {
	targetURL := fmt.Sprintf("%s/job/requests/%s", s.baseUrl, id)
	return s.httpClient.MakeRequest("GET", targetURL, nil)
}

func (s Service) GetConfiguration(name string) (SubmissionConfiguration, error) {
	targetURL := fmt.Sprintf("%s/job/deployed/%s", s.baseUrl, name)
	response, err := s.httpClient.MakeRequest("GET", targetURL, nil)
	if err != nil {
		fmt.Println(err)
	}
	var jsonMap SubmissionConfiguration
	if err := json.Unmarshal([]byte(response), &jsonMap); err != nil {
		return SubmissionConfiguration{}, fmt.Errorf("error unmarshaing response %w", err)
	}

	return jsonMap, nil
}

func (s Service) GetLogs(id string) (string, error) {
	targetURL := fmt.Sprintf("%s/job/logs/%s", s.baseUrl, id)
	return s.httpClient.MakeRequest("GET", targetURL, nil)
}

type Config struct {
	BaseUrl      string
	GetTokenFunc func() (string, error)
	HTTPClient   *http.Client
}

func New(c Config) (*Service, error) {
	s := &Service{
		httpClient: http.NewClient(c.GetTokenFunc),
		baseUrl:    c.BaseUrl,
	}
	return s, nil
}
