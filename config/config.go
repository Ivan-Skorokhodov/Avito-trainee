package config

import "os"

type Config struct {
	Database struct {
		Host     string
		Port     string
		User     string
		Password string
		Name     string
	}

	Server struct {
		Port string
	}
}

func LoadConfig() *Config {
	// Load environment variables from docker compose
	return &Config{
		Database: struct {
			Host     string
			Port     string
			User     string
			Password string
			Name     string
		}{
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Name:     os.Getenv("DB_NAME"),
		},
		Server: struct {
			Port string
		}{
			Port: os.Getenv("APP_PORT"),
		},
	}
}
