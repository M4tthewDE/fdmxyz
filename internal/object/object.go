package object

import (
	"encoding/json"

	"github.com/nicklaw5/helix/v2"
)

type Webhook struct {
	ID       string
	UserID   string
	Typing   string
	Status   WebhookStatus
	Callback string
}

type WebhookStatus int

const (
	PENDING   WebhookStatus = iota
	CONFIRMED               = iota
)

type Authentication struct {
	Token     string
	ExpiresIn int
}

type EventSubNotification struct {
	Subscription helix.EventSubSubscription `json:"subscription"`
	Challenge    string                     `json:"challenge"`
	Event        json.RawMessage            `json:"event"`
}
