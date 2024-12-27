package config

import (
	"log"
	"os"
)

type EnvConfig struct {
	DATABASE_URL string
	SERVER_PORT  string
	JWT_SECRET   string
}

func LoadEnvConfig() *EnvConfig {
	return &EnvConfig{
		DATABASE_URL: MustGetEnv("DATABASE_URL"),
		SERVER_PORT:  MustGetEnv("SERVER_PORT"),
		JWT_SECRET:   MustGetEnv("JWT_SECRET"),
	}
}

func MustGetEnv(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		log.Fatalf("FATAL: Environment variable %s is not set!", key)
	}
	return value
}
