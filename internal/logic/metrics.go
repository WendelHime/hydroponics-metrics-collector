package logic

import (
	"context"

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
	return l.repository.WriteMeasurement(ctx, m)
}
