package api

import (
	"github.com/WendelHime/hydroponics-metrics-collector/internal/api/endpoints"
	"github.com/WendelHime/hydroponics-metrics-collector/internal/api/middlewares"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog"
	"github.com/go-chi/render"
	"github.com/rs/zerolog"
)

func NewRouter(logger zerolog.Logger, metricsEndpoints endpoints.MetricsEndpoints, userEndpoints endpoints.UserEndpoints, nonce string) chi.Router {
	mux := chi.NewRouter()
	mux.Use(httplog.RequestLogger(logger))
	mux.Use(render.SetContentType(render.ContentTypeJSON))

	// public endpoints
	mux.Post("/users", userEndpoints.CreateAccount)
	mux.Post("/signin", userEndpoints.SignIn)

	// private endpoints for iot device
	mux.Group(func(r chi.Router) {
		r.Use(middlewares.EnsureValidToken)
		r.Use(middlewares.HasScope("write:metrics"))

		r.Post("/metrics", metricsEndpoints.RegisterMetric)
	})

	// private endpoints for binding user to device
	mux.Group(func(r chi.Router) {
		r.Use(middlewares.EnsureValidToken)
		r.Use(middlewares.HasScope("write:device read:device"))
		r.Use(middlewares.UserMatches)

		r.Post("/users/{userID}/devices", userEndpoints.AddDevice)
		r.Get("/users/{userID}/devices", userEndpoints.GetDevices)
	})

	return mux
}
