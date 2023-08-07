package oauth

import (
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
)

type TokenStore struct {
	db *sqlx.DB
}

const (
	queryInsertAccessToken = `INSERT INTO oauth_access_tokens (
			access_token,
			client_id,
			user_id,
			expires,
			scope
		) VALUES (
			:access_token,
			:client_id,
			:user_id,
			:expires,
			:scope
		)`

	querySelectAccessToken = `SELECT 
			access_token,
			client_id,
			user_id,
			expires,
			scope
		FROM
			oauth_access_tokens`

	querySelectClients = `SELECT
			client_id,
			client_secret,
			redirect_uri,
			grant_types
		FROM 
			oauth_clients`

	querySelectUser = `
			SELECT
				id,
				username,
				password
			FROM
				user`
)

func NewTokenStore(db *sqlx.DB) TokenStore {
	return TokenStore{
		db: db,
	}
}

func (a *TokenStore) createAccessToken(accessToken OauthAccessToken) error {
	stmt, err := a.db.PrepareNamed(queryInsertAccessToken)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(accessToken)
	if err != nil {
		return err
	}

	return nil
}

func (a *TokenStore) resolveAccessTokenByAccessToken(accessToken string) (oauthAccessToken OauthAccessToken, err error) {
	err = a.db.Get(&oauthAccessToken, querySelectAccessToken+" WHERE access_token = ?", accessToken)
	switch {
	case err == sql.ErrNoRows:
		err = errors.New(ErrorClientNotFound)
		return
	case err != nil:
		return
	}

	return
}

func (a *TokenStore) resolveAllClients(db *sqlx.DB) ([]OauthClient, error) {
	var clients []OauthClient

	err := a.db.Get(&clients, querySelectClients)
	if err != nil {
		return []OauthClient{}, err
	}

	return clients, nil
}

func (a *TokenStore) resolveClientByClientID(clientID string) (client OauthClient, err error) {
	err = a.db.Get(&client, querySelectClients+" WHERE client_id = ?", clientID)
	switch {
	case err == sql.ErrNoRows:
		err = errors.New(ErrorClientNotFound)
		return
	case err != nil:
		return
	}

	return
}

func (a *TokenStore) resolveByTelephoneOrEmail(username string) (User, error) {
	var user User

	err := a.db.Get(&user, querySelectUser+" WHERE telephone = ? OR  email = ?", username, username)
	switch {
	case err == sql.ErrNoRows:
		return User{}, errors.New(ErrorClientNotFound)
	case err != nil:
		return User{}, err
	}

	return user, nil
}
