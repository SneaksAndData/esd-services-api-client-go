package crystal

type Connector interface {
	CreateRun(algorithmName string, parameters interface{}, customConfig AlgorithmConfiguration, tag string, token string) (string, error)
	RetrieveRun(runId string, algorithmName string, token string) (string, error)
	SubmitResult(runId string, algorithmName string, cause string, message string, sasUri string, token string) (string, error)
}

type connector struct {
	schedulerUrl string
	receiverUrl  string
	apiVersion   string
}

func NewConnector(schedulerUrl, receiverUrl, apiVersion string) Connector {
	return &connector{
		schedulerUrl: schedulerUrl,
		receiverUrl:  receiverUrl,
		apiVersion:   apiVersion,
	}
}
