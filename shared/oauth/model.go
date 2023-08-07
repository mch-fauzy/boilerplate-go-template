package oauth

import (
	"golang.org/x/crypto/bcrypt"
	"strconv"
	"time"

	"github.com/guregu/null"
)

type TokenType string

const (
	Bearer TokenType = "Bearer"
)

var scope = struct {
	User string
}{
	"user",
}

// Credential is
type Credential struct {
	GrantType    GrantType
	ClientID     string
	ClientSecret string
	Username     string
	Password     string
}

type OauthAccessToken struct {
	AccessToken string      `json:"accessToken" db:"access_token"`
	ClientID    string      `json:"clientId" db:"client_id"`
	UserID      null.String `json:"userId" db:"user_id"`
	Expires     time.Time   `json:"expires" db:"expires"`
	Scope       null.String `json:"scope" db:"scope"`
}

func (o *OauthAccessToken) Generate(accessToken string, clientID string, userID *int, withScope bool, config Config) OauthAccessToken {
	if userID != nil {
		o.UserID = null.StringFrom(strconv.Itoa(*userID))
	}

	if withScope {
		o.Scope = null.StringFrom(scope.User)
	}

	o.ClientID = clientID
	o.AccessToken = accessToken
	o.Expires = time.Now().Add(time.Second * time.Duration(config.Expiration))

	return *o
}

func (o *OauthAccessToken) VerifyExpireIn() bool {
	now := time.Now()
	if now.After(o.Expires) {
		return false
	}

	return true
}

func (o *OauthAccessToken) VerifyUserLoggedIn() bool {
	if o.UserID.Valid && !o.Scope.Valid {
		return true
	}
	return false
}

func (o *OauthAccessToken) toCreateTokenResponse() *TokenResponse {
	return &TokenResponse{
		AccessToken: o.AccessToken,
		ExpiresIn:   o.Expires,
		TokenType:   string(Bearer),
		Scope:       scope.User,
	}
}

type OauthClient struct {
	ClientID     string `json:"clientId" db:"client_id"`
	ClientSecret string `json:"clientSecret" db:"client_secret"`
	RedirectURI  string `json:"redirectUri" db:"redirect_uri"`
	GrantTypes   string `json:"grantTypes" db:"grant_types"`
}

func (o *OauthClient) VerifyClient(credential Credential) bool {
	if o.ClientID != credential.ClientID {
		return false
	}

	if o.ClientSecret != credential.ClientSecret {
		return false
	}

	return true
}

type TokenResponse struct {
	AccessToken string
	ExpiresIn   time.Time
	TokenType   string
	Scope       string
}

type User struct {
	ID       int    `json:"id" db:"id"`
	Username string `json:"username" db:"username"`
	Password string `json:"password" db:"password"`
}

func (u *User) ValidCredential(credential Credential) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(credential.Password))
	if err != nil {
		return false
	}

	return true
}
