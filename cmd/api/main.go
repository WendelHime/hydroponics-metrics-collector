package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/httplog"

	"github.com/InfluxCommunity/influxdb3-go/influx"
	"github.com/WendelHime/hydroponics-metrics-collector/internal/api"
	"github.com/WendelHime/hydroponics-metrics-collector/internal/logic"
	"github.com/WendelHime/hydroponics-metrics-collector/internal/storage"
)

func main() {
	database := os.Getenv("DATABASE")
	hostURL := os.Getenv("INFLUXDB_HOST")
	authToken := os.Getenv("INFLUXDB_TOKEN")

	logger := httplog.NewLogger("hydroponics-metrics-collector", httplog.Options{
		LogLevel:        "info",
		LevelFieldName:  "level",
		JSON:            true,
		TimeFieldFormat: time.RFC3339Nano,
		TimeFieldName:   "timestamp",
	})

	repository := storage.NewRepository(database, influx.Configs{
		HostURL:   hostURL,
		AuthToken: authToken,
	})

	logger.Debug().Str("database", database).Str("hostURL", hostURL).Str("authToken", authToken).Msg("Creating repository with the environment variables")
	l := logic.NewLogic(repository)
	endpoints := api.NewMetricsEndpoints(l)
	r := api.NewRouter(logger, endpoints)

	server := &http.Server{Addr: "0.0.0.0:8080", Handler: r}

	// Server run context
	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	// Listen for syscall signals for process to interrupt/quit
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		logger.Info().Msg("Service started listening at port 8080")
		<-sig

		logger.Info().Msg("Received graceful shutdown signal")

		// Shutdown signal with grace period of 30 seconds
		shutdownCtx, cancel := context.WithTimeout(serverCtx, 30*time.Second)
		defer cancel()

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				logger.Fatal().Msg("graceful shutdown timed out.. forcing exit.")
			}
		}()

		// Trigger graceful shutdown
		err := server.Shutdown(shutdownCtx)
		if err != nil {
			logger.Fatal().Err(err)
		}
		serverStopCtx()
	}()

	// Run the server
	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		logger.Fatal().Err(err)
	}

	// Wait for server context to be stopped
	<-serverCtx.Done()
}
