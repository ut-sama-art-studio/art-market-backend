package oauth

import (
	"os"

	"golang.org/x/oauth2"
)

var OAuthConfig *oauth2.Config
var FrontEndURL string

// InitOAuth initializes the OAuth configuration for Discord.
func InitOAuth() {
	OAuthConfig = &oauth2.Config{
		ClientID:     os.Getenv("DISCORD_CLIENT_ID"),
		ClientSecret: os.Getenv("DISCORD_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("DISCORD_REDIRECT_URL"),
		Scopes:       []string{"identify", "email"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://discord.com/api/oauth2/authorize",
			TokenURL: "https://discord.com/api/oauth2/token",
		},
	}

	if os.Getenv("ENV") == "production" {
		FrontEndURL = "https://utsama-art-studio.vercel.app"
	} else {
		FrontEndURL = "http://localhost:3000"
	}
}
