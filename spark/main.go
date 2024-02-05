package spark

import (
	"encoding/json"
	"fmt"
	"github.com/SneaksAndData/esd-services-api-client-go/shared/http"
	"golang.org/x/exp/slices"
	"log"
)

const codeRoot = "/ecco/dist"

var failedStages = []interface{}{
	"FAILED",
	"SCHEDULING_FAILED",
	"RETRIES_EXCEEDED",
	"SUBMISSION_FAILED",
	"STALE",
}
var successStages = []interface{}{"COMPLETED"}

type Service struct {
	httpClient *http.Client
	baseUrl    string
}

type JobRequest struct {
	inputs      []JobSocket
	outputs     []JobSocket
	args        interface{}
	tag         string
	parallelism int
}

type JobSocket struct {
	alias      string
	dataPath   string
	dataFormat string
}

type SubmissionConfiguration struct {
	RootPath          string      `json:"rootPath"`
	ProjectName       string      `json:"projectName"`
	Runnable          string      `json:"runnable"`
	SubmissionDetails interface{} `json:"submissionDetails"`
}

type submission struct {
	Id    string
	Stage interface{}
}

func (s Service) SubmitJob(request JobRequest, sparkJobName string) (string, error) {
	// TODO: check if submission already exists
	targetURL := fmt.Sprintf("%s/job/submit/%s", s.baseUrl, sparkJobName)
	return s.httpClient.MakeRequest("POST", targetURL, request)
}

func (s Service) CheckExistingSubmission(tag string) (string, error) {
	//fmt.Println(fmt.Sprintf("Looking for existing submission of %s", tag))
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

	//var jsonResponse map[string]interface{}
	var arr []string
	if err := json.Unmarshal([]byte(response), &arr); err != nil {
		return "", fmt.Errorf("error unmarshaling response: %w", err)
	}

	var runningSubmissions []submission
	//var jsonResponse map[string]interface{}
	//json.Unmarshal([]byte(response), &jsonResponse)
	for _, id := range arr {
		stage, err := s.GetLifecycleStage(id)
		if err != nil {
			return "", fmt.Errorf("error getting lifecycle stage for %s: %w", id, err)
		}
		if !slices.Contains(successStages, stage) && !slices.Contains(failedStages, stage) {
			log.Printf("Found a running submission of %s: %s", tag, id)
			runningSubmissions = append(runningSubmissions, submission{Id: id, Stage: stage})
		}
	}

	if len(runningSubmissions) == 0 {
		log.Println("None of found submissions are active")
		return "", nil
	}

	if len(runningSubmissions) > 1 {
		return "", fmt.Errorf("fatal: more than one submission of %s is running: %+v. Please review their status and restart/terminate the task accordingly", tag, runningSubmissions)
	}
	a := runningSubmissions[0]
	fmt.Println(a)
	run, err := json.Marshal(runningSubmissions[0])
	if err != nil {
		return "", fmt.Errorf("error marshaling running submission: %w", err)
	}
	str := string(run)

	return str, err

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
	json.Unmarshal([]byte(response), &jsonMap)

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
