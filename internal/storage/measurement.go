package storage

import (
	"context"
	"time"

	"github.com/WendelHime/hydroponics-metrics-collector/internal/shared/errors"
	"github.com/WendelHime/hydroponics-metrics-collector/internal/shared/models"
)

// MetricRepository implement functions for persisting data
type MetricRepository interface {
	WriteMeasurement(ctx context.Context, request ...models.SensorRequest) error
}

// SensorMeasurement represents the database data structure
type SensorMeasurement struct {
	Table            string    `lp:"measurement"`
	SensorID         string    `lp:"tag,sensor_id"`
	SensorVersion    string    `lp:"tag,sensor_version"`
	Alias            string    `lp:"tag,alias"`
	Temperature      float64   `lp:"field,temperature"`
	Humidity         float64   `lp:"field,humidity"`
	PH               float64   `lp:"field,ph"`
	TDS              float64   `lp:"field,tds"`
	EC               float64   `lp:"field,ec"`
	WaterTemperature float64   `lp:"field,water_temperature"`
	Timestamp        time.Time `lp:"timestamp"`
}

type repository struct {
	database string
	cli      InfluxClient
}

// InfluxClient represents functions from influxdb client used by the storage layer
//
//go:generate mockgen -destination measurement_mock.go -package storage github.com/WendelHime/hydroponics-metrics-collector/internal/storage InfluxClient
type InfluxClient interface {
	WriteData(ctx context.Context, database string, points ...any) error
}

func NewRepository(database string, client InfluxClient) MetricRepository {
	return &repository{database: database, cli: client}
}

func parseRequestToMeasurement(r models.SensorRequest) SensorMeasurement {
	measurement := SensorMeasurement{
		Table:            "metrics",
		SensorID:         r.SensorID,
		SensorVersion:    r.SensorVersion,
		Alias:            r.Alias,
		Temperature:      r.Temperature,
		Humidity:         r.Humidity,
		PH:               r.PH,
		TDS:              r.TDS,
		EC:               r.EC,
		WaterTemperature: r.WaterTemperature,
		Timestamp:        r.Time,
	}
	return measurement
}

func parseRequestsToMeasurements(requests ...models.SensorRequest) []any {
	measurements := make([]any, len(requests))
	for i, r := range requests {
		measurements[i] = parseRequestToMeasurement(r)
	}
	return measurements
}

func (r repository) WriteMeasurement(ctx context.Context, request ...models.SensorRequest) error {
	measurements := parseRequestsToMeasurements(request...)
	err := r.cli.WriteData(ctx, r.database, measurements...)
	if err != nil {
		return errors.InternalServerErr.WithMsg("failed to write data").WithErr(err)
	}

	return nil
}
