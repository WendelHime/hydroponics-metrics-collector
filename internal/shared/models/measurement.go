package models

import "time"

// SensorRequest is used to represent metrics registered by any sensors connected to the raspberry
type SensorRequest struct {
	SensorID      string    `json:"sensor_id" validate:"required"`
	SensorVersion string    `json:"sensor_version" validate:"required"`
	Alias         string    `json:"alias" validate:"required"`
	Temperature   float64   `json:"temperature"`
	Humidity      float64   `json:"humidity"`
	PH            float64   `json:"ph"`
	TDS           float64   `json:"tds"`
	EC            float64   `json:"ec"`
	Timestamp     float64   `json:"timestamp" validate:"required"`
	Time          time.Time `json:"-"`
}
