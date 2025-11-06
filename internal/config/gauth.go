package config

import (
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func NewGoogleAuth(config *viper.Viper) *oauth2.Config {
	secretKey := config.GetString("oauth.secret.key")
	redirectURL := config.GetString("oauth.redirect.url")
	clientID := config.GetString("oauth.client.id")
	return &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: secretKey,
		RedirectURL:  redirectURL,
		Endpoint:     google.Endpoint,
	}
}
