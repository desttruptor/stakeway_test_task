package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type HealthChecker interface {
	CheckHealth() error
}

type HealthHandler struct {
	repo HealthChecker
}

func NewHealthHandler(repo HealthChecker) *HealthHandler {
	return &HealthHandler{repo: repo}
}

func (h *HealthHandler) HealthCheck(w http.ResponseWriter, _ *http.Request) {
	status := map[string]string{
		"status": "up",
	}

	err := h.repo.CheckHealth()
	if err != nil {
		status["status"] = "down"
		status["database"] = "disconnected"
		w.WriteHeader(http.StatusServiceUnavailable)
	} else {
		status["database"] = "connected"
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(status)
	if err != nil {
		slog.Error("Can't encode health status", "error", err)
		return
	}
}
