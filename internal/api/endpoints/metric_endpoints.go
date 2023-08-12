package endpoints

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/rs/zerolog/log"

	"github.com/WendelHime/hydroponics-metrics-collector/internal/logic"
	"github.com/WendelHime/hydroponics-metrics-collector/internal/shared/errors"
	"github.com/WendelHime/hydroponics-metrics-collector/internal/shared/models"
)

type MetricsEndpoints struct {
	logic logic.MetricLogic
}

func NewMetricsEndpoints(logic logic.MetricLogic) MetricsEndpoints {
	return MetricsEndpoints{logic: logic}
}

func (e MetricsEndpoints) RegisterMetric(w http.ResponseWriter, r *http.Request) {
	var sensorRequest models.SensorRequest
	err := render.Bind(r, &sensorRequest)
	if err != nil {
		log.Warn().Err(err).Msg("failed to decode sensor request")
		errors.RenderErr(w, r, err)
		return
	}

	err = e.logic.WriteSensorMetrics(r.Context(), sensorRequest)
	if err != nil {
		log.Error().Err(err).Msg("failed to write sensor metrics")
		errors.RenderErr(w, r, err)
		return
	}

	render.Status(r, http.StatusCreated)
}
