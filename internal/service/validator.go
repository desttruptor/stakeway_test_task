package services

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/google/uuid"
	"log/slog"
	"regexp"
	"stakeway_test_task/internal/models"
	"stakeway_test_task/internal/repository"
	"stakeway_test_task/internal/utils"
	"time"
)

type RequestRepo interface {
	CreateRequest(request *models.ValidatorRequest) error
	GetRequestByID(id string) (*models.ValidatorRequest, error)
	GetKeysByRequestID(requestID string) ([]string, error)
	UpdateRequestStatus(id string, status models.Status, errorMessage string) error
	SaveValidatorKey(key *models.ValidatorKey) error
}

type ValidatorService struct {
	repo   RequestRepo
	logger *slog.Logger
}

func NewValidatorService(repo *repository.ValidatorRepository, slog *slog.Logger) *ValidatorService {
	return &ValidatorService{repo: repo, logger: slog}
}

func (s *ValidatorService) CreateValidatorRequest(input *models.ValidatorRequestInput) (*models.ValidatorRequestResponse, error) {
	if input.NumValidators <= 0 {
		return nil, fmt.Errorf("number of validators must be positive")
	}

	if !isValidEthereumAddress(input.FeeRecipient) {
		return nil, fmt.Errorf("invalid Ethereum address format")
	}

	requestID := uuid.New().String()
	request := &models.ValidatorRequest{
		ID:            requestID,
		NumValidators: input.NumValidators,
		FeeRecipient:  input.FeeRecipient,
		Status:        models.StatusStarted,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	err := s.repo.CreateRequest(request)
	if err != nil {
		return nil, err
	}

	go s.processValidatorCreation(requestID, input.NumValidators, input.FeeRecipient)

	return &models.ValidatorRequestResponse{
		RequestID: requestID,
		Message:   "Validator creation in progress",
	}, nil
}

func (s *ValidatorService) GetRequestStatus(requestID string) (*models.ValidatorStatusResponse, error) {
	request, err := s.repo.GetRequestByID(requestID)
	if err != nil {
		return nil, err
	}

	response := &models.ValidatorStatusResponse{
		Status: request.Status,
	}

	if request.Status == models.StatusSuccessful {
		keys, err := s.repo.GetKeysByRequestID(requestID)
		if err != nil {
			return nil, err
		}
		response.Keys = keys
	} else if request.Status == models.StatusFailed {
		response.Message = request.ErrorMessage
	}

	return response, nil
}

func (s *ValidatorService) processValidatorCreation(requestID string, numValidators int, feeRecipient string) {
	s.logger.Info("Starting validator creation process",
		"request_id", requestID,
		"num_validators", numValidators)

	startTime := time.Now()

	utils.TasksTotal.WithLabelValues("started").Inc()

	var err error

	for i := 0; i < numValidators; i++ {
		time.Sleep(20 * time.Millisecond)

		key, err := generateRandomKey()
		if err != nil {
			s.logger.Error("Failed to generate key",
				"error", err,
				"request_id", requestID)

			err = s.repo.UpdateRequestStatus(requestID, models.StatusFailed, "Error generating validator keys")
			if err != nil {
				utils.TasksTotal.WithLabelValues("failed").Inc()
				s.logger.Error("Failed to update request status", "error", err)
			}
			return
		}

		validatorKey := &models.ValidatorKey{
			ID:           uuid.New().String(),
			RequestID:    requestID,
			Key:          key,
			FeeRecipient: feeRecipient,
		}

		err = s.repo.SaveValidatorKey(validatorKey)
		if err != nil {
			s.logger.Error("Failed to save validator key",
				"error", err,
				"request_id", requestID)

			err = s.repo.UpdateRequestStatus(requestID, models.StatusFailed, "Error saving validator keys")
			if err != nil {
				utils.TasksTotal.WithLabelValues("failed").Inc()
				s.logger.Error("Failed to update request status", "error", err)
			}
			return
		}

		s.logger.Info("Generated validator key",
			"key", key,
			"index", i+1,
			"request_id", requestID)
	}

	err = s.repo.UpdateRequestStatus(requestID, models.StatusSuccessful, "")
	if err != nil {
		s.logger.Error("Failed to update request status",
			"error", err,
			"request_id", requestID)
		utils.TasksTotal.WithLabelValues("failed").Inc()
	} else {
		s.logger.Info("Validator creation completed successfully",
			"request_id", requestID,
			"duration_ms", time.Since(startTime).Milliseconds())
		utils.TasksTotal.WithLabelValues("successful").Inc()
	}

	utils.TaskDuration.Observe(time.Since(startTime).Seconds())
}

func generateRandomKey() (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func isValidEthereumAddress(address string) bool {
	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	return re.MatchString(address)
}
