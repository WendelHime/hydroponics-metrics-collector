package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/auth0/go-auth0/authentication"
	"github.com/auth0/go-auth0/management"
	"github.com/go-chi/httplog"

	"github.com/InfluxCommunity/influxdb3-go/influx"
	"github.com/WendelHime/hydroponics-metrics-collector/internal/api"
	"github.com/WendelHime/hydroponics-metrics-collector/internal/api/endpoints"
	"github.com/WendelHime/hydroponics-metrics-collector/internal/logic"
	"github.com/WendelHime/hydroponics-metrics-collector/internal/services"
	"github.com/WendelHime/hydroponics-metrics-collector/internal/shared/errors"
	"github.com/WendelHime/hydroponics-metrics-collector/internal/storage"
)

func main() {
	database := os.Getenv("DATABASE")
	hostURL := os.Getenv("INFLUXDB_HOST")
	authToken := os.Getenv("INFLUXDB_TOKEN")
	auth0Domain := os.Getenv("AUTH0_DOMAIN")
	auth0ClientID := os.Getenv("AUTH0_CLIENTID")
	auth0ClientSecret := os.Getenv("AUTH0_CLIENT_SECRET")
	authAudience := os.Getenv("AUTH_AUDIENCE")
	authNonce := os.Getenv("AUTH0_NONCE")
	env := os.Getenv("ENV")
	roleID := os.Getenv("USER_ROLE_ID")
	projectID := os.Getenv("PROJECT_ID")

	ctx := context.Background()
	logger := httplog.NewLogger("hydroponics-metrics-collector", httplog.Options{
		LogLevel:        "info",
		LevelFieldName:  "level",
		JSON:            true,
		TimeFieldFormat: time.RFC3339Nano,
		TimeFieldName:   "timestamp",
	})

	influxCli, err := influx.New(influx.Configs{
		HostURL:   hostURL,
		AuthToken: authToken,
	})
	if err != nil {
		panic(errors.InternalServerErr.WithMsg("failed to create influx client").WithErr(err).Error())
	}

	metricsRepository := storage.NewRepository(database, influxCli)
	firestoreCli, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		panic(errors.InternalServerErr.WithMsg("failed to create firestore client").WithErr(err).Error())
	}

	userDeviceRepository := storage.NewUserDeviceRepository(firestoreCli)
	logger.Debug().Str("database", database).Str("hostURL", hostURL).Str("authToken", authToken).Msg("Creating repository with the environment variables")
	metricsLogic := logic.NewMetricLogic(metricsRepository, userDeviceRepository)
	metricsEndpoints := endpoints.NewMetricsEndpoints(metricsLogic)

	authCli, err := authentication.New(
		ctx,
		auth0Domain,
		authentication.WithClientID(auth0ClientID),
		authentication.WithClientSecret(auth0ClientSecret))
	if err != nil {
		panic(errors.InternalServerErr.WithMsg("failed to create auth0 authentication client").WithDetails("err", err.Error()).Error())
	}

	managementCli, err := management.New(auth0Domain, management.WithClientCredentials(ctx, auth0ClientID, auth0ClientSecret))
	if err != nil {
		panic(errors.InternalServerErr.WithMsg("failed to create auth0 management client").WithDetails("err", err.Error()).Error())
	}

	authService := services.NewAuthService(authCli.OAuth, env, authAudience, authNonce)
	userService := services.NewUserService(managementCli.User, managementCli.Role)

	userLogic := logic.NewUserLogic(userService, authService, userDeviceRepository, roleID)
	userEndpoints := endpoints.NewUserEndpoint(userLogic)
	r := api.NewRouter(logger, metricsEndpoints, userEndpoints, authNonce)

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
			influxCli.Close()
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				logger.Fatal().Msg("graceful shutdown timed out.. forcing exit.")
			}
		}()

		// Trigger graceful shutdown
		err = server.Shutdown(shutdownCtx)
		if err != nil {
			logger.Fatal().Err(err)
		}
		serverStopCtx()
	}()

	// Run the server
	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		logger.Fatal().Err(err)
	}

	// Wait for server context to be stopped
	<-serverCtx.Done()
}
