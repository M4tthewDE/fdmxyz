package twitch

import (
	"time"

	"github.com/m4tthewde/fdmxyz/internal/config"
	"github.com/m4tthewde/fdmxyz/internal/db"
	"github.com/m4tthewde/fdmxyz/internal/object"
	"github.com/nicklaw5/helix/v2"
)

type AuthenticationHandler struct {
	config       *config.Config
	mongoHandler *db.MongoHandler
}

func NewAuthenticationHandler(config *config.Config) *AuthenticationHandler {
	authHandler := AuthenticationHandler{
		config:       config,
		mongoHandler: &db.MongoHandler{Config: config},
	}

	return &authHandler
}

func (ah *AuthenticationHandler) IsValid() bool {
	// check if auth is stored in db
	// for example the case if first time startup
	auth := ah.mongoHandler.GetAuth()
	if auth == nil {
		return false
	}

	isValid, err := ah.ValidateToken(auth.Token)
	if err != nil {
		panic(err)
	}

	return isValid
}

func (ah *AuthenticationHandler) GenerateToken() error {
	// request token from twitch
	resp, err := ah.RequestToken()
	if err != nil {
		return err
	}

	var auth object.Authentication
	auth.Token = resp.Data.AccessToken
	auth.ExpiresIn = resp.Data.ExpiresIn

	// save auth in database
	ah.mongoHandler.SaveAuth(&auth)

	return nil
}

func (ah *AuthenticationHandler) RegenerateTokenJob() {
	for {
		// if token expires in less than one day
		auth := ah.mongoHandler.GetAuth()
		if auth.ExpiresIn < 1000000 {
			resp, err := ah.RequestToken()
			if err != nil {
				panic(err)
			}
			// save new token details in databsae
			auth.Token = resp.Data.AccessToken
			auth.ExpiresIn = resp.Data.ExpiresIn

			// TODO errors out!!!!!!
			ah.mongoHandler.UpdateAuth(auth)
		}

		time.Sleep(1 * time.Hour)
	}
}

func (ah *AuthenticationHandler) GetAuth() string {
	auth := ah.mongoHandler.GetAuth()
	if auth == nil {
		panic("No auth found")
	}
	return auth.Token
}

func (ah *AuthenticationHandler) ValidateToken(token string) (bool, error) {
	client, err := helix.NewClient(&helix.Options{
		ClientID: ah.config.Twitch.ClientID,
	})
	if err != nil {
		return false, nil
	}

	isValid, _, err := client.ValidateToken(token)
	if err != nil {
		return false, err
	}

	return isValid, nil
}

func (ah *AuthenticationHandler) RequestToken() (*helix.AppAccessTokenResponse, error) {
	client, err := helix.NewClient(&helix.Options{
		ClientID:     ah.config.Twitch.ClientID,
		ClientSecret: ah.config.Twitch.Secret,
	})
	if err != nil {
		return nil, err
	}

	resp, err := client.RequestAppAccessToken([]string{})
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (ah *AuthenticationHandler) RefreshToken(refreshToken string) (*helix.RefreshTokenResponse, error) {
	client, err := helix.NewClient(&helix.Options{
		ClientID:     ah.config.Twitch.ClientID,
		ClientSecret: ah.config.Twitch.Secret,
	})
	if err != nil {
		return nil, err
	}

	resp, err := client.RefreshUserAccessToken(refreshToken)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
