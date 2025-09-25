package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"job-service/internal/handler"
	"job-service/internal/repository"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ .env file not found, using system environment")
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL not set")
	}

	db, err := sqlx.Connect("postgres", dbURL)
	if err != nil {
		log.Fatal("Failed to connect DB:", err)
	}
	defer db.Close()
	fmt.Println("✅ Connected to PostgreSQL (Job Service)")

	repo := repository.NewJobRepo(db)
	h := handler.NewJobHandler(repo)

	r := chi.NewRouter()
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong from job-service"))
	})

	r.Get("/jobs", h.GetJobs)
	r.Post("/jobs", h.CreateJob)

	port := os.Getenv("JOB_SERVICE_PORT")
	if port == "" {
		port = "8002"
	}
	fmt.Println("Job service running on port", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
