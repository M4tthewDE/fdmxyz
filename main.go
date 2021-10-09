package main

import (
	"github.com/m4tthewde/fdmxyz/internal/api"
	"github.com/m4tthewde/fdmxyz/internal/config"
	auth "github.com/m4tthewde/fdmxyz/internal/twitch"
)

func main() {
	config := config.GetConfig("config.yml")

	authHandler := auth.NewAuthenticationHandler(config)

	// check if auth is valid
	if !authHandler.IsValid() {
		err := authHandler.GenerateToken()
		if err != nil {
			panic(err)
		}
	}

	go authHandler.RegenerateTokenJob()

	server := api.NewServer(config)

	server.Run()

	for {
		select {}
	}
}
