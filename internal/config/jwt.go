package config

import "time"

type JWTConfig struct {
	SecretKey       string
	TokenExpiration time.Duration
	RefreshDuration time.Duration
	IssuerName      string
}

func NewJWTConfig(jwtSecret string) JWTConfig {
	return JWTConfig{
		SecretKey:       jwtSecret,
		TokenExpiration: 15 * time.Minute,
		RefreshDuration: 30 * 24 * time.Hour,
		IssuerName:      "ecom",
	}
}
