package services

import (
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"log/slog"
	"os"
	"stakeway_test_task/internal/mocks"
	"stakeway_test_task/internal/models"
	"testing"
	"time"
)

func setupValidatorServiceTest(t *testing.T) (*mocks.RequestRepo, *ValidatorService) {
	mockRepo := mocks.NewRequestRepo(t)

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	service := &ValidatorService{
		repo:   mockRepo,
		logger: logger,
	}

	return mockRepo, service
}

func TestCreateValidatorRequest(t *testing.T) {
	t.Run("successful request creation", func(t *testing.T) {
		mockRepo, service := setupValidatorServiceTest(t)

		mockRepo.On("CreateRequest", mock.AnythingOfType("*models.ValidatorRequest")).
			Return(nil)

		input := &models.ValidatorRequestInput{
			NumValidators: 3,
			FeeRecipient:  "0x1234567890abcdef1234567890abcdef12345678",
		}

		response, err := service.CreateValidatorRequest(input)

		assert.NoError(t, err)
		assert.NotEmpty(t, response.RequestID)
		assert.Equal(t, "Validator creation in progress", response.Message)

		mockRepo.AssertCalled(t, "CreateRequest", mock.AnythingOfType("*models.ValidatorRequest"))

		time.Sleep(10 * time.Millisecond)
	})

	t.Run("validation error - negative validators", func(t *testing.T) {
		_, service := setupValidatorServiceTest(t)

		input := &models.ValidatorRequestInput{
			NumValidators: 0,
			FeeRecipient:  "0x1234567890abcdef1234567890abcdef12345678",
		}

		response, err := service.CreateValidatorRequest(input)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Contains(t, err.Error(), "number of validators must be positive")
	})

	t.Run("validation error - invalid ethereum address", func(t *testing.T) {
		_, service := setupValidatorServiceTest(t)

		input := &models.ValidatorRequestInput{
			NumValidators: 3,
			FeeRecipient:  "invalid-address",
		}

		response, err := service.CreateValidatorRequest(input)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Contains(t, err.Error(), "invalid Ethereum address format")
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo, service := setupValidatorServiceTest(t)

		mockRepo.On("CreateRequest", mock.AnythingOfType("*models.ValidatorRequest")).
			Return(errors.New("database error"))

		input := &models.ValidatorRequestInput{
			NumValidators: 3,
			FeeRecipient:  "0x1234567890abcdef1234567890abcdef12345678",
		}

		response, err := service.CreateValidatorRequest(input)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Contains(t, err.Error(), "database error")
	})
}

func TestGetRequestStatus(t *testing.T) {
	t.Run("successful status retrieval", func(t *testing.T) {
		mockRepo, service := setupValidatorServiceTest(t)

		requestID := uuid.New().String()

		mockRepo.On("GetRequestByID", requestID).
			Return(&models.ValidatorRequest{
				ID:            requestID,
				NumValidators: 2,
				FeeRecipient:  "0x1234567890abcdef1234567890abcdef12345678",
				Status:        models.StatusSuccessful,
				CreatedAt:     time.Now(),
				UpdatedAt:     time.Now(),
			}, nil)

		mockRepo.On("GetKeysByRequestID", requestID).
			Return([]string{"key1", "key2"}, nil)

		response, err := service.GetRequestStatus(requestID)

		assert.NoError(t, err)
		assert.Equal(t, models.StatusSuccessful, response.Status)
		assert.Equal(t, []string{"key1", "key2"}, response.Keys)
		assert.Empty(t, response.Message)
	})

	t.Run("failed status retrieval", func(t *testing.T) {
		mockRepo, service := setupValidatorServiceTest(t)

		requestID := uuid.New().String()
		errorMessage := "Test error message"

		mockRepo.On("GetRequestByID", requestID).
			Return(&models.ValidatorRequest{
				ID:            requestID,
				NumValidators: 2,
				FeeRecipient:  "0x1234567890abcdef1234567890abcdef12345678",
				Status:        models.StatusFailed,
				CreatedAt:     time.Now(),
				UpdatedAt:     time.Now(),
				ErrorMessage:  errorMessage,
			}, nil)

		response, err := service.GetRequestStatus(requestID)

		assert.NoError(t, err)
		assert.Equal(t, models.StatusFailed, response.Status)
		assert.Empty(t, response.Keys)
		assert.Equal(t, errorMessage, response.Message)
	})

	t.Run("request not found", func(t *testing.T) {
		mockRepo, service := setupValidatorServiceTest(t)

		requestID := "non-existent-id"

		mockRepo.On("GetRequestByID", requestID).
			Return(nil, errors.New("request not found"))

		response, err := service.GetRequestStatus(requestID)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Contains(t, err.Error(), "request not found")
	})

	t.Run("error getting keys", func(t *testing.T) {
		mockRepo, service := setupValidatorServiceTest(t)

		requestID := uuid.New().String()

		mockRepo.On("GetRequestByID", requestID).
			Return(&models.ValidatorRequest{
				ID:            requestID,
				NumValidators: 2,
				FeeRecipient:  "0x1234567890abcdef1234567890abcdef12345678",
				Status:        models.StatusSuccessful,
				CreatedAt:     time.Now(),
				UpdatedAt:     time.Now(),
			}, nil)

		mockRepo.On("GetKeysByRequestID", requestID).
			Return(nil, errors.New("database error"))

		response, err := service.GetRequestStatus(requestID)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Contains(t, err.Error(), "database error")
	})
}

func TestProcessValidatorCreation(t *testing.T) {
	t.Run("successful key generation", func(t *testing.T) {
		mockRepo, service := setupValidatorServiceTest(t)

		requestID := uuid.New().String()
		numValidators := 2
		feeRecipient := "0x1234567890abcdef1234567890abcdef12345678"

		mockRepo.On("SaveValidatorKey", mock.AnythingOfType("*models.ValidatorKey")).
			Return(nil)

		mockRepo.On("UpdateRequestStatus", requestID, models.StatusSuccessful, "").
			Return(nil)

		service.processValidatorCreation(requestID, numValidators, feeRecipient)

		mockRepo.AssertNumberOfCalls(t, "SaveValidatorKey", numValidators)
		mockRepo.AssertCalled(t, "UpdateRequestStatus", requestID, models.StatusSuccessful, "")
	})

	t.Run("error saving validator key", func(t *testing.T) {
		mockRepo, service := setupValidatorServiceTest(t)

		requestID := uuid.New().String()
		numValidators := 2
		feeRecipient := "0x1234567890abcdef1234567890abcdef12345678"

		mockRepo.On("SaveValidatorKey", mock.AnythingOfType("*models.ValidatorKey")).
			Return(errors.New("database error"))

		mockRepo.On("UpdateRequestStatus", requestID, models.StatusFailed, mock.Anything).
			Return(nil)

		service.processValidatorCreation(requestID, numValidators, feeRecipient)

		mockRepo.AssertNumberOfCalls(t, "SaveValidatorKey", 1) // только первая попытка
		mockRepo.AssertCalled(t, "UpdateRequestStatus", requestID, models.StatusFailed, mock.Anything)
	})
}
