package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/ut-sama-art-studio/art-market-backend/services/oauth"
	"github.com/ut-sama-art-studio/art-market-backend/tests"
)

func (app *Application) AddAPIRoutes(router *chi.Mux) {
	router.Get("/hello-world", tests.HelloWorldHandler)

	oauth.InitOAuth()
	router.Route("/auth/discord", func(r chi.Router) {
		r.HandleFunc("/login", oauth.HandleDiscordLogin)
		r.HandleFunc("/callback", oauth.HandleDiscordCallback)
	})
}
