package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	LogLevel string `env:"LOG_LEVEL" default:"info"`
	APIURL   string
	PORT     string
}

func Load(path string) (*Config, error) {
	if err := godotenv.Load(path); err != nil {
		return nil, err
	}

	cfg := &Config{
		APIURL: os.Getenv("API_URL"),
		PORT:   os.Getenv("PORT"),
	}

	return cfg, nil
}
