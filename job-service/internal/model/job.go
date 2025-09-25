package model

type Job struct {
	ID          int    `db:"id" json:"id"`
	Title       string `db:"title" json:"title"`
	Description string `db:"description" json:"description"`
	Company     string `db:"company" json:"company"`
}
