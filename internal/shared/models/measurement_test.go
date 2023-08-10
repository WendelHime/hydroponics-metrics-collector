package models

import (
	"testing"
	"time"

	localErrs "github.com/WendelHime/hydroponics-metrics-collector/internal/shared/errors"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestSensorRequestBind(t *testing.T) {
	var tests = []struct {
		name               string
		assert             func(t *testing.T, err error)
		givenSensorRequest *SensorRequest
	}{
		{
			name: "bind with success",
			assert: func(t *testing.T, err error) {
				assert.Nil(t, err)
			},
			givenSensorRequest: &SensorRequest{
				SensorID:         uuid.NewString(),
				SensorVersion:    "1.0.0",
				Alias:            "lettuce 1",
				Temperature:      0.0,
				Humidity:         0.0,
				PH:               0.0,
				TDS:              0.0,
				EC:               0.0,
				WaterTemperature: 0.0,
				Timestamp:        float64(time.Now().Unix()),
				Time:             time.Time{},
			},
		},
		{
			name: "missing any required field should return a bad request",
			assert: func(t *testing.T, err error) {
				if assert.Error(t, err) {
					assert.ErrorIs(t, err, localErrs.BadRequestErr)
				}
			},
			givenSensorRequest: &SensorRequest{
				SensorID:         "",
				SensorVersion:    "1.0.0",
				Alias:            "lettuce 1",
				Temperature:      0.0,
				Humidity:         0.0,
				PH:               0.0,
				TDS:              0.0,
				EC:               0.0,
				WaterTemperature: 0.0,
				Timestamp:        float64(time.Now().Unix()),
				Time:             time.Time{},
			},
		},
		{
			name: "timestamp before 30 days should return an error",
			assert: func(t *testing.T, err error) {
				if assert.Error(t, err) {
					assert.ErrorIs(t, err, localErrs.BadRequestErr)
				}
			},
			givenSensorRequest: &SensorRequest{
				SensorID:         uuid.NewString(),
				SensorVersion:    "1.0.0",
				Alias:            "lettuce 1",
				Temperature:      0.0,
				Humidity:         0.0,
				PH:               0.0,
				TDS:              0.0,
				EC:               0.0,
				WaterTemperature: 0.0,
				Timestamp:        float64(time.Now().AddDate(0, 0, -30).Unix()),
				Time:             time.Time{},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			err := tt.givenSensorRequest.Bind(nil)
			tt.assert(t, err)
		})
	}
}
