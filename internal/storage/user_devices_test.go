package storage

import (
	"context"
	"testing"

	localErrs "github.com/WendelHime/hydroponics-metrics-collector/internal/shared/errors"

	"github.com/stretchr/testify/assert"
)

func TestGetDevicesFromUser(t *testing.T) {
	ctx := context.Background()
	cli := newFirestoreTestClient(ctx)
	defer cli.Close()

	var tests = []struct {
		name        string
		assert      func(t *testing.T, devices []string, err error)
		setup       func() UserDeviceRepository
		givenUserID string
	}{
		{
			name: "should return not found when user doesn't have a sensor",
			assert: func(t *testing.T, devices []string, err error) {
				assert.Empty(t, devices)
				if assert.Error(t, err) {
					assert.ErrorIs(t, err, localErrs.NotFoundErr)
				}
			},
			setup: func() UserDeviceRepository {
				return NewUserDeviceRepository(cli)
			},
			givenUserID: "randomID",
		},
		{
			name: "retrieve user devices with success",
			assert: func(t *testing.T, devices []string, err error) {
				assert.Nil(t, err)
				assert.Contains(t, devices, "sensorID")
			},
			setup: func() UserDeviceRepository {
				cli.Collection("user_devices").Doc("userID").Set(ctx, UserDevices{UserID: "userID", Devices: []string{"sensorID"}})
				return NewUserDeviceRepository(cli)
			},
			givenUserID: "userID",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			repository := tt.setup()
			devices, err := repository.GetDevicesFromUser(ctx, tt.givenUserID)
			tt.assert(t, devices, err)
		})
	}
}
