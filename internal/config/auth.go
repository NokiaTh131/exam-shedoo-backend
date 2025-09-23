package config

import (
	"os"

	_ "github.com/joho/godotenv/autoload"
)

type AuthConfig struct {
	TokenURL     string
	RedirectURL  string
	ClientID     string
	ClientSecret string
	Scope        string
	BasicInfoURL string
	JWTSecret    string
	IsProd       bool
	CookieDomain string
}

func LoadAuthConfig() *AuthConfig {
	return &AuthConfig{
		TokenURL:     os.Getenv("CMU_ENTRAID_GET_TOKEN_URL"),
		RedirectURL:  os.Getenv("CMU_ENTRAID_REDIRECT_URL"),
		ClientID:     os.Getenv("CMU_ENTRAID_CLIENT_ID"),
		ClientSecret: os.Getenv("CMU_ENTRAID_CLIENT_SECRET"),
		Scope:        os.Getenv("SCOPE"),
		BasicInfoURL: os.Getenv("CMU_ENTRAID_GET_BASIC_INFO"),
		JWTSecret:    os.Getenv("JWT_SECRET"),
		IsProd:       os.Getenv("PROD") == "true",
		CookieDomain: os.Getenv("COOKIE_DOMAIN"),
	}
}
