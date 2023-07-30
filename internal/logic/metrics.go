package logic

import (
	"context"
	"math"
	"time"

	"github.com/WendelHime/hydroponics-metrics-collector/internal/shared/models"
	"github.com/WendelHime/hydroponics-metrics-collector/internal/storage"
)

type MetricLogic interface {
	WriteSensorMetrics(ctx context.Context, m models.SensorRequest) error
}

type metricLogic struct {
	repository storage.MetricRepository
}

func NewLogic(repository storage.MetricRepository) MetricLogic {
	return &metricLogic{repository: repository}
}

func (l *metricLogic) WriteSensorMetrics(ctx context.Context, m models.SensorRequest) error {
	// parse timestamp
	sec, dec := math.Modf(m.Timestamp)
	m.Time = time.Unix(int64(sec), int64(dec*1e9))

	return l.repository.WriteMeasurement(ctx, m)
}
