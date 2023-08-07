package services

import (
	"context"
	"errors"
	"strings"
	"testing"

	localErrs "github.com/WendelHime/hydroponics-metrics-collector/internal/shared/errors"
	"github.com/WendelHime/hydroponics-metrics-collector/internal/shared/models"
	"github.com/auth0/go-auth0"
	"github.com/auth0/go-auth0/management"
	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"
)

func TestCreateAccount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	successAccount := models.User{
		Name:     "Random User",
		Email:    "random@test.com",
		Password: "UltraSecr3tPassword!",
		Role:     "user",
	}

	ctx := context.Background()
	var tests = []struct {
		name         string
		setup        func() UserService
		assert       func(t *testing.T, err error)
		givenAccount models.User
	}{
		{
			name: "creating account with success",
			setup: func() UserService {
				userManager := NewMockUserManager(ctrl)
				roleManager := NewMockRoleManager(ctrl)
				userManager.EXPECT().Create(gomock.Any(), &management.User{
					Connection:   auth0.String("Username-Password-Authentication"),
					Name:         auth0.String(successAccount.Name),
					Email:        auth0.String(successAccount.Email),
					Password:     auth0.String(successAccount.Password),
					VerifyEmail:  auth0.Bool(true),
					UserMetadata: &map[string]interface{}{"role": successAccount.Role},
				}).Return(nil)
				return NewUserService(userManager, roleManager)
			},
			assert: func(t *testing.T, err error) {
				assert.Nil(t, err)
			},
			givenAccount: successAccount,
		},
		{
			name: "should return conflict if user already exists",
			setup: func() UserService {
				userManager := NewMockUserManager(ctrl)
				roleManager := NewMockRoleManager(ctrl)
				userManager.EXPECT().Create(gomock.Any(), &management.User{
					Connection:   auth0.String("Username-Password-Authentication"),
					Name:         auth0.String(successAccount.Name),
					Email:        auth0.String(successAccount.Email),
					Password:     auth0.String(successAccount.Password),
					VerifyEmail:  auth0.Bool(true),
					UserMetadata: &map[string]interface{}{"role": successAccount.Role},
				}).Return(auth0Error{StatusCode: 409})
				return NewUserService(userManager, roleManager)
			},
			assert: func(t *testing.T, err error) {
				assert.NotNil(t, err)
				assert.ErrorIs(t, err, localErrs.AlreadyExistsErr)
			},
			givenAccount: successAccount,
		},
		{
			name: "should return internal server error on random error",
			setup: func() UserService {
				userManager := NewMockUserManager(ctrl)
				roleManager := NewMockRoleManager(ctrl)
				userManager.EXPECT().Create(gomock.Any(), &management.User{
					Connection:   auth0.String("Username-Password-Authentication"),
					Name:         auth0.String(successAccount.Name),
					Email:        auth0.String(successAccount.Email),
					Password:     auth0.String(successAccount.Password),
					VerifyEmail:  auth0.Bool(true),
					UserMetadata: &map[string]interface{}{"role": successAccount.Role},
				}).Return(errors.New("random error"))
				return NewUserService(userManager, roleManager)
			},
			assert: func(t *testing.T, err error) {
				assert.NotNil(t, err)
				assert.ErrorIs(t, err, localErrs.InternalServerErr)
			},
			givenAccount: successAccount,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			service := tt.setup()
			err := service.CreateAccount(ctx, tt.givenAccount)
			tt.assert(t, err)
		})
	}
}

func TestAssignRoleToUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	roleID := "randomRoleID"
	userID := "randomUserID"

	var tests = []struct {
		name        string
		setup       func() UserService
		assert      func(t *testing.T, err error)
		givenRoleID string
		givenUserID string
	}{
		{
			name: "role assigned to user with success",
			setup: func() UserService {
				userManager := NewMockUserManager(ctrl)
				roleManager := NewMockRoleManager(ctrl)
				roleManager.EXPECT().AssignUsers(gomock.Any(), roleID, []*management.User{{Connection: auth0.String("Username-Password-Authentication"), ID: &userID}}).Return(nil)
				return NewUserService(userManager, roleManager)
			},
			assert: func(t *testing.T, err error) {
				assert.Nil(t, err)
			},
			givenRoleID: roleID,
			givenUserID: userID,
		},
		{
			name: "role assigned to user with random error",
			setup: func() UserService {
				userManager := NewMockUserManager(ctrl)
				roleManager := NewMockRoleManager(ctrl)
				roleManager.EXPECT().AssignUsers(gomock.Any(), roleID, []*management.User{{Connection: auth0.String("Username-Password-Authentication"), ID: &userID}}).Return(errors.New("random error"))
				return NewUserService(userManager, roleManager)
			},
			assert: func(t *testing.T, err error) {
				assert.NotNil(t, err)
				assert.ErrorIs(t, err, localErrs.InternalServerErr)
			},
			givenRoleID: roleID,
			givenUserID: userID,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			userService := tt.setup()
			err := userService.AssignRoleToUser(context.Background(), tt.givenRoleID, tt.givenUserID)
			tt.assert(t, err)
		})
	}
}

