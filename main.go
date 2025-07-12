package main

import (
	// "fmt"
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/troyboy95/RSS-project-test/internal/database"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	// fmt.Println("Hello, World!")
	// feed, err := urlToFeed("https://wagslane.dev/index.xml")
	// if err != nil {
	// 	log.Fatal("Error fetching RSS feed:", err)
	// }
	// log.Println("Fetched feed:", feed)
	godotenv.Load(".env")

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT environment variable is not found!")
	}

	db := os.Getenv("DB_URL")
	if db == "" {
		log.Fatal("DB_URL environment variable is not found!")
	}

	conn, err := sql.Open("postgres", db)
	if err != nil {
		log.Fatal("Error connecting to the database:", err)
	}

	dbConn := database.New(conn)
	apiCfg := apiConfig{
		DB: dbConn,
	}

	go startScraper(dbConn, 10, time.Minute)

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1Router := chi.NewRouter()
	v1Router.Get("/ready", readinessHandler)
	v1Router.Get("/error", handlerError)
	v1Router.Post("/users", apiCfg.createUserHandler)
	v1Router.Get("/users", apiCfg.MiddlewareAuth(apiCfg.getUserHandler))
	v1Router.Post("/feeds", apiCfg.MiddlewareAuth(apiCfg.createFeedHandler))
	v1Router.Get("/feeds", apiCfg.getFeedsHandler)

	v1Router.Post("/feed_follows", apiCfg.MiddlewareAuth(apiCfg.createFeedFollowHandler))
	v1Router.Get("/feed_follows", apiCfg.MiddlewareAuth(apiCfg.getFeedFollowsHandler))
	v1Router.Delete("/feed_follows/{FeedFollowID}", apiCfg.MiddlewareAuth(apiCfg.deleteFeedFollowHandler))

	v1Router.Get("/users/posts", apiCfg.MiddlewareAuth(apiCfg.getUserPostsHandler))

	router.Mount("/v1", v1Router)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	log.Printf("Server is running on port: %v", port)

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal("Error starting server:", err)
	}

}
