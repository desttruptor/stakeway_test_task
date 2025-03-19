package middleware

import (
	"github.com/gorilla/mux"
	"net/http"
	"stakeway_test_task/internal/utils"
	"strconv"
	"time"
)

func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rw := &responseWriter{w, http.StatusOK}

		next.ServeHTTP(rw, r)

		route := mux.CurrentRoute(r)
		path, _ := route.GetPathTemplate()

		duration := time.Since(start).Seconds()
		utils.RequestsTotal.WithLabelValues(path, r.Method, strconv.Itoa(rw.statusCode)).Inc()
		utils.RequestDuration.WithLabelValues(path, r.Method).Observe(duration)
	})
}
