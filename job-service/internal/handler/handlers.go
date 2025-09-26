package handler

import (
	"encoding/json"
	"job-service/internal/model"
	"job-service/internal/repository"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type JobHandler struct {
	repo *repository.JobRepo
}

func NewJobHandler(repo *repository.JobRepo) *JobHandler {
	return &JobHandler{repo: repo}
}

// GET /jobs
func (h *JobHandler) GetJobs(w http.ResponseWriter, r *http.Request) {
	jobs, err := h.repo.GetJobs()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, jobs)
}

// GET /jobs/{id}
func (h *JobHandler) GetJobByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid job id", http.StatusBadRequest)
		return
	}
	job, err := h.repo.GetJobByID(id)
	if err != nil {
		http.Error(w, "job not found", http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, job)
}

// POST /jobs
func (h *JobHandler) CreateJob(w http.ResponseWriter, r *http.Request) {
	var job model.Job
	if err := json.NewDecoder(r.Body).Decode(&job); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if job.Title == "" || job.Description == "" || job.Company == "" {
		http.Error(w, "title, description, and company are required", http.StatusBadRequest)
		return
	}
	if err := h.repo.CreateJob(&job); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusCreated, job)
}

// PUT /jobs/{id}
func (h *JobHandler) UpdateJob(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid job id", http.StatusBadRequest)
		return
	}

	var job model.Job
	if err := json.NewDecoder(r.Body).Decode(&job); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	job.ID = id

	if err := h.repo.UpdateJob(&job); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, job)
}

// DELETE /jobs/{id}
func (h *JobHandler) DeleteJob(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid job id", http.StatusBadRequest)
		return
	}
	if err := h.repo.DeleteJob(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "job deleted"})
}

// GET /jobs/category/{categoryID}
func (h *JobHandler) GetJobsByCategory(w http.ResponseWriter, r *http.Request) {
	categoryIDStr := chi.URLParam(r, "categoryID")
	categoryID, err := strconv.Atoi(categoryIDStr)
	if err != nil {
		http.Error(w, "invalid category id", http.StatusBadRequest)
		return
	}

	jobs, err := h.repo.GetJobsByCategory(categoryID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jobs)
}

// helper
func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
