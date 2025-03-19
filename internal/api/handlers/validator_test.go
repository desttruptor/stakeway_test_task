package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"stakeway_test_task/internal/models"
	"testing"
)

type MockValidatorService struct {
	mock.Mock
}

func (m *MockValidatorService) CreateValidatorRequest(input *models.ValidatorRequestInput) (*models.ValidatorRequestResponse, error) {
	args := m.Called(input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.ValidatorRequestResponse), args.Error(1)
}

func (m *MockValidatorService) GetRequestStatus(requestID string) (*models.ValidatorStatusResponse, error) {
	args := m.Called(requestID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.ValidatorStatusResponse), args.Error(1)
}

func TestCreateValidator(t *testing.T) {
	t.Run("successful validator creation", func(t *testing.T) {
		mockService := new(MockValidatorService)

		expectedResponse := &models.ValidatorRequestResponse{
			RequestID: "test-uuid",
			Message:   "Validator creation in progress",
		}

		mockService.On("CreateValidatorRequest", mock.AnythingOfType("*models.ValidatorRequestInput")).
			Return(expectedResponse, nil)

		handler := &ValidatorHandler{service: mockService}

		requestBody := map[string]interface{}{
			"num_validators": 3,
			"fee_recipient":  "0x1234567890abcdef1234567890abcdef12345678",
		}
		body, _ := json.Marshal(requestBody)
		req := httptest.NewRequest(http.MethodPost, "/validators", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		handler.CreateValidator(w, req)

		assert.Equal(t, http.StatusAccepted, w.Code)

		var response models.ValidatorRequestResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, expectedResponse.RequestID, response.RequestID)
		assert.Equal(t, expectedResponse.Message, response.Message)

		mockService.AssertCalled(t, "CreateValidatorRequest", mock.AnythingOfType("*models.ValidatorRequestInput"))
	})

	t.Run("invalid request body", func(t *testing.T) {
		mockService := new(MockValidatorService)

		handler := &ValidatorHandler{service: mockService}

		req := httptest.NewRequest(http.MethodPost, "/validators", bytes.NewBuffer([]byte("invalid-json")))
		w := httptest.NewRecorder()

		handler.CreateValidator(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		mockService.AssertNotCalled(t, "CreateValidatorRequest", mock.Anything)
	})

	t.Run("service returns error", func(t *testing.T) {
		mockService := new(MockValidatorService)

		mockService.On("CreateValidatorRequest", mock.AnythingOfType("*models.ValidatorRequestInput")).
			Return(nil, errors.New("service error"))

		handler := &ValidatorHandler{service: mockService}

		requestBody := map[string]interface{}{
			"num_validators": 3,
			"fee_recipient":  "0x1234567890abcdef1234567890abcdef12345678",
		}
		body, _ := json.Marshal(requestBody)
		req := httptest.NewRequest(http.MethodPost, "/validators", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		handler.CreateValidator(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		mockService.AssertCalled(t, "CreateValidatorRequest", mock.AnythingOfType("*models.ValidatorRequestInput"))
	})
}

func TestGetValidatorStatus(t *testing.T) {
	t.Run("successful status retrieval", func(t *testing.T) {
		mockService := new(MockValidatorService)

		expectedResponse := &models.ValidatorStatusResponse{
			Status: models.StatusSuccessful,
			Keys:   []string{"key1", "key2"},
		}

		mockService.On("GetRequestStatus", "test-uuid").
			Return(expectedResponse, nil)

		handler := &ValidatorHandler{service: mockService}

		req := httptest.NewRequest(http.MethodGet, "/validators/test-uuid", nil)
		w := httptest.NewRecorder()

		vars := map[string]string{
			"request_id": "test-uuid",
		}
		req = mux.SetURLVars(req, vars)

		handler.GetValidatorStatus(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response models.ValidatorStatusResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, expectedResponse.Status, response.Status)
		assert.Equal(t, expectedResponse.Keys, response.Keys)

		mockService.AssertCalled(t, "GetRequestStatus", "test-uuid")
	})

	t.Run("request not found", func(t *testing.T) {
		mockService := new(MockValidatorService)

		mockService.On("GetRequestStatus", "non-existent-id").
			Return(nil, errors.New("request not found"))

		handler := &ValidatorHandler{service: mockService}

		req := httptest.NewRequest(http.MethodGet, "/validators/non-existent-id", nil)
		w := httptest.NewRecorder()

		vars := map[string]string{
			"request_id": "non-existent-id",
		}
		req = mux.SetURLVars(req, vars)

		handler.GetValidatorStatus(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)

		mockService.AssertCalled(t, "GetRequestStatus", "non-existent-id")
	})
}
