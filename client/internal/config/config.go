package config

import (
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/viper"
)

type AppConfig struct {
	ClientPort int
	ServerPort int
	ServerHost string
}

func ReadFromEnv() *AppConfig {
	viper.SetConfigFile("./../.env")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("error while reading config %s", err)
	}

	// for docker container networking
	host := os.Getenv("SERVER_HOST")
	if host == "" {
		host = viper.GetString("SERVER_HOST")
	}

	return &AppConfig{
		ServerPort: viper.GetInt("PORT"),
		ClientPort: viper.GetInt("CLIENT_PORT"),
		ServerHost: host,
	}
}
