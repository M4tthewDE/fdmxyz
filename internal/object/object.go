package object

type Webhook struct {
	Id      string
	User_id string
	Typing  WebhookType
	Status  WebhookStatus
}

type WebhookType int

const (
	FOLLOW WebhookType = iota
	SUB                = iota
)

type WebhookStatus int

const (
	PENDING   WebhookStatus = iota
	CONFIRMED               = iota
)
