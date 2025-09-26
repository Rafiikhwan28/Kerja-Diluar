package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"job-service/internal/handler"
	"job-service/internal/repository"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ .env file not found, using system environment")
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL not set")
	}

	// Connect PostgreSQL
	db, err := sqlx.Connect("postgres", dbURL)
	if err != nil {
		log.Fatal("❌ Failed to connect DB:", err)
	}
	defer db.Close()
	fmt.Println("✅ Connected to PostgreSQL (Job Service)")

	// Init repo & handler
	repo := repository.NewJobRepo(db)
	h := handler.NewJobHandler(repo)

	// Setup router
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// Health check
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong from job-service"))
	})

	// Job routes
	r.Route("/jobs", func(r chi.Router) {
		r.Get("/", h.GetJobs)                                // GET /jobs
		r.Post("/", h.CreateJob)                             // POST /jobs
		r.Get("/{id}", h.GetJobByID)                         // GET /jobs/{id}
		r.Put("/{id}", h.UpdateJob)                          // PUT /jobs/{id}
		r.Delete("/{id}", h.DeleteJob)                       // DELETE /jobs/{id}
		r.Get("/category/{categoryID}", h.GetJobsByCategory) // GET /jobs/category/2
	})

	// Run server
	port := os.Getenv("JOB_SERVICE_PORT")
	if port == "" {
		port = "8002"
	}
	fmt.Println("🚀 Job service running on port", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
