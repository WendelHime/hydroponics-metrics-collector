package logic

import (
	"context"
	"errors"

	"github.com/WendelHime/hydroponics-metrics-collector/internal/services"
	localErrs "github.com/WendelHime/hydroponics-metrics-collector/internal/shared/errors"
	"github.com/WendelHime/hydroponics-metrics-collector/internal/shared/models"
	"github.com/WendelHime/hydroponics-metrics-collector/internal/storage"
)

// UserLogic contain user correlated logic
//
//go:generate mockgen -destination user_mock.go -package logic github.com/WendelHime/hydroponics-metrics-collector/internal/logic UserLogic
type UserLogic interface {
	CreateAccount(ctx context.Context, account models.User) error
	Login(ctx context.Context, credentials models.Credentials) (models.Token, error)
	AddDevice(ctx context.Context, userID string, newDevice string) error
	GetDevices(ctx context.Context, userID string) ([]string, error)
}

func NewUserLogic(userService services.UserService, authService services.Authenticator, deviceRepo storage.UserDeviceRepository, roleID string) UserLogic {
	return &userLogic{
		userService:          userService,
		authService:          authService,
		userDeviceRepository: deviceRepo,
		roleID:               roleID,
	}
}

type userLogic struct {
	userService          services.UserService
	authService          services.Authenticator
	userDeviceRepository storage.UserDeviceRepository
	roleID               string
}

func (l *userLogic) CreateAccount(ctx context.Context, account models.User) error {
	err := l.userService.CreateAccount(ctx, account)
	if err != nil {
		return err
	}

	user, err := l.userService.GetUser(ctx, account.Email)
	if err != nil {
		return err
	}

	err = l.userService.AssignRoleToUser(ctx, l.roleID, user.ID)
	if err != nil {
		return err
	}

	return nil
}

func (l *userLogic) Login(ctx context.Context, credentials models.Credentials) (models.Token, error) {
	user, err := l.userService.GetUser(ctx, credentials.Email)
	if err != nil {
		return models.Token{}, err
	}

	if !user.EmailVerified {
		return models.Token{}, localErrs.ForbiddenErr
	}

	permissions, err := l.userService.GetRolePermissions(ctx, user.Role)
	if err != nil {
		return models.Token{}, err
	}

	credentials.Scope = permissions

	token, err := l.authService.SignIn(ctx, credentials)
	if err != nil {
		return models.Token{}, err
	}
	return token, nil
}

func (l *userLogic) AddDevice(ctx context.Context, userID string, newDevice string) error {
	var localErr *localErrs.Error
	currentDevices, err := l.userDeviceRepository.GetDevicesFromUser(ctx, userID)
	if err != nil {
		if errors.As(err, &localErr) {
			if localErr.StatusCode == 404 {
				err = l.userDeviceRepository.CreateUserDevice(ctx, userID, newDevice)
				if err != nil {
					return err
				}
			}
		}
		return err
	}

	err = l.userDeviceRepository.AddDeviceToUser(ctx, userID, newDevice, currentDevices)
	return err
}

func (l *userLogic) GetDevices(ctx context.Context, userID string) ([]string, error) {
	currentDevices, err := l.userDeviceRepository.GetDevicesFromUser(ctx, userID)
	if err != nil {
		return []string{}, err
	}

	return currentDevices, nil
}
