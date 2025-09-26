package model

import (
	"database/sql"
	"encoding/json"
	"time"
)

type Job struct {
	ID          int       `db:"id" json:"id"`
	Title       string    `db:"title" json:"title"`
	Description string    `db:"description" json:"description"`
	Company     string    `db:"company" json:"company"`
	Location    *string   `db:"location" json:"location"`
	Salary      *string   `db:"salary" json:"salary"`
	CategoryID  *int      `db:"category_id" json:"category_id,omitempty"`
	CreatedBy   *int      `db:"created_by" json:"created_by,omitempty"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

// NullString wrapper biar JSON rapi
type NullString struct {
	sql.NullString
}

func (ns NullString) MarshalJSON() ([]byte, error) {
	if ns.Valid {
		return json.Marshal(ns.String)
	}
	return json.Marshal(nil) // kalau NULL -> tampil null di JSON
}
