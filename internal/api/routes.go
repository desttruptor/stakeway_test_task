package api

import (
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log/slog"
	"stakeway_test_task/internal/api/handlers"
	"stakeway_test_task/internal/api/middleware"
	"stakeway_test_task/internal/repository"
	services "stakeway_test_task/internal/service"
)

func SetupRoutes(repo *repository.ValidatorRepository, logger *slog.Logger) *mux.Router {
	r := mux.NewRouter()

	// services
	validatorService := services.NewValidatorService(repo, logger)

	// handlers
	validatorHandler := handlers.NewValidatorHandler(validatorService)
	healthHandler := handlers.NewHealthHandler(repo)

	// middleware
	r.Use(middleware.MetricsMiddleware)
	r.Use(middleware.LoggingMiddleware(logger))

	// routes
	r.HandleFunc("/validators", validatorHandler.CreateValidator).Methods("POST")
	r.HandleFunc("/validators/{request_id}", validatorHandler.GetValidatorStatus).Methods("GET")
	r.HandleFunc("/health", healthHandler.HealthCheck).Methods("GET")

	r.Handle("/metrics", promhttp.Handler())

	return r
}
