package rest

import (
	"context"
	"encoding/json"
	"html/template"
	"net/http"
	"strconv"

	"github.com/elangreza/scheduler/internal"
)

// RootHandler renders the create task form at the root page
func (h *Handler) RootHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/home.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func NewHandler(svc svc) *Handler {
	return &Handler{svc: svc}
}

type (
	svc interface {
		CreateTask(ctx context.Context, req internal.CreateTaskParams) error
		ListTask(ctx context.Context) ([]internal.Task, error)
		DeleteTask(ctx context.Context, id int) error
		UpdateTask(ctx context.Context, id int, req internal.UpdateTaskParams) error
	}

	Handler struct {
		svc
	}
)

// ListTaskHandler returns all tasks as JSON
func (h *Handler) ListTaskHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.svc.ListTask(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

// DeleteTaskHandler deletes a task by id (expects ?id=)
func (h *Handler) DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	if err := h.svc.DeleteTask(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// UpdateTaskHandler updates a task by id (expects ?id=, and JSON body)
func (h *Handler) UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	var req internal.UpdateTaskParams
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := h.svc.UpdateTask(r.Context(), id, req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) CreateTask(w http.ResponseWriter, r *http.Request) {
	// Support both API and form POST
	if r.Method == http.MethodPost && r.Header.Get("Content-Type") == "application/json" {
		var req internal.CreateTaskParams
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := h.svc.CreateTask(r.Context(), req); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		return
	}

	// Handle form POST
	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		req := internal.CreateTaskParams{
			Name:        r.FormValue("title"),
			Description: r.FormValue("description"),
		}
		if err := h.svc.CreateTask(r.Context(), req); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Render form
	tmpl, err := template.ParseFiles("templates/home.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}
