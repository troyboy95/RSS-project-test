package main

import (
	// "fmt"
	"database/sql"
	"log"
	"net/http"
	"os"

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

	apiCfg := apiConfig{
		DB: database.New(conn),
	}

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

	router.Mount("/v1", v1Router)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	log.Printf("Server is running on port: %v", port)

	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal("Error starting server:", err)
	}

}
