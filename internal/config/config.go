package config

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	DB            DBConfig     `envconfig:"DB"`
	Server        ServerConfig `envconfig:"APP"`
	AirlineAPIKey string       `envconfig:"AIRLINE_API_KEY"`
}

type DBConfig struct {
	Host     string
	Name     string
	Port     int
	User     string
	Password string
	Schema   string
}

type ServerConfig struct {
	Port int
}

func Load() (Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Printf("warning: .env file not found: %v", err)
		return Config{}, err
	}

	var appConfig Config

	if err := envconfig.Process("", &appConfig); err != nil {
		log.Fatalf("error in processing config: %v", err)
		return Config{}, err
	}

	return appConfig, nil
}
