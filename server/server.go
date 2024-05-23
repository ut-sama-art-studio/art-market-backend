package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/joho/godotenv"
	"github.com/ut-sama-art-studio/art-market-backend/database"
	"github.com/ut-sama-art-studio/art-market-backend/graph"
)

type Config struct {
	Port string
}

type Application struct {
	Config Config
	// Models
}

func (app *Application) Serve(port string) error {
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	return http.ListenAndServe(":"+port, nil)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	port := os.Getenv("PORT")
	fmt.Println("API is listening on port", port)

	cfg := Config{
		Port: port,
	}

	dsn := os.Getenv("DSN")
	dbConn, err := database.ConnectPostgres(dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer dbConn.DB.Close()

	app := &Application{
		Config: cfg,
		// TODO: add models
	}

	err = app.Serve(port)
	if err != nil {
		log.Fatal(err)
	}
}
