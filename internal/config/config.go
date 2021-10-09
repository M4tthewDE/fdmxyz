package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

func GetConfig(file string) *Config {
	f, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}

	decoder := yaml.NewDecoder(f)

	var config Config

	err = decoder.Decode(&config)
	if err != nil {
		log.Fatal(err)
	}

	f.Close()

	return &config
}

type Config struct {
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Database struct {
		User       string `yaml:"user"`
		Password   string `yaml:"password"`
		Name       string `yaml:"name"`
		Collection string `yaml:"collection"`
	} `yaml:"database"`
	Twitch struct {
		ClientID string `yaml:"client_id"`
		Secret   string `yaml:"secret"`
	} `yaml:"twitch"`
	Api struct {
		BaseURL string `yaml:"base_url"`
	} `yaml:"api"`
}
