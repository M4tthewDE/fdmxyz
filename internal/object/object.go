package object

import (
	"encoding/json"

	"github.com/nicklaw5/helix/v2"
)

type Webhook struct {
	Id       string
	User_id  string
	Typing   string
	Status   WebhookStatus
	Secret   string
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
	Secret    string
}

type EventSubNotification struct {
	Subscription helix.EventSubSubscription `json:"subscription"`
	Challenge    string                     `json:"challenge"`
	Event        json.RawMessage            `json:"event"`
}
