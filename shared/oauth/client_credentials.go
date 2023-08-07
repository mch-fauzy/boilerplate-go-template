package oauth

import (
	"errors"
)

type ClientCredentialsAuth struct {
	tokenStore TokenStore
	config     Config
}

func (c *ClientCredentialsAuth) Create(credential Credential) (oauthAccessToken OauthAccessToken, err error) {
	client, err := c.tokenStore.resolveClientByClientID(credential.ClientID)
	if err != nil {
		return
	}

	if !client.VerifyClient(credential) {
		err = errors.New(ErrorInvalidClient)
		return
	}

	accessToken, err := generateAccessToken()
	if err != nil {
		err = errors.New(ErrorGenerateAccessToken)
		return
	}

	oauthAccessToken = new(OauthAccessToken).Generate(accessToken, credential.ClientID, nil, true, c.config)
	err = c.tokenStore.createAccessToken(oauthAccessToken)
	if err != nil {
		return
	}

	return
}
