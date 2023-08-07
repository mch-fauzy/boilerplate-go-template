package oauth

import (
	"github.com/jmoiron/sqlx"
)

type GrantType string

const (
	ClientCredentials GrantType = "client_credentials"
	Password          GrantType = "password"
)

type Token struct {
	config          Config
	tokenRepository TokenStore
}

func New(db *sqlx.DB, config Config) *Token {
	return &Token{
		config:          config,
		tokenRepository: NewTokenStore(db),
	}
}

type Config struct {
	Expiration  int64
	ClientScope []string
}

// Create is function to store NewToken into database
func (t *Token) Create(credential Credential) (*TokenResponse, error) {
	grant, err := NewGrant(t.tokenRepository, t.config).Create(credential)
	if err != nil {
		return &TokenResponse{}, err
	}

	return grant.toCreateTokenResponse(), nil
}

// ParseWithAccessToken is function to exchange valid token into token info
func (t *Token) ParseWithAccessToken(accessToken string) (OauthAccessToken, error) {
	return NewParser(t.tokenRepository).Parse(accessToken)
}

// ClientScopeAllowed is function that is used to limit the client
// set * to allowed all client example in confing, ex : ClientScope: ["*"] or keep it empty
// set clientId to limit scope, ex : ClientScope: ["client_web"]
func (t *Token) ClientScopeAllowed(clientID string) bool {
	if len(t.config.ClientScope) == 0 {
		return true
	}

	if len(t.config.ClientScope) > 0 {
		if t.config.ClientScope[0] == "*" {
			return true
		}
	}

	for _, c := range t.config.ClientScope {
		if c == clientID {
			return true
		}
	}

	return false
}
