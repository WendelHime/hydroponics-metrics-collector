package logic

import (
	"context"
	"testing"

	localErrs "github.com/WendelHime/hydroponics-metrics-collector/internal/shared/errors"
	"github.com/WendelHime/hydroponics-metrics-collector/internal/shared/models"
	"github.com/WendelHime/hydroponics-metrics-collector/internal/storage"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"
)

func TestWriteSensorMetrics(t *testing.T) {
	userID := uuid.NewString()
	device1 := uuid.NewString()
	device2 := uuid.NewString()
	var tests = []struct {
		name         string
		setup        func(ctrl *gomock.Controller) MetricLogic
		givenMetrics []models.SensorRequest
		assert       func(t *testing.T, err error)
	}{
		{
			name: "write metric with success",
			setup: func(ctrl *gomock.Controller) MetricLogic {
				userDeviceRepository := storage.NewMockUserDeviceRepository(ctrl)
				userDeviceRepository.EXPECT().GetDevicesFromUser(gomock.Any(), userID).Return([]string{device1, device2}, nil).Times(1)
				metricRepository := storage.NewMockMetricRepository(ctrl)
				metricRepository.EXPECT().WriteMeasurement(gomock.Any(), []models.SensorRequest{
					{SensorID: device1, UserID: userID},
				}).Return(nil).Times(1)
				return NewMetricLogic(metricRepository, userDeviceRepository)
			},
			givenMetrics: []models.SensorRequest{{SensorID: device1, UserID: userID}},
			assert: func(t *testing.T, err error) {
				assert.Nil(t, err)
			},
		},
		{
			name: "user doesn't have any devices",
			setup: func(ctrl *gomock.Controller) MetricLogic {
				userDeviceRepository := storage.NewMockUserDeviceRepository(ctrl)
				userDeviceRepository.EXPECT().GetDevicesFromUser(gomock.Any(), userID).Return([]string{}, localErrs.NotFoundErr).Times(1)
				return NewMetricLogic(nil, userDeviceRepository)
			},
			givenMetrics: []models.SensorRequest{{SensorID: device1, UserID: userID}},
			assert: func(t *testing.T, err error) {
				if assert.Error(t, err) {
					assert.ErrorIs(t, err, localErrs.NotFoundErr)
				}
			},
		},
		{
			name: "provided device isn't correlated to the user",
			setup: func(ctrl *gomock.Controller) MetricLogic {
				userDeviceRepository := storage.NewMockUserDeviceRepository(ctrl)
				userDeviceRepository.EXPECT().GetDevicesFromUser(gomock.Any(), userID).Return([]string{device1}, nil).Times(1)
				return NewMetricLogic(nil, userDeviceRepository)
			},
			givenMetrics: []models.SensorRequest{{SensorID: device2, UserID: userID}},
			assert: func(t *testing.T, err error) {
				if assert.Error(t, err) {
					assert.ErrorIs(t, err, localErrs.ForbiddenErr)
				}
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			logic := tt.setup(ctrl)
			err := logic.WriteSensorMetrics(context.Background(), tt.givenMetrics)
			tt.assert(t, err)
		})
	}
}
