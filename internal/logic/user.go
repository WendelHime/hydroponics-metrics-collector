package logic

import (
	"context"

	"github.com/WendelHime/hydroponics-metrics-collector/internal/services"
	localErrs "github.com/WendelHime/hydroponics-metrics-collector/internal/shared/errors"
	"github.com/WendelHime/hydroponics-metrics-collector/internal/shared/models"
)

// UserLogic contain user correlated logic
//
//go:generate mockgen -destination user_mock.go -package logic github.com/WendelHime/hydroponics-metrics-collector/internal/logic UserLogic
type UserLogic interface {
	CreateAccount(ctx context.Context, account models.User) error
	Login(ctx context.Context, credentials models.Credentials) (models.Token, error)
}

func NewUserLogic(userService services.UserService, authService services.Authenticator, roleID string) UserLogic {
	return &userLogic{
		userService: userService,
		authService: authService,
		roleID:      roleID,
	}
}

type userLogic struct {
	userService services.UserService
	authService services.Authenticator
	roleID      string
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
