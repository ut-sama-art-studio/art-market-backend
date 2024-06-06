package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/ut-sama-art-studio/art-market-backend/database"
	"github.com/ut-sama-art-studio/art-market-backend/graph"
	"github.com/ut-sama-art-studio/art-market-backend/graph/resolvers"
)

type Config struct {
	Port string
}

type Application struct {
	Config Config
	// Models
}

func (app *Application) Serve(router *chi.Mux) {
	server := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &resolvers.Resolver{}}))

	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", server)

	address := fmt.Sprintf(":%s", app.Config.Port)
	log.Printf("Connect to http://localhost%s/ for GraphQL playground", address)

	if err := http.ListenAndServe(address, router); err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("PORT")
	fmt.Println("API is listening on port", port)

	cfg := Config{Port: port}
	router := chi.NewRouter()

	dbString := os.Getenv("DBSTRING")

	if err := database.InitDB(dbString); err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer database.CloseDB()

	app := &Application{
		Config: cfg,
	}

	app.Serve(router)
}
