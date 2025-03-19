package models

import (
	"time"
)

type Status string

const (
	StatusStarted    Status = "started"
	StatusSuccessful Status = "successful"
	StatusFailed     Status = "failed"
)

type ValidatorRequest struct {
	ID            string    `json:"request_id"`
	NumValidators int       `json:"num_validators"`
	FeeRecipient  string    `json:"fee_recipient"`
	Status        Status    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	ErrorMessage  string    `json:"error_message,omitempty"`
}

type ValidatorKey struct {
	ID           string `json:"id"`
	RequestID    string `json:"request_id"`
	Key          string `json:"key"`
	FeeRecipient string `json:"fee_recipient"`
}

type ValidatorRequestInput struct {
	NumValidators int    `json:"num_validators"`
	FeeRecipient  string `json:"fee_recipient"`
}

type ValidatorRequestResponse struct {
	RequestID string `json:"request_id"`
	Message   string `json:"message"`
}

type ValidatorStatusResponse struct {
	Status  Status   `json:"status"`
	Keys    []string `json:"keys,omitempty"`
	Message string   `json:"message,omitempty"`
}
