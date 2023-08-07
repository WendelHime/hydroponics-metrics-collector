// Package services contains implementations that refer to external services
package services

import (
	"context"
	"errors"
	"strings"

	localErrs "github.com/WendelHime/hydroponics-metrics-collector/internal/shared/errors"
	"github.com/WendelHime/hydroponics-metrics-collector/internal/shared/models"

	"github.com/auth0/go-auth0"
	"github.com/auth0/go-auth0/management"
	"github.com/rs/zerolog/log"
)

// UserService holds functions for user management
//
//go:generate mockgen -destination user_management_mock.go -package services github.com/WendelHime/hydroponics-metrics-collector/internal/services UserService,UserManager,RoleManager
type UserService interface {
	CreateAccount(ctx context.Context, account models.User) error
	GetUser(ctx context.Context, email string) (models.User, error)
	AssignRoleToUser(ctx context.Context, roleID, userID string) error
	GetRolePermissions(ctx context.Context, roleID string) (string, error)
}

// UserManager interface user management functionalities from oauth service
type UserManager interface {
	Create(ctx context.Context, u *management.User, opts ...management.RequestOption) error
	ListByEmail(ctx context.Context, email string, opts ...management.RequestOption) (us []*management.User, err error)
}

// RoleManager interface role management functionalities from oauth service
type RoleManager interface {
	AssignUsers(ctx context.Context, id string, users []*management.User, opts ...management.RequestOption) error
	Permissions(ctx context.Context, id string, opts ...management.RequestOption) (p *management.PermissionList, err error)
}

type userService struct {
	userManager UserManager
	roleManager RoleManager
	connection  string
}

// NewUserService builds a new user service that allows the application to manage users
func NewUserService(userManager UserManager, roleManager RoleManager) UserService {
	return &userService{
		userManager: userManager,
		roleManager: roleManager,
		connection:  "Username-Password-Authentication",
	}
}

func (u *userService) CreateAccount(ctx context.Context, account models.User) error {
	err := u.userManager.Create(ctx, &management.User{
		Connection:   auth0.String(u.connection),
		Name:         auth0.String(account.Name),
		Email:        auth0.String(account.Email),
		Password:     auth0.String(account.Password),
		VerifyEmail:  auth0.Bool(true),
		UserMetadata: &map[string]interface{}{"role": account.Role},
	})
	if err != nil {
		var mngmtErr management.Error
		if errors.As(err, &mngmtErr) {
			status := mngmtErr.Status()
			if status == 409 {
				return localErrs.AlreadyExistsErr
			}
		}
		return localErrs.InternalServerErr.WithMsg("failed to create account").WithErr(err)
	}

	return nil
}

func (u *userService) AssignRoleToUser(ctx context.Context, roleID, userID string) error {
	err := u.roleManager.AssignUsers(ctx, roleID, []*management.User{{Connection: auth0.String(u.connection), ID: &userID}})
	if err != nil {
		log.Warn().Err(err).Msg("failed to assign role to user")
		return localErrs.InternalServerErr.WithErr(err).WithMsg("failed to assign role to user").WithDetails("roleID", roleID).WithDetails("userID", userID)
	}
	return nil
}

func (u *userService) GetRolePermissions(ctx context.Context, roleID string) (string, error) {
	permissions, err := u.roleManager.Permissions(ctx, roleID)
	if err != nil {
		log.Warn().Err(err).Msg("failed to get role permissions")
		return "", localErrs.InternalServerErr.WithErr(err).WithMsg("failed to retrieve role permissions").WithDetails("roleID", roleID)
	}
	return permissions.String(), nil
}

func (u *userService) GetUser(ctx context.Context, email string) (models.User, error) {
	users, err := u.userManager.ListByEmail(ctx, strings.ToLower(email))
	if err != nil {
		log.Warn().Err(err).Msg("failed when retrieving user")
		return models.User{}, localErrs.InternalServerErr.WithErr(err).WithMsg("failed retrieving users with provided email")
	}
	if len(users) > 0 {
		return models.User{
			ID:            users[0].GetID(),
			Name:          users[0].GetName(),
			Email:         users[0].GetEmail(),
			Role:          users[0].GetUserMetadata()["role"].(string),
			EmailVerified: users[0].GetEmailVerified(),
		}, nil
	}
	return models.User{}, localErrs.NotFoundErr.WithMsg("user not found").WithDetails("email", email)
}
