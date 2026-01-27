package config

import (
	"fmt"
	"os"
)

type Config struct {
	MongoURI  string
	DBName    string
	JWTSecret string
	Port      string
}

func Load() (*Config, error) {
	cfg := &Config{
		MongoURI:  os.Getenv("MONGODB_URI"),
		DBName:    os.Getenv("DB_NAME"),
		JWTSecret: os.Getenv("JWT_SECRET"),
		Port:      os.Getenv("PORT"),
	}

	if cfg.MongoURI == "" {
		return nil, fmt.Errorf("MONGODB_URI is required")
	}
	if cfg.DBName == "" {
		return nil, fmt.Errorf("DB_NAME is required")
	}
	if cfg.Port == "" {
		return nil, fmt.Errorf("PORT is required")
	}

	if cfg.JWTSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET is required")
	}

	return cfg, nil
}
