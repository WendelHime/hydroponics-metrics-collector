package api

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"

	"github.com/WendelHime/hydroponics-metrics-collector/internal/logic"
	"github.com/WendelHime/hydroponics-metrics-collector/internal/shared/models"
)

type MetricsEndpoints struct {
	logic logic.Logic
}

func NewMetricsEndpoints(logic logic.Logic) MetricsEndpoints {
	return MetricsEndpoints{logic: logic}
}

func (e MetricsEndpoints) RegisterMetric(w http.ResponseWriter, r *http.Request) {
	var sensorRequest models.SensorRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&sensorRequest)
	if err != nil {
		log.Warn().Err(err).Msg("failed to decode sensor request")
		http.Error(w, "failed to decode request", http.StatusBadRequest)
		return
	}

	err = e.logic.WriteSensorMetrics(r.Context(), sensorRequest)
	if err != nil {
		// TODO add error package and validate returned errors
		log.Error().Err(err).Msg("failed to write sensor metrics")
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(200)
}
