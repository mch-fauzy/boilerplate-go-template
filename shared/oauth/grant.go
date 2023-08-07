package oauth

type AuthorizationMethod interface {
	Create(credential Credential) (OauthAccessToken, error)
}

type Grant struct {
	TokenStore TokenStore
	Config     Config
}

func NewGrant(tokenStore TokenStore, config Config) *Grant {
	return &Grant{
		TokenStore: tokenStore,
		Config:     config,
	}
}

func (g *Grant) Create(credential Credential) (OauthAccessToken, error) {
	authMap := make(map[GrantType]AuthorizationMethod)
	authMap[ClientCredentials] = &ClientCredentialsAuth{tokenStore: g.TokenStore, config: g.Config}
	authMap[Password] = &PasswordAuth{tokenStore: g.TokenStore, config: g.Config}

	return authMap[credential.GrantType].Create(credential)
}
