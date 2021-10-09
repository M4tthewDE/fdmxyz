package main

import (
	"github.com/m4tthewde/fdmxyz/internal/api"
	"github.com/m4tthewde/fdmxyz/internal/config"
)

func main() {
	config := config.GetConfig("config.yml")

	server := api.NewServer(config)

	server.Run()
}
