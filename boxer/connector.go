// Package boxer provides an interface for interacting with Boxer authorization service.
// It defines several interfaces to work with users, tokens, and claims.
package boxer

// Connector is an interface that groups the Token, Claim, and User interfaces.
type Connector interface {
	Token
	Claim
	User
}

// User is an interface for managing users in Boxer.
type User interface {
	AddUser(user string, provider string, token string) (string, error)
	RemoveUser(user string, provider string, token string) (string, error)
}

// Token is an interface for managing authentication tokens.
type Token interface {
	GetToken() (string, error)
}

// Claim is an interface for managing claims.
type Claim interface {
	GetClaim(user string, provider string, token string) (string, error)
	AddClaim(user string, provider string, claims []string, token string) (string, error)
	RemoveClaim(user string, provider string, claims []string, token string) (string, error)
}

// connector is an implementation of the Connector interface.
type connector struct {
	tokenUrl string        // tokenUrl is the URL used to retrieve the Boxer internal token e.g. http://boxer.test.sneaksanddata.com.
	claimUrl string        // claimUrl is the URL used to manage claims and users e.g. http://boxer-claim.test.sneaksanddata.com.
	auth     ExternalToken // auth represents external token-based authentication, only used for generating the token
}

// Input represents the configuration inputs for creating a new Connector.
type Input struct {
	TokenUrl string
	ClaimUrl string
	Auth     ExternalToken
}

// NewConnector creates a new Connector instance with the provided configuration.
// It returns an implementation of the Connector interface.
func NewConnector(i Input) Connector {
	return &connector{
		tokenUrl: i.TokenUrl,
		claimUrl: i.ClaimUrl,
		auth:     i.Auth,
	}
}
