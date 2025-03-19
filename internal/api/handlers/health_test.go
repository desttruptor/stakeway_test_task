package handlers

import (
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"stakeway_test_task/internal/mocks"
	"testing"
)

func TestHealthCheck(t *testing.T) {
	t.Run("successful health check", func(t *testing.T) {
		mockHealthChecker := mocks.NewHealthChecker(t)

		mockHealthChecker.On("CheckHealth").Return(nil)

		handler := NewHealthHandler(mockHealthChecker)

		req := httptest.NewRequest(http.MethodGet, "/health", nil)
		w := httptest.NewRecorder()

		handler.HealthCheck(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "up", response["status"])
		assert.Equal(t, "connected", response["database"])

		mockHealthChecker.AssertCalled(t, "CheckHealth")
	})

	t.Run("failed health check", func(t *testing.T) {
		mockHealthChecker := mocks.NewHealthChecker(t)

		mockHealthChecker.On("CheckHealth").Return(errors.New("database connection error"))

		handler := NewHealthHandler(mockHealthChecker)

		req := httptest.NewRequest(http.MethodGet, "/health", nil)
		w := httptest.NewRecorder()

		handler.HealthCheck(w, req)

		assert.Equal(t, http.StatusServiceUnavailable, w.Code)

		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "down", response["status"])
		assert.Equal(t, "disconnected", response["database"])

		mockHealthChecker.AssertCalled(t, "CheckHealth")
	})
}
