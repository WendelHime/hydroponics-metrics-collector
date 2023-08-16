package logic

import (
	"context"
	"slices"

	localErrs "github.com/WendelHime/hydroponics-metrics-collector/internal/shared/errors"
	"github.com/WendelHime/hydroponics-metrics-collector/internal/shared/models"
	"github.com/WendelHime/hydroponics-metrics-collector/internal/storage"
)

type MetricLogic interface {
	WriteSensorMetrics(ctx context.Context, metrics []models.SensorRequest) error
}

type metricLogic struct {
	metricRepository     storage.MetricRepository
	userDeviceRepository storage.UserDeviceRepository
}

func NewMetricLogic(repository storage.MetricRepository, userDeviceRepository storage.UserDeviceRepository) MetricLogic {
	return &metricLogic{metricRepository: repository, userDeviceRepository: userDeviceRepository}
}

func (l *metricLogic) WriteSensorMetrics(ctx context.Context, m []models.SensorRequest) error {
	// caching inmem user devices
	userDevices := make(map[string][]string)
	for _, request := range m {
		if _, ok := userDevices[request.UserID]; !ok {
			devices, err := l.userDeviceRepository.GetDevicesFromUser(ctx, m[0].UserID)
			if err != nil {
				return err
			}
			userDevices[request.UserID] = devices
		}
		devices := userDevices[request.UserID]

		// checking if sensor ID is correlated to user devices
		if !slices.Contains(devices, request.SensorID) {
			return localErrs.ForbiddenErr
		}
	}

	// if everything succeed, write measurement
	return l.metricRepository.WriteMeasurement(ctx, m...)
}
