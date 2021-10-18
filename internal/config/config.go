package config

import (
	"log"
	"math/rand"
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

	config.Secret = randSeq(64)

	return &config
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
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
	API struct {
		BaseURL string `yaml:"base_url"`
	} `yaml:"api"`
	Secret string `yaml:"secret"`
}
