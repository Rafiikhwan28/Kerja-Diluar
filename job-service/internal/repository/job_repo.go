package repository

import (
	"job-service/internal/model"

	"github.com/jmoiron/sqlx"
)

type JobRepo struct {
	db *sqlx.DB
}

func NewJobRepo(db *sqlx.DB) *JobRepo {
	return &JobRepo{db: db}
}

func (r *JobRepo) GetJobs() ([]model.Job, error) {
	var jobs []model.Job
	err := r.db.Select(&jobs, "SELECT id, title, description, company FROM jobs")
	return jobs, err
}

func (r *JobRepo) CreateJob(job *model.Job) error {
	query := `INSERT INTO jobs (title, description, company) VALUES ($1, $2, $3) RETURNING id`
	return r.db.QueryRow(query, job.Title, job.Description, job.Company).Scan(&job.ID)
}
