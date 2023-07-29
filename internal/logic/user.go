package logic

import (
	"context"

	localErrs "github.com/WendelHime/hydroponics-metrics-collector/internal/shared/errors"
	"github.com/WendelHime/hydroponics-metrics-collector/internal/shared/models"
	"github.com/WendelHime/hydroponics-metrics-collector/internal/storage"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
)

type UserLogic interface {
	CreateAccount(ctx context.Context, account models.User) error
	Login(ctx context.Context, credentials models.Credentials) (string, error)
}

func NewUserLogic(repository storage.UserRepository, roleID string) UserLogic {
	// roleID := "rol_wVYCnzv9toZu1ISL"
	return &userLogic{
		repository: repository,
		roleID:     roleID,
	}
}

type userLogic struct {
	repository storage.UserRepository
	roleID     string
}

func (l *userLogic) CreateAccount(ctx context.Context, account models.User) error {
	validate := validator.New()
	err := validate.Struct(account)
	if err != nil {
		log.Warn().Err(err).Msg("failed to validate request")
		return err
	}

	err = l.repository.CreateAccount(ctx, account)
	if err != nil {
		return err
	}

	user, err := l.repository.GetUser(ctx, account.Email)
	if err != nil {
		return err
	}

	err = l.repository.AssignRoleToUser(ctx, l.roleID, user.ID)
	if err != nil {
		return err
	}

	return nil
}

func (l *userLogic) Login(ctx context.Context, credentials models.Credentials) (string, error) {
	validate := validator.New()
	err := validate.Struct(credentials)
	if err != nil {
		log.Warn().Err(err).Msg("failed to validate request")
		return "", err
	}

	user, err := l.repository.GetUser(ctx, credentials.Email)
	if err != nil {
		return "", err
	}

	if !user.EmailVerified {
		return "", localErrs.ForbiddenErr
	}

	permissions, err := l.repository.GetRolePermissions(ctx, user.Role)
	if err != nil {
		return "", err
	}

	credentials.Scope = permissions

	token, err := l.repository.SignIn(ctx, credentials)
	if err != nil {
		return "", err
	}
	return token, nil
}
