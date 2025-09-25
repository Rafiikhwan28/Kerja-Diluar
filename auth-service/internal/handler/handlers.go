package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"kerjadiluar/auth/internal/model"
	"kerjadiluar/auth/internal/repository"
)

type Handler struct {
	repo      repository.UserRepo
	jwtSecret string
}

func NewHandler(r repository.UserRepo, jwtSecret string) *Handler {
	if jwtSecret == "" {
		jwtSecret = "dev-secret"
	}
	return &Handler{repo: r, jwtSecret: jwtSecret}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Decode request body
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// Validasi input
	if in.Name == "" || in.Email == "" || in.Password == "" {
		http.Error(w, "name, email, and password are required", http.StatusBadRequest)
		return
	}

	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "failed to hash password", http.StatusInternalServerError)
		return
	}

	// Buat user object
	u := &model.User{
		Name:         in.Name,
		Email:        in.Email,
		PasswordHash: string(hash),
	}

	// Simpan ke DB lewat repo
	if err := h.repo.CreateUser(u); err != nil {
		http.Error(w, "failed to create user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Response sukses
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"id":         u.ID,
		"name":       u.Name,
		"email":      u.Email,
		"created_at": u.CreatedAt,
	})
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var creds model.Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	u, err := h.repo.GetByEmail(creds.Email)
	if err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(creds.Password)); err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": u.ID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})
	signed, err := token.SignedString([]byte(h.jwtSecret))
	if err != nil {
		http.Error(w, "token error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"token": signed})
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, _ := strconv.Atoi(idStr)
	u, err := h.repo.GetByID(id)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(u)
}

/* AuthMiddleware - validate JWT and set user id header */
func (h *Handler) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth == "" {
			http.Error(w, "missing auth", http.StatusUnauthorized)
			return
		}
		parts := strings.Split(auth, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "bad auth", http.StatusUnauthorized)
			return
		}
		tokenStr := parts[1]
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) { return []byte(h.jwtSecret), nil })
		if err != nil || !token.Valid {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}
		claims := token.Claims.(jwt.MapClaims)
		if uidf, ok := claims["user_id"].(float64); ok {
			r.Header.Set("X-User-ID", strconv.Itoa(int(uidf)))
			next(w, r)
			return
		}
		http.Error(w, "invalid token claims", http.StatusUnauthorized)
	}
}
