package spark

import (
	"encoding/json"
	"fmt"
	"github.com/SneaksAndData/esd-services-api-client-go/shared/http"
)

const codeRoot = "/ecco/dist"

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
	rootPath          string
	projectName       string
	runnable          string
	submissionDetails interface{}
}

func (s Service) SubmitJob(request JobRequest, sparkJobName string) (string, error) {
	targetURL := fmt.Sprintf("%s/job/submit/%s", s.baseUrl, sparkJobName)
	return s.httpClient.MakeRequest("POST", targetURL, request)
}

func (s Service) GetLifecycleStage(id string) (interface{}, error) {
	targetURL := fmt.Sprintf("%s/job/requests/%s", s.baseUrl, id)
	fmt.Println(targetURL)
	response, err := s.httpClient.MakeRequest("GET", targetURL, nil)
	if err != nil {
		fmt.Println(err)
	}
	var jsonMap map[string]interface{}
	json.Unmarshal([]byte(response), &jsonMap)
	return jsonMap["lifeCycleStage"], nil
}

func (s Service) GetRuntimeInfo(id string) (string, error) {
	targetURL := fmt.Sprintf("%s/job/requests/%s", id)
	return s.httpClient.MakeRequest("GET", targetURL, nil)
}

func (s Service) GetConfiguration(name string) (SubmissionConfiguration, error) {
	targetURL := fmt.Sprintf("%s/job/deployed/%s", name)
	response, err := s.httpClient.MakeRequest("GET", targetURL, nil)
	if err != nil {
		fmt.Println(err)
	}
	var jsonMap SubmissionConfiguration
	json.Unmarshal([]byte(response), &jsonMap)

	return jsonMap, nil
}

func (s Service) GetLogs(id string) (string, error) {
	targetURL := fmt.Sprintf("%s/job/logs/%s", id)
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
