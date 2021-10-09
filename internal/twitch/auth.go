package twitch

import (
	"time"

	"github.com/m4tthewde/fdmxyz/internal/config"
	"github.com/m4tthewde/fdmxyz/internal/db"
	"github.com/m4tthewde/fdmxyz/internal/object"
)

type AuthenticationHandler struct {
	config        *config.Config
	mongoHandler  *db.MongoHandler
	twitchHandler *TwitchHandler
}

func NewAuthenticationHandler(config *config.Config) *AuthenticationHandler {
	authHandler := AuthenticationHandler{
		config:        config,
		mongoHandler:  &db.MongoHandler{Config: config},
		twitchHandler: &TwitchHandler{},
	}
	authHandler.twitchHandler = &TwitchHandler{
		Config:      config,
		AuthHandler: &authHandler,
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

	isValid, err := ah.twitchHandler.ValidateToken(auth.Token)
	if err != nil {
		panic(err)
	}

	return isValid
}

func (ah *AuthenticationHandler) GenerateToken() error {
	// request token from twitch
	resp, err := ah.twitchHandler.RequestToken()
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
			resp, err := ah.twitchHandler.RequestToken()
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
