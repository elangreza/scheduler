package rest

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/elangreza/scheduler/internal"
)

func NewHandler(svc svc) *Handler {
	return &Handler{svc: svc}
}

type (
	svc interface {
		CreateTask(ctx context.Context, req internal.CreateTaskParams) error
	}

	Handler struct {
		svc
	}
)

func (h *Handler) CreateTask(w http.ResponseWriter, r *http.Request) {
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
}
