package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/iamarju/rss-aggregator/internal/database"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {

	godotenv.Load(".env")

	portString := os.Getenv("PORT")

	if portString == "" {
		log.Fatal("Port not found in the env")
	}

	dbUrl := os.Getenv("DB_URL")

	if dbUrl == "" {
		log.Fatal("DB_URL not found in the env")
	}

	conn, err := sql.Open("postgres", dbUrl)

	if err != nil {
		log.Fatal("Cannot connect to dabatase: %v", err)

	}

	db := database.New(conn)
	apiCgf := apiConfig{
		DB: db,
	}

	go startScraping(db, 10, time.Minute)

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1Router := chi.NewRouter()
	v1Router.Get("/healthz", handlerReadiness)
	v1Router.Get("/err", handlerErr)
	v1Router.Get("/users", apiCgf.middlewareAuth(apiCgf.handlerGetUser))
	v1Router.Post("/users", apiCgf.handlerCreateUser)
	v1Router.Post("/feeds", apiCgf.middlewareAuth(apiCgf.handlerCreateFeed))
	v1Router.Get("/feeds", apiCgf.handlerGetFeed)
	v1Router.Post("/feed_follows", apiCgf.middlewareAuth(apiCgf.handlerCreateFeedFollow))
	v1Router.Get("/feed_follows", apiCgf.middlewareAuth(apiCgf.handlerGetFeedFollow))
	v1Router.Delete("/feed_follows/{feedFollowId}", apiCgf.middlewareAuth(apiCgf.handlerDeleteFeedFollow))
	router.Mount("/v1", v1Router)

	server := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}

	log.Printf("server starting at port: %s", portString)

	error := server.ListenAndServe()

	if error != nil {
		log.Fatal(error)
	}

	fmt.Println("Port: $s", portString)
}
