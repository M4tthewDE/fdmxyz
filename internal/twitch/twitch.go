package twitch

import (
	"github.com/m4tthewde/fdmxyz/internal/config"
	"github.com/m4tthewde/fdmxyz/internal/object"
	"github.com/nicklaw5/helix/v2"
)

type TwitchHandler struct {
	Config      *config.Config
	AuthHandler *AuthenticationHandler
}

func (th *TwitchHandler) getClient() (*helix.Client, error) {
	client, err := helix.NewClient(&helix.Options{
		ClientID:       th.Config.Twitch.ClientID,
		AppAccessToken: th.AuthHandler.GetAuth(),
	})

	return client, err
}

func (th *TwitchHandler) RegisterWebhook(webhook *object.Webhook) (*helix.EventSubSubscriptionsResponse, error) {
	client, err := th.getClient()
	if err != nil {
		return nil, err
	}

	eventSubSubscription := helix.EventSubSubscription{
		Type:    webhook.Typing,
		Version: "1",
		Condition: helix.EventSubCondition{
			BroadcasterUserID: webhook.User_id,
		},
		Transport: helix.EventSubTransport{
			Method:   "webhook",
			Callback: th.Config.Api.BaseURL + webhook.Callback,
			Secret:   webhook.Secret,
		},
	}

	resp, err := client.CreateEventSubSubscription(&eventSubSubscription)
	return resp, err
}

func (th *TwitchHandler) ValidateToken(token string) (bool, error) {
	client, err := helix.NewClient(&helix.Options{
		ClientID: th.Config.Twitch.ClientID,
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

func (th *TwitchHandler) RequestToken() (*helix.AppAccessTokenResponse, error) {
	client, err := helix.NewClient(&helix.Options{
		ClientID:     th.Config.Twitch.ClientID,
		ClientSecret: th.Config.Twitch.Secret,
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

func (th *TwitchHandler) RefreshToken(refreshToken string) (*helix.RefreshTokenResponse, error) {
	client, err := helix.NewClient(&helix.Options{
		ClientID:     th.Config.Twitch.ClientID,
		ClientSecret: th.Config.Twitch.Secret,
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
