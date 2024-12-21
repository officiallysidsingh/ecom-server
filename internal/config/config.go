package config

import (
	"log"
	"os"
)

type AppConfig struct {
	DATABASE_URL string
	SERVER_PORT  string
}

func LoadConfig() *AppConfig {
	return &AppConfig{
		DATABASE_URL: MustGetEnv("DATABASE_URL"),
		SERVER_PORT:  MustGetEnv("SERVER_PORT"),
	}
}

func MustGetEnv(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		log.Fatalf("FATAL: Environment variable %s is not set!", key)
	}
	return value
}
