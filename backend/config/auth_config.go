package config

import "os"

type AuthConfig struct {
	GoogleClientID     string
	GoogleClientSecret string
	JWTSecret          string
	CookieSecret       string
}

func NewAuthConfig() *AuthConfig {
	return &AuthConfig{
		GoogleClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		GoogleClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		JWTSecret:          os.Getenv("JWT_SECRET"),
		CookieSecret:       os.Getenv("COOKIE_SECRET"),
	}
}
