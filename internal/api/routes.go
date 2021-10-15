package api

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"

	"github.com/m4tthewde/fdmxyz/internal/config"
	"github.com/m4tthewde/fdmxyz/internal/db"
	"github.com/m4tthewde/fdmxyz/internal/object"
	"github.com/m4tthewde/fdmxyz/internal/twitch"
	"github.com/nicklaw5/helix/v2"
)

type RouteHandler struct {
	config        *config.Config
	mongoHandler  *db.MongoHandler
	twitchHandler *twitch.Handler
}

func (rh *RouteHandler) get() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
	}
}

func (rh *RouteHandler) register() func(
	w http.ResponseWriter,
	r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()
		if !params.Has("type") || !params.Has("user_id") {
			http.Error(w, "Missing url parameter", http.StatusBadRequest)
			return
		}

		// create new webhook object
		var webhook object.Webhook
		webhook.UserID = params.Get("user_id")
		webhook.Status = object.PENDING
		webhook.Secret = randSeq(64)

		switch params.Get("type") {
		case "follow":
			webhook.Typing = "channel.follow"
			webhook.Callback = "/twitch/follow"
		case "sub":
			webhook.Typing = "channel.subscribe"
			webhook.Callback = "/twitch/subscribe"
		}
		// save webhook object in db
		rh.mongoHandler.SaveWebhook(webhook)

		// register webhook at twitch
		resp, err := rh.twitchHandler.RegisterWebhook(&webhook)
		if err != nil {
			log.Println(err)
			http.Error(w, "error registering webhook", http.StatusInternalServerError)
		}

		// TODO write good feedback
		_, err = w.Write([]byte(resp.Data.EventSubSubscriptions[0].ID))
		if err != nil {
			log.Println(err)
			http.Error(w, "error registering webhook", http.StatusInternalServerError)
		}
	}
}

func (rh *RouteHandler) delete() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
	}
}

func (rh *RouteHandler) twitchFollow() func(
	w http.ResponseWriter,
	r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			return
		}
		defer r.Body.Close()
		// verify that the notification came from twitch using the secret.
		if !helix.VerifyEventSubNotification(
			rh.mongoHandler.GetPendingWebhook().Secret,
			r.Header,
			string(body)) {
			http.Error(w, "no valid signature", http.StatusBadRequest)
			return
		}

		var vals object.EventSubNotification
		err = json.NewDecoder(bytes.NewReader(body)).Decode(&vals)
		if err != nil {
			log.Println(err)
			return
		}
		// if there's a challenge in the request,
		// respond with only the challenge to verify your eventsub.
		if vals.Challenge != "" {
			_, err = w.Write([]byte(vals.Challenge))
			if err != nil {
				log.Println(err)
				http.Error(w, "error registering webhook", http.StatusInternalServerError)
			}
			return
		}

		var followEvent helix.EventSubChannelFollowEvent
		err = json.NewDecoder(bytes.NewReader(vals.Event)).Decode(&followEvent)
		if err != nil {
			panic(err)
		}

		// TODO too many open files after some time, something is leaking
		log.Println(followEvent)
	}
}

func (rh *RouteHandler) twitchSubscribe() func(
	w http.ResponseWriter,
	r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			return
		}
		defer r.Body.Close()
		// verify that the notification came from twitch using the secret.
		if !helix.VerifyEventSubNotification(
			rh.mongoHandler.GetPendingWebhook().Secret,
			r.Header,
			string(body)) {
			http.Error(w, "no valid signature", http.StatusBadRequest)
			return
		}

		var vals object.EventSubNotification
		err = json.NewDecoder(bytes.NewReader(body)).Decode(&vals)
		if err != nil {
			log.Println(err)
			return
		}
		// if there's a challenge in the request,
		// respond with only the challenge to verify your eventsub.
		if vals.Challenge != "" {
			_, err = w.Write([]byte(vals.Challenge))
			if err != nil {
				log.Println(err)
				http.Error(w, "error registering webhook", http.StatusInternalServerError)
			}
			return
		}

		var subscribeEvent helix.EventSubChannelSubscribeEvent
		err = json.NewDecoder(bytes.NewReader(vals.Event)).Decode(&subscribeEvent)
		if err != nil {
			panic(err)
		}

		log.Println(subscribeEvent)
	}
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
