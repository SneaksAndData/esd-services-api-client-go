package beast

type Connector interface {
	RunJob(request JobRequest, jobName string, clientTag string) (string, error)
	GetRuntimeInfo(id string, token string) (string, error)
	GetConfiguration() (string, error)
	GetLogs() (string, error)
}

type connector struct {
	url string
}

func NewConnector(url string) Connector {
	return &connector{
		url: url,
	}
}
