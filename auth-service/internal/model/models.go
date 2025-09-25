package model

import "time"

type User struct {
	ID           int       `db:"id" json:"id"`
	Name         string    `db:"name" json:"name"`
	Email        string    `db:"email" json:"email"`
	PasswordHash string    `db:"password_hash" json:"-"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
}

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
