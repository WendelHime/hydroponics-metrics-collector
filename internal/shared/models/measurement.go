package models

import (
	"net/http"
	"time"

	localErrs "github.com/WendelHime/hydroponics-metrics-collector/internal/shared/errors"
	"github.com/go-playground/validator/v10"
)

// SensorRequest is used to represent metrics registered by any sensors connected to the raspberry
type SensorRequest struct {
	SensorID         string    `json:"sensor_id" validate:"required"`
	SensorVersion    string    `json:"sensor_version" validate:"required"`
	Alias            string    `json:"alias" validate:"required"`
	Temperature      float64   `json:"temperature"`
	Humidity         float64   `json:"humidity"`
	PH               float64   `json:"ph"`
	TDS              float64   `json:"tds"`
	EC               float64   `json:"ec"`
	WaterTemperature float64   `json:"water_temperature"`
	Timestamp        float64   `json:"timestamp" validate:"required"`
	Time             time.Time `json:"-"`
}

func (s *SensorRequest) Bind(r *http.Request) error {
	validate := validator.New()
	err := validate.Struct(s)
	if err != nil {
		return localErrs.BadRequestErr.WithErr(err).WithMsg("failed to validate request")
	}
	return nil
}