func TestGetRolePermissions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	roleID := "randomRoleID"

	var tests = []struct {
		name        string
		setup       func() UserService
		assert      func(t *testing.T, permissions string, err error)
		givenRoleID string
	}{
		{
			name: "getting scopes from user with success",
			setup: func() UserService {
				userManager := NewMockUserManager(ctrl)
				roleManager := NewMockRoleManager(ctrl)
				roleManager.EXPECT().Permissions(gomock.Any(), roleID).Return(
					&management.PermissionList{
						Permissions: []*management.Permission{
							{
								Name: auth0.String("test"),
							},
						},
					},
					nil)
				return NewUserService(userManager, roleManager)
			},
			assert: func(t *testing.T, permissions string, err error) {
				assert.Nil(t, err)
				assert.NotEmpty(t, permissions)
			},
			givenRoleID: roleID,
		},
		{
			name: "role assigned to user with random error",
			setup: func() UserService {
				userManager := NewMockUserManager(ctrl)
				roleManager := NewMockRoleManager(ctrl)
				roleManager.EXPECT().Permissions(gomock.Any(), roleID).Return(nil, errors.New("random error"))
				return NewUserService(userManager, roleManager)
			},
			assert: func(t *testing.T, permissions string, err error) {
				assert.Empty(t, permissions)
				assert.NotNil(t, err)
				assert.ErrorIs(t, err, localErrs.InternalServerErr)
			},
			givenRoleID: roleID,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			userService := tt.setup()
			permissions, err := userService.GetRolePermissions(context.Background(), tt.givenRoleID)
			tt.assert(t, permissions, err)
		})
	}
}

func TestGetUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	successAccount := models.User{
		Name:     "Random User",
		Email:    "Random@test.com",
		Password: "UltraSecr3tPassword!",
		Role:     "user",
	}

	ctx := context.Background()
	var tests = []struct {
		name       string
		setup      func() UserService
		assert     func(t *testing.T, user models.User, err error)
		givenEmail string
	}{
		{
			name: "getting account with success",
			setup: func() UserService {
				userManager := NewMockUserManager(ctrl)
				roleManager := NewMockRoleManager(ctrl)
				userManager.EXPECT().ListByEmail(gomock.Any(), strings.ToLower(successAccount.Email)).Return(
					[]*management.User{
						{
							Name:         auth0.String(successAccount.Name),
							Email:        auth0.String(successAccount.Email),
							Password:     auth0.String(successAccount.Password),
							UserMetadata: &map[string]interface{}{"role": successAccount.Role},
						},
					},
					nil,
				)
				return NewUserService(userManager, roleManager)
			},
			assert: func(t *testing.T, user models.User, err error) {
				assert.Nil(t, err)
				assert.Equal(t, successAccount.Name, user.Name)
				assert.Equal(t, successAccount.Email, user.Email)
				assert.Equal(t, successAccount.Role, user.Role)
				assert.Empty(t, user.Password)
			},
			givenEmail: successAccount.Email,
		},
		{
			name: "user not found",
			setup: func() UserService {
				userManager := NewMockUserManager(ctrl)
				roleManager := NewMockRoleManager(ctrl)
				userManager.EXPECT().ListByEmail(gomock.Any(), strings.ToLower(successAccount.Email)).Return(
					[]*management.User{},
					nil,
				)
				return NewUserService(userManager, roleManager)
			},
			assert: func(t *testing.T, user models.User, err error) {
				assert.Equal(t, models.User{}, user)
				assert.NotNil(t, err)
				assert.ErrorIs(t, err, localErrs.NotFoundErr)
			},
			givenEmail: successAccount.Email,
		},
		{
			name: "should return internal server error on random error",
			setup: func() UserService {
				userManager := NewMockUserManager(ctrl)
				roleManager := NewMockRoleManager(ctrl)
				userManager.EXPECT().ListByEmail(gomock.Any(), strings.ToLower(successAccount.Email)).Return(
					[]*management.User{},
					errors.New("random error"),
				)
				return NewUserService(userManager, roleManager)
			},
			assert: func(t *testing.T, user models.User, err error) {
				assert.Equal(t, models.User{}, user)
				assert.NotNil(t, err)
				assert.ErrorIs(t, err, localErrs.InternalServerErr)
			},
			givenEmail: successAccount.Email,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			service := tt.setup()
			user, err := service.GetUser(ctx, tt.givenEmail)
			tt.assert(t, user, err)
		})
	}
}
