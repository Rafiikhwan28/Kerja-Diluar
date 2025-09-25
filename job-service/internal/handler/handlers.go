package handler

import (
	"encoding/json"
	"job-service/internal/model"
	"job-service/internal/repository"
	"net/http"
)

type JobHandler struct {
	repo *repository.JobRepo
}

func NewJobHandler(repo *repository.JobRepo) *JobHandler {
	return &JobHandler{repo: repo}
}

func (h *JobHandler) GetJobs(w http.ResponseWriter, r *http.Request) {
	jobs, err := h.repo.GetJobs()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jobs)
}

func (h *JobHandler) CreateJob(w http.ResponseWriter, r *http.Request) {
	var job model.Job
	if err := json.NewDecoder(r.Body).Decode(&job); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if job.Title == "" || job.Description == "" || job.Company == "" {
		http.Error(w, "title, description, and company required", http.StatusBadRequest)
		return
	}
	if err := h.repo.CreateJob(&job); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(job)
}
