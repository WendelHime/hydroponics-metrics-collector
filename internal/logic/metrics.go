package logic

import (
	"context"
	"math"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"

	"github.com/WendelHime/hydroponics-metrics-collector/internal/shared/models"
	"github.com/WendelHime/hydroponics-metrics-collector/internal/storage"
)

type Logic interface {
	WriteSensorMetrics(ctx context.Context, m models.SensorRequest) error
}

type logic struct {
	repository storage.Repository
}

func NewLogic(repository storage.Repository) Logic {
	return &logic{repository: repository}
}

func (l logic) WriteSensorMetrics(ctx context.Context, m models.SensorRequest) error {
	validate := validator.New()
	err := validate.Struct(m)
	if err != nil {
		log.Warn().Err(err).Msg("failed to validate request")
		return err
	}

	// parse timestamp
	sec, dec := math.Modf(m.Timestamp)
	m.Time = time.Unix(int64(sec), int64(dec*1e9))

	// TODO rewrite this function to accept a list of requests instead of a single one, we should be able to aggregate requests in the future
	return l.repository.WriteMeasurement(ctx, m)
}
