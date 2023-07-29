package storage

import (
	"context"
	"time"

	"github.com/InfluxCommunity/influxdb3-go/influx"
	"github.com/rs/zerolog/log"

	"github.com/WendelHime/hydroponics-metrics-collector/internal/shared/models"
)

// MetricRepository implement functions for persisting data
type MetricRepository interface {
	WriteMeasurement(ctx context.Context, request models.SensorRequest) error
}

type SensorMeasurement struct {
	Table         string    `lp:"measurement"`
	SensorID      string    `lp:"tag,sensor_id"`
	SensorVersion string    `lp:"tag,sensor_version"`
	Alias         string    `lp:"tag,alias"`
	Temperature   float64   `lp:"field,temperature"`
	Humidity      float64   `lp:"field,humidity"`
	PH            float64   `lp:"field,ph"`
	TDS           float64   `lp:"field,tds"`
	EC            float64   `lp:"field,ec"`
	Timestamp     time.Time `lp:"timestamp"`
}

type repository struct {
	database string
	config   influx.Configs
}

func NewRepository(database string, config influx.Configs) MetricRepository {
	return &repository{database: database, config: config}
}

func parseRequestToMeasurement(r models.SensorRequest) SensorMeasurement {
	measurement := SensorMeasurement{
		Table:         "metrics",
		SensorID:      r.SensorID,
		SensorVersion: r.SensorVersion,
		Alias:         r.Alias,
		Temperature:   r.Temperature,
		Humidity:      r.Humidity,
		PH:            r.PH,
		TDS:           r.TDS,
		EC:            r.EC,
		Timestamp:     r.Time,
	}
	return measurement
}

func parseRequestsToMeasurements(requests ...models.SensorRequest) []SensorMeasurement {
	measurements := make([]SensorMeasurement, len(requests))
	for i, r := range requests {
		measurements[i] = parseRequestToMeasurement(r)
	}
	return measurements
}

func (r repository) WriteMeasurement(ctx context.Context, request models.SensorRequest) error {
	cli, err := influx.New(r.config)
	if err != nil {
		log.Error().Err(err).Msg("failed to create influx client")
		return err
	}
	defer cli.Close()

	measurements := parseRequestToMeasurement(request)
	err = cli.WriteData(ctx, r.database, measurements)
	if err != nil {
		log.Error().Err(err).Msg("failed to write data")
		return err
	}

	return nil
}
