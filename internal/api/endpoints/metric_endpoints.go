package endpoints

import (
	"math"
	"net/http"
	"time"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"

	"github.com/WendelHime/hydroponics-metrics-collector/internal/logic"
	localErrs "github.com/WendelHime/hydroponics-metrics-collector/internal/shared/errors"
	"github.com/WendelHime/hydroponics-metrics-collector/internal/shared/models"
)

type MetricsEndpoints struct {
	logic logic.MetricLogic
}

func NewMetricsEndpoints(logic logic.MetricLogic) MetricsEndpoints {
	return MetricsEndpoints{logic: logic}
}

type RegisterMetricRequest struct {
	Metrics []models.SensorRequest `json:"metrics" validate:"required,dive"`
}

func (s *RegisterMetricRequest) Bind(r *http.Request) error {
	validate := validator.New()
	err := validate.Struct(s)
	if err != nil {
		return localErrs.BadRequestErr.WithErr(err).WithMsg("failed to validate request")
	}

	for _, v := range s.Metrics {
		// parse timestamp
		sec, dec := math.Modf(v.Timestamp)
		v.Time = time.Unix(int64(sec), int64(dec*1e9))

		// we only store dates from the last 30 days
		if v.Time.Before(time.Now().AddDate(0, 0, -30)) {
			return localErrs.BadRequestErr.WithMsg("timestamp before 30 days is not acceptable")
		}
	}

	return nil
}

func (e MetricsEndpoints) RegisterMetric(w http.ResponseWriter, r *http.Request) {
	var request RegisterMetricRequest
	err := render.Bind(r, &request)
	if err != nil {
		log.Warn().Err(err).Msg("failed to decode sensor request")
		localErrs.RenderErr(w, r, err)
		return
	}

	err = e.logic.WriteSensorMetrics(r.Context(), request.Metrics)
	if err != nil {
		log.Error().Err(err).Msg("failed to write sensor metrics")
		localErrs.RenderErr(w, r, err)
		return
	}

	render.Status(r, http.StatusCreated)
}
