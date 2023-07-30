package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog"
	"github.com/go-chi/render"
	"github.com/rs/zerolog"
)

func NewRouter(logger zerolog.Logger, metricsEndpoints MetricsEndpoints, userEndpoints UserEndpoints) chi.Router {
	mux := chi.NewRouter()
	mux.Use(httplog.RequestLogger(logger))
	mux.Use(render.SetContentType(render.ContentTypeJSON))

	mux.Post("/metrics", metricsEndpoints.RegisterMetric)
	mux.Post("/users", userEndpoints.CreateAccount)
	mux.Post("/signin", userEndpoints.SignIn)

	return mux
}
