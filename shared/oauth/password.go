package oauth

import (
	"errors"
)

type PasswordAuth struct {
	tokenStore TokenStore
	config     Config
}

func (c *PasswordAuth) Create(credential Credential) (oauthAccessToken OauthAccessToken, err error) {
	client, err := c.tokenStore.resolveClientByClientID(credential.ClientID)
	if err != nil {
		return
	}

	if !client.VerifyClient(credential) {
		err = errors.New(ErrorInvalidClient)
		return
	}

	user, err := c.tokenStore.resolveByTelephoneOrEmail(credential.Username)
	if err != nil {
		return
	}

	if !user.ValidCredential(credential) {
		err = errors.New(ErrorInvalidPassword)
		return
	}

	accessToken, err := generateAccessToken()
	if err != nil {
		err = errors.New(ErrorGenerateAccessToken)
		return
	}

	oauthAccessToken = new(OauthAccessToken).Generate(accessToken, credential.ClientID, &user.ID, false, c.config)

	err = c.tokenStore.createAccessToken(oauthAccessToken)
	if err != nil {
		return
	}

	return
}
