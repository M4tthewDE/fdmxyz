package api

import (
	"net/http"

	"github.com/m4tthewde/fdmxyz/internal/config"
	"github.com/m4tthewde/fdmxyz/internal/db"
	"github.com/m4tthewde/fdmxyz/internal/object"
)

type RouteHandler struct {
	locked       bool
	config       *config.Config
	mongoHandler *db.MongoHandler
}

func (rh *RouteHandler) get() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
	}
}

func (rh *RouteHandler) register() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()
		if !params.Has("type") || !params.Has("user_id") {
			http.Error(w, "Missing url parameter", http.StatusBadRequest)
			return
		}

		// create new webhook object
		var webhook object.Webhook
		webhook.user_id = params.Get("user_id")
		webhook.status = object.PENDING

		switch params.Get("type") {
		case "follow":
			webhook.typing = object.FOLLOW
		case "sub":
			webhook.typing = object.SUB
		}
		// save webhook object in db
		rh.mongoHandler.SaveWebhook(webhook)
		// register webhook at twitch
	}
}

func (rh *RouteHandler) delete() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
	}
}

func (rh *RouteHandler) twitch() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
	}
}
