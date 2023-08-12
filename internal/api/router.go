package api

import (
	"github.com/WendelHime/hydroponics-metrics-collector/internal/api/endpoints"
	"github.com/WendelHime/hydroponics-metrics-collector/internal/api/middlewares"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog"
	"github.com/go-chi/render"
	"github.com/rs/zerolog"
)

func NewRouter(logger zerolog.Logger, metricsEndpoints endpoints.MetricsEndpoints, userEndpoints endpoints.UserEndpoints) chi.Router {
	mux := chi.NewRouter()
	mux.Use(httplog.RequestLogger(logger))
	mux.Use(render.SetContentType(render.ContentTypeJSON))

	// private endpoints for metrics
	mux.Group(func(r chi.Router) {
		r.Use(middlewares.EnsureValidToken)
		r.Use(middlewares.HasScope("write:metrics"))

		r.Post("/metrics", metricsEndpoints.RegisterMetric)
	})

	// private endpoints for binding user to sensor
	// TODO: create endpoints for binding user to sensor

	// public endpoints
	mux.Post("/users", userEndpoints.CreateAccount)
	mux.Post("/signin", userEndpoints.SignIn)

	return mux
}
