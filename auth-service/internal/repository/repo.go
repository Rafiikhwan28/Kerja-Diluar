package repository

import (
	"kerjadiluar/auth/internal/model"

	"github.com/jmoiron/sqlx"
)

type UserRepo interface {
	CreateUser(u *model.User) error
	GetByEmail(email string) (*model.User, error)
	GetByID(id int) (*model.User, error)
}

type userRepo struct {
	db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) UserRepo {
	return &userRepo{db: db}
}

func (r *userRepo) CreateUser(u *model.User) error {
	return r.db.QueryRowx(
		`INSERT INTO users (name,email,password_hash) VALUES ($1,$2,$3) RETURNING id, created_at`,
		u.Name, u.Email, u.PasswordHash,
	).Scan(&u.ID, &u.CreatedAt)
}

func (r *userRepo) GetByEmail(email string) (*model.User, error) {
	var u model.User
	err := r.db.Get(&u, "SELECT * FROM users WHERE email=$1", email)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *userRepo) GetByID(id int) (*model.User, error) {
	var u model.User
	err := r.db.Get(&u, "SELECT id,name,email,created_at FROM users WHERE id=$1", id)
	if err != nil {
		return nil, err
	}
	return &u, nil
}
