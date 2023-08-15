package storage

import (
	"context"
	"testing"

	localErrs "github.com/WendelHime/hydroponics-metrics-collector/internal/shared/errors"
	"github.com/google/uuid"

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

func TestAddUserDevice(t *testing.T) {
	ctx := context.Background()
	cli := newFirestoreTestClient(ctx)
	defer cli.Close()

	var tests = []struct {
		name                string
		assert              func(t *testing.T, givenUserID string, err error)
		setup               func(givenUserID string, givenCurrentDevices []string) UserDeviceRepository
		givenUserID         string
		givenNewDevice      string
		givenCurrentDevices []string
	}{
		{
			name: "Add user device with success",
			assert: func(t *testing.T, givenUserID string, err error) {
				assert.Nil(t, err)
				doc, err := cli.Collection("user_devices").Doc(givenUserID).Get(ctx)
				assert.Nil(t, err)
				var userDevices UserDevices
				err = doc.DataTo(&userDevices)
				assert.Nil(t, err)

				assert.Len(t, userDevices.Devices, 2)
				assert.Contains(t, userDevices.Devices, "new_device")
			},
			setup: func(givenUserID string, givenCurrentDevices []string) UserDeviceRepository {
				cli.Collection("user_devices").Doc(givenUserID).Set(ctx, UserDevices{
					UserID:  givenUserID,
					Devices: givenCurrentDevices,
				})
				return NewUserDeviceRepository(cli)
			},
			givenUserID:         "userID1",
			givenNewDevice:      "new_device",
			givenCurrentDevices: []string{"old_device"},
		},
		{
			name: "Add user device to unexistent account should succeed",
			assert: func(t *testing.T, _ string, err error) {
				if assert.Error(t, err) {
					assert.ErrorIs(t, err, localErrs.InternalServerErr)
				}
			},
			setup: func(_ string, _ []string) UserDeviceRepository {
				return NewUserDeviceRepository(cli)
			},
			givenUserID:         "userID2",
			givenNewDevice:      "new_device",
			givenCurrentDevices: make([]string, 0),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			repo := tt.setup(tt.givenUserID, tt.givenCurrentDevices)
			err := repo.AddDeviceToUser(ctx, tt.givenUserID, tt.givenNewDevice, tt.givenCurrentDevices)
			tt.assert(t, tt.givenUserID, err)
		})
	}
}

func TestCreateUserDevice(t *testing.T) {
	ctx := context.Background()
	cli := newFirestoreTestClient(ctx)
	defer cli.Close()

	var tests = []struct {
		name           string
		assert         func(t *testing.T, givenUserID string, givenNewDevice string, err error)
		setup          func() UserDeviceRepository
		givenUserID    string
		givenNewDevice string
	}{
		{
			name: "create user device with success",
			assert: func(t *testing.T, givenUserID string, givenNewDevice string, err error) {
				assert.Nil(t, err)

				doc, err := cli.Collection("user_devices").Doc(givenUserID).Get(ctx)
				assert.Nil(t, err)
				var userDevices UserDevices
				err = doc.DataTo(&userDevices)
				assert.Nil(t, err)

				assert.Len(t, userDevices.Devices, 1)
				assert.Contains(t, userDevices.Devices, givenNewDevice)
			},
			setup: func() UserDeviceRepository {
				return NewUserDeviceRepository(cli)
			},
			givenUserID:    uuid.NewString(),
			givenNewDevice: uuid.NewString(),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			repo := tt.setup()
			err := repo.CreateUserDevice(ctx, tt.givenUserID, tt.givenNewDevice)
			tt.assert(t, tt.givenUserID, tt.givenNewDevice, err)
		})
	}
}
