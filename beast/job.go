package beast

import (
	"fmt"
	"github.com/SneaksAndData/esd-services-api-client-go/shared/http"
)

type JobRequest struct {
	Inputs              []JobSocket
	Outputs             []JobSocket
	Overwrite           bool
	ClientTag           string
	ExpectedParallelism int
	ExtraArgs           map[string]interface{}
}

//type JobParams struct {
//	ClientTag           string
//	ProjectInputs       []JobSocket
//	ProjectOutputs      []JobSocket
//	OverwriteOutputs    bool
//	ExpectedParallelism int
//	ExtraArguments      map[string]interface{}
//}

type JobSocket struct {
	alias      string
	dataPath   string
	dataFormat string
}

func (c connector) RunJob(request JobRequest, jobName string, clientTag string) (string, error) {
	//targetURL := fmt.Sprintf("%s/job/submit/%s", c.url, jobName)
	//payload := JobRequest{
	//	Inputs:              request.Inputs,
	//	Outputs:             request.Outputs,
	//	Overwrite:           request.Overwrite,
	//	ClientTag:           request.ClientTag,
	//	ExpectedParallelism: request.ExpectedParallelism,
	//	ExtraArgs: map[string]interface{}
	//}
	return "", nil
}
func (c connector) GetRuntimeInfo(id string, token string) (string, error) {
	targetURL := fmt.Sprintf("%s/job/requests/%s", c.url, id)
	fmt.Println(targetURL)
	client := http.NewClient(token)
	return client.MakeRequest("GET", targetURL, nil)
}

func (c connector) GetLogs() (string, error) {
	return "", nil
}

func (c connector) GetConfiguration() (string, error) {
	return "", nil
}

func existingSubmission(tag string, baseUrl string, token string) (string, error) {
	fmt.Println(fmt.Sprintf("Looking for existing submission of %s", tag))
	targetURL := fmt.Sprintf("%s/job/requests/tags/%s", baseUrl, tag)
	client := http.NewClient(token)
	submission, err := client.MakeRequest("GET", targetURL, nil)
	if err != nil {
		return "", err
	}
	if submission == "" {
		fmt.Println(fmt.Sprintf("No previous submission found for %s", tag))
		return "", nil
	}

	//runningSubmission := []string
	//for i, s := range runningSubmission {
	//	checkSubmissionURL := fmt.Sprintf("%s/job/requests/tags/%s", baseUrl, tag)
	//}
	return "", nil
}
