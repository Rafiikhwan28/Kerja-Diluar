package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

var db *sqlx.DB
var jwtSecret string

func main() {
	// Load .env
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ .env file not found, using system environment")
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL not set")
	}

	jwtSecret = os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET not set")
	}

	// Connect DB
	var err error
	db, err = sqlx.Connect("postgres", dbURL)
	if err != nil {
		log.Fatal("Failed to connect DB:", err)
	}
	defer db.Close()
	fmt.Println("✅ Connected to PostgreSQL (Auth Service)")

	// Router
	r := chi.NewRouter()
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong from auth-service"))
	})

	// Auth endpoints
	r.Post("/register", RegisterHandler)
	r.Post("/login", LoginHandler)

	port := os.Getenv("AUTH_SERVICE_PORT")
	if port == "" {
		port = "8001"
	}
	fmt.Println("Auth service running on port", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

// ------------------ Handlers ------------------

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if req.Name == "" || req.Email == "" || req.Password == "" {
		http.Error(w, "name/email/password required", http.StatusBadRequest)
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "failed to hash password", http.StatusInternalServerError)
		return
	}

	// Insert to DB (note: column name password_hash sesuai DB)
	var userID int
	query := `INSERT INTO users (name, email, password_hash) VALUES ($1, $2, $3) RETURNING id`
	err = db.QueryRow(query, req.Name, req.Email, string(hashedPassword)).Scan(&userID)
	if err != nil {
		http.Error(w, "failed to register user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	resp := map[string]interface{}{
		"message": "user registered successfully",
		"user_id": userID,
		"email":   req.Email,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	// Ambil user dari DB
	var user struct {
		ID           int    `db:"id"`
		PasswordHash string `db:"password_hash"`
	}
	err := db.Get(&user, "SELECT id, password_hash FROM users WHERE email=$1", req.Email)
	if err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	// Cek password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	// Buat JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})
	signed, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		http.Error(w, "token generation failed", http.StatusInternalServerError)
		return
	}

	resp := map[string]string{"token": signed}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
