package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog"
	"github.com/rs/zerolog"
)

func NewRouter(logger zerolog.Logger, metricsEndpoints MetricsEndpoints) chi.Router {
	mux := chi.NewRouter()
	mux.Use(httplog.RequestLogger(logger))
	mux.Post("/metrics", metricsEndpoints.RegisterMetric)
	return mux
}
