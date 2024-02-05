package algorithm

import (
	"fmt"

	"github.com/SneaksAndData/esd-services-api-client-go/shared/http"
)

type Service struct {
	httpClient   *http.Client
	schedulerUrl string
	receiverUrl  string
	apiVersion   string
}

type Payload struct {
	AlgorithmParameters interface{}
	AlgorithmName       string
	CustomConfiguration interface{}
	Tag                 string
}

func (s Service) RetrieveRun(runId string, algorithmName string) (string, error) {
	targetURL := fmt.Sprintf("%s/algorithm/%s/results/%s/requests/%s", s.schedulerUrl, s.apiVersion, algorithmName, runId)

	return s.httpClient.MakeRequest("GET", targetURL, nil)
}

func (s Service) CreateRun(algorithmName string, input map[string]interface{}, tag string) (string, error) {
	targetURL := fmt.Sprintf("%s/algorithm/%s/run/%s", s.schedulerUrl, s.apiVersion, algorithmName)
	body := Payload{
		AlgorithmParameters: input["algorithm_parameters"],
		AlgorithmName:       algorithmName,
		CustomConfiguration: input["custom_configuration"],
		Tag:                 tag,
	}
	return s.httpClient.MakeRequest("POST", targetURL, body)

}

type Config struct {
	GetTokenFunc func() (string, error)
	HTTPClient   *http.Client
	SchedulerUrl string
	ApiVersion   string
}

func New(c Config) (*Service, error) {
	s := &Service{
		httpClient:   http.NewClient(c.GetTokenFunc),
		schedulerUrl: c.SchedulerUrl,
		apiVersion:   c.ApiVersion,
	}
	return s, nil
}
