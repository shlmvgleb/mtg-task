package config

import (
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/viper"
)

type PostgresConfig struct {
	Port     int
	Host     string
	Password string
	User     string
	DbName   string
}

type AppConfig struct {
	Port     int
	Postgres *PostgresConfig
}

func ReadFromEnv() *AppConfig {
	viper.SetConfigFile("./../.env")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("error while reading config %s", err)
	}

	// for docker container networking
	host := os.Getenv("POSTGRES_HOST")
	if host == "" {
		host = viper.GetString("POSTGRES_HOST")
	}

	return &AppConfig{
		Port: viper.GetInt("PORT"),
		Postgres: &PostgresConfig{
			Port:     viper.GetInt("POSTGRES_PORT"),
			Host:     host,
			User:     viper.GetString("POSTGRES_USER"),
			Password: viper.GetString("POSTGRES_PWD"),
			DbName:   viper.GetString("POSTGRES_DB_NAME"),
		},
	}
}
