package storage

import (
	"context"

	localErrs "github.com/WendelHime/hydroponics-metrics-collector/internal/shared/errors"

	"cloud.google.com/go/firestore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserDevices struct {
	UserID  string   `firestore:"user_id"`
	Devices []string `firestore:"devices,omitempty"`
}

type UserDeviceRepository interface {
	GetDevicesFromUser(ctx context.Context, userID string) ([]string, error)
	AddDeviceToUser(ctx context.Context, userID, newDevice string, currentDevices []string) error
	CreateUserDevice(ctx context.Context, userID, newDevice string) error
}

type userDeviceRepository struct {
	client *firestore.Client
}

func NewUserDeviceRepository(client *firestore.Client) UserDeviceRepository {
	return &userDeviceRepository{client: client}
}

func (u *userDeviceRepository) GetDevicesFromUser(ctx context.Context, userID string) ([]string, error) {
	doc, err := u.client.Collection("user_devices").Doc(userID).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, localErrs.NotFoundErr.WithMsg("user without correlated devices").WithErr(err)
		}
		return nil, localErrs.InternalServerErr.WithMsg("failed to retrieve user devices").WithErr(err)
	}

	var userDevices UserDevices
	err = doc.DataTo(&userDevices)
	if err != nil {
		return nil, localErrs.InternalServerErr.WithMsg("failed to parse user devices struct").WithErr(err)
	}

	return userDevices.Devices, nil
}

func (u *userDeviceRepository) AddDeviceToUser(ctx context.Context, userID string, newDevice string, currentDevices []string) error {
	_, err := u.client.Collection("user_devices").Doc(userID).Update(ctx, []firestore.Update{
		{Path: "devices", Value: append(currentDevices, newDevice)},
	})

	if err != nil {
		return err
	}

	return nil
}

func (u *userDeviceRepository) CreateUserDevice(ctx context.Context, userID string, newDevice string) error {
	userDevice := UserDevices{UserID: userID, Devices: []string{newDevice}}
	_, err := u.client.Collection("user_devices").Doc(userID).Set(ctx, userDevice)
	if err != nil {
		return err
	}

	return nil
}
