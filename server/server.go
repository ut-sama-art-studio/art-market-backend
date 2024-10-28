package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"github.com/ut-sama-art-studio/art-market-backend/database"
	"github.com/ut-sama-art-studio/art-market-backend/graph"
	"github.com/ut-sama-art-studio/art-market-backend/graph/directives"
	"github.com/ut-sama-art-studio/art-market-backend/graph/resolvers"
	"github.com/ut-sama-art-studio/art-market-backend/middlewares"
	"github.com/ut-sama-art-studio/art-market-backend/services/files"
)

type Config struct {
	Port      string
	ApiPrefix string
}

type Application struct {
	Config Config
}

func (app *Application) Serve(router *chi.Mux) {
	apiRouter := chi.NewRouter()
	// middlewares
	apiRouter.Use(middlewares.AuthMiddleware)
	apiRouter.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{os.Getenv("FRONT_END_URL_LOCALHOST"), os.Getenv("FRONT_END_URL")},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,

		// Debug: os.Getenv("ENV") == "development",
	}))

	// Initializing the S3 client
	err := files.InitS3Client()
	files.LogBucketObjects()
	if err != nil {
		log.Fatalf("failed to initialize S3 client: %v", err)
	}

	router.Mount(app.Config.ApiPrefix, apiRouter)
	app.AddAPIRoutes(apiRouter)

	// Add graphql
	graphConfig := graph.Config{Resolvers: &resolvers.Resolver{}}
	graphConfig.Directives.Auth = directives.AuthDirective
	server := handler.NewDefaultServer(graph.NewExecutableSchema(graphConfig))
	if os.Getenv("NODE_ENV") != "production" {
		apiRouter.Handle("/", playground.Handler("GraphQL playground", app.Config.ApiPrefix+"/graphql"))
		log.Printf("Connect to http://localhost:%s%s for GraphQL playground", app.Config.Port, app.Config.ApiPrefix)
	}
	apiRouter.Handle("/graphql", server)

	address := fmt.Sprintf(":%s", app.Config.Port)
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

	cfg := Config{Port: port, ApiPrefix: "/api"}
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
