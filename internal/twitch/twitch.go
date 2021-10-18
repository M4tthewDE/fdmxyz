package twitch

import (
	"github.com/m4tthewde/fdmxyz/internal/config"
	"github.com/m4tthewde/fdmxyz/internal/object"
	"github.com/nicklaw5/helix/v2"
)

type Handler struct {
	Config      *config.Config
	AuthHandler *AuthenticationHandler
}

func (th *Handler) getClient() (*helix.Client, error) {
	client, err := helix.NewClient(&helix.Options{
		ClientID:       th.Config.Twitch.ClientID,
		AppAccessToken: th.AuthHandler.GetAuth(),
	})

	return client, err
}

func (th *Handler) RegisterWebhook(webhook *object.Webhook) (
	*helix.EventSubSubscriptionsResponse, error) {
	client, err := th.getClient()
	if err != nil {
		return nil, err
	}

	eventSubSubscription := helix.EventSubSubscription{
		Type:    webhook.Typing,
		Version: "1",
		Condition: helix.EventSubCondition{
			BroadcasterUserID: webhook.UserID,
		},
		Transport: helix.EventSubTransport{
			Method:   "webhook",
			Callback: th.Config.API.BaseURL + webhook.Callback,
			Secret:   th.Config.Secret,
		},
	}

	resp, err := client.CreateEventSubSubscription(&eventSubSubscription)
	return resp, err
}
