package config

import (
	"errors"
	"os"
)

type Config struct {
	HTTPPort string
	password string
}

func Load() (*Config, error) {
	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		 return  nil, errors.New("HTTP_PORT environment variable is required")
	}

	password := os.Getenv("LOGIN_PASSWORD")
	if password == "" {
		return nil, errors.New("LOGIN_PASSWORD environment variable is required")
	}

	return &Config{
		HTTPPort: httpPort,
		password: password,
	}, nil
}