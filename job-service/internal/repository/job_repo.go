package repository

import (
	"database/sql"
	"errors"
	"job-service/internal/model"

	"github.com/jmoiron/sqlx"
)

type JobRepo struct {
	db *sqlx.DB
}

func NewJobRepo(db *sqlx.DB) *JobRepo {
	return &JobRepo{db: db}
}

// Get all jobs
func (r *JobRepo) GetJobs() ([]model.Job, error) {
	var jobs []model.Job
	err := r.db.Select(&jobs, `
		SELECT id, title, description, company, location, salary, category_id, created_by, created_at, updated_at
		FROM jobs
		ORDER BY id DESC
	`)
	return jobs, err
}

// Get job by ID
func (r *JobRepo) GetJobByID(id int) (*model.Job, error) {
	var job model.Job
	err := r.db.Get(&job, `
		SELECT id, title, description, company, location, salary, category_id, created_by, created_at, updated_at
		FROM jobs
		WHERE id = $1
	`, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil // tidak error fatal, tapi job tidak ditemukan
	}
	return &job, err
}

// Create job
func (r *JobRepo) CreateJob(job *model.Job) error {
	query := `
		INSERT INTO jobs (title, description, company, location, salary, category_id, created_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`
	// RETURNING cukup id + timestamp (category_id & created_by sudah dikirim dari input, tidak perlu dikembalikan lagi)
	return r.db.QueryRow(
		query,
		job.Title,
		job.Description,
		job.Company,
		job.Location,
		job.Salary,
		job.CategoryID,
		job.CreatedBy,
	).Scan(&job.ID, &job.CreatedAt, &job.UpdatedAt)
}

// Update job
func (r *JobRepo) UpdateJob(job *model.Job) error {
	query := `
		UPDATE jobs
		SET title = $1, description = $2, company = $3, location = $4, salary = $5, category_id = $6, updated_at = NOW()
		WHERE id = $7
		RETURNING updated_at
	`
	return r.db.QueryRow(
		query,
		job.Title,
		job.Description,
		job.Company,
		job.Location,
		job.Salary,
		job.CategoryID,
		job.ID,
	).Scan(&job.UpdatedAt)
}

// Delete job
func (r *JobRepo) DeleteJob(id int) error {
	_, err := r.db.Exec("DELETE FROM jobs WHERE id = $1", id)
	return err
}

// Get jobs by category ID
func (r *JobRepo) GetJobsByCategory(categoryID int) ([]model.Job, error) {
	var jobs []model.Job
	query := `
		SELECT id, title, description, company, location, salary, category_id, created_by, created_at, updated_at
		FROM jobs
		WHERE category_id = $1
		ORDER BY id DESC
	`
	err := r.db.Select(&jobs, query, categoryID)
	return jobs, err
}
