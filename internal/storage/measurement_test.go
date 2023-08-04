package storage

import (
	"context"
	"errors"
	"testing"
	"time"

	localErrs "github.com/WendelHime/hydroponics-metrics-collector/internal/shared/errors"
	"github.com/WendelHime/hydroponics-metrics-collector/internal/shared/models"
	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"
)

func TestWriteMeasurement(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	now := time.Now()
	var tests = []struct {
		name             string
		assert           func(t *testing.T, err error)
		metricRepository func() MetricRepository
		givenRequests    []models.SensorRequest
	}{
		{
			name: "Creating multiple metrics with success",
			assert: func(t *testing.T, err error) {
				assert.Nil(t, err)
			},
			metricRepository: func() MetricRepository {
				db := "hydroponics"
				mock := NewMockInfluxClient(ctrl)
				mock.EXPECT().WriteData(gomock.Any(), db, []any{
					SensorMeasurement{
						Table:            "metrics",
						SensorID:         "test",
						SensorVersion:    "0.0.1",
						Alias:            "test",
						PH:               7.0,
						EC:               1403,
						TDS:              707,
						Humidity:         50.0,
						WaterTemperature: 25.0,
						Temperature:      25.0,
						Timestamp:        now,
					},
					SensorMeasurement{
						Table:            "metrics",
						SensorID:         "test2",
						SensorVersion:    "0.0.1",
						Alias:            "test2",
						PH:               7.0,
						EC:               1403,
						TDS:              707,
						Humidity:         50.0,
						WaterTemperature: 25.0,
						Temperature:      25.0,
						Timestamp:        now,
					},
				}).Return(nil)
				return NewRepository(db, mock)
			},
			givenRequests: []models.SensorRequest{
				{
					SensorID:         "test",
					SensorVersion:    "0.0.1",
					Alias:            "test",
					PH:               7.0,
					EC:               1403,
					TDS:              707,
					Humidity:         50.0,
					WaterTemperature: 25.0,
					Temperature:      25.0,
					Time:             now,
				},
				{
					SensorID:         "test2",
					SensorVersion:    "0.0.1",
					Alias:            "test2",
					PH:               7.0,
					EC:               1403,
					TDS:              707,
					Humidity:         50.0,
					WaterTemperature: 25.0,
					Temperature:      25.0,
					Time:             now,
				},
			},
		},
		{
			name: "Should return internal server error when there's unexpected error",
			assert: func(t *testing.T, err error) {
				var returnedErr *localErrs.Error
				if assert.ErrorAs(t, err, &returnedErr) {
					assert.Equal(t, 500, returnedErr.StatusCode)
					return
				}
				assert.Fail(t, "expected internal server error")
			},
			metricRepository: func() MetricRepository {
				db := "hydroponics"
				mock := NewMockInfluxClient(ctrl)
				mock.EXPECT().WriteData(gomock.Any(), db, gomock.Any()).Return(errors.New("random error"))
				return NewRepository(db, mock)
			},
			givenRequests: []models.SensorRequest{
				{
					SensorID:         "test",
					SensorVersion:    "0.0.1",
					Alias:            "test",
					PH:               7.0,
					EC:               1403,
					TDS:              707,
					Humidity:         50.0,
					WaterTemperature: 25.0,
					Temperature:      25.0,
					Time:             time.Now(),
				},
				{
					SensorID:         "test2",
					SensorVersion:    "0.0.1",
					Alias:            "test2",
					PH:               7.0,
					EC:               1403,
					TDS:              707,
					Humidity:         50.0,
					WaterTemperature: 25.0,
					Temperature:      25.0,
					Time:             time.Now(),
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			err := tt.metricRepository().WriteMeasurement(context.Background(), tt.givenRequests...)
			tt.assert(t, err)
		})
	}
}
