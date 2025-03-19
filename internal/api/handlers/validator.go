package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"stakeway_test_task/internal/models"
)

type Validator interface {
	CreateValidatorRequest(input *models.ValidatorRequestInput) (*models.ValidatorRequestResponse, error)
	GetRequestStatus(requestID string) (*models.ValidatorStatusResponse, error)
}

type ValidatorHandler struct {
	service Validator
}

func NewValidatorHandler(service Validator) *ValidatorHandler {
	return &ValidatorHandler{service: service}
}

func (h *ValidatorHandler) CreateValidator(w http.ResponseWriter, r *http.Request) {
	var input models.ValidatorRequestInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	response, err := h.service.CreateValidatorRequest(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *ValidatorHandler) GetValidatorStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	requestID := vars["request_id"]

	status, err := h.service.GetRequestStatus(requestID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(status)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
