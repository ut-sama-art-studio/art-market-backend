package oauth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ut-sama-art-studio/art-market-backend/services/users"
	"github.com/ut-sama-art-studio/art-market-backend/utils/jwt"
	"golang.org/x/oauth2"
)

// HandleDiscordLogin initiates the OAuth2 flow with Discord.
func HandleDiscordLogin(w http.ResponseWriter, r *http.Request) {
	// Generate the OAuth2 URL for Discord
	url := OAuthConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// HandleDiscordCallback handles the OAuth2 callback from Discord.
func HandleDiscordCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Code not found", http.StatusBadRequest)
		return
	}

	// Exchange the code for a token
	token, err := OAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Use the token to get user info from Discord
	client := OAuthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://discord.com/api/users/@me")
	if err != nil {
		http.Error(w, "Failed to get user info: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Decode the user info from the response
	var userInfo map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		http.Error(w, "Failed to decode user info: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Extract user details
	discordID, ok := userInfo["id"].(string)
	if !ok {
		http.Error(w, "Invalid user info: ID not found", http.StatusInternalServerError)
		return
	}
	oauthID := DiscordIdToOauthId(discordID)
	username, ok := userInfo["username"].(string)
	if !ok {
		http.Error(w, "Invalid user info: Username not found", http.StatusInternalServerError)
		return
	}
	email, _ := userInfo["email"].(string) // Email is optional

	existingUser, err := users.GetUserByOauthID(oauthID)
	if err != nil {
		http.Error(w, "Failed to fetch user: "+err.Error(), http.StatusInternalServerError)
		return
	}
	// If user does not exist, create a new one
	if existingUser == nil {
		newUser := users.User{
			OauthID: oauthID,
			Name:    username,
			Email:   &email,
		}
		newUserID, err := newUser.Insert()
		if err != nil {
			http.Error(w, "Failed to create user: "+err.Error(), http.StatusInternalServerError)
			return
		}
		existingUser, err = users.GetUserByID(newUserID)
		if err != nil {
			http.Error(w, "Failed to fetch user: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	jwtToken, err := jwt.GenerateToken(existingUser)
	if err != nil {
		http.Error(w, "Failed to create JWT: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Redirect to the frontend with the token
	redirectionUrl := FrontEndURL + "?token=" + jwtToken
	http.Redirect(w, r, redirectionUrl, http.StatusSeeOther)
}

func DiscordIdToOauthId(id string) string {
	return fmt.Sprintf("discord%s", id)
}
