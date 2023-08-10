package logic

import (
	"context"
	"errors"
	"testing"

	"github.com/WendelHime/hydroponics-metrics-collector/internal/services"
	localErrs "github.com/WendelHime/hydroponics-metrics-collector/internal/shared/errors"
	"github.com/WendelHime/hydroponics-metrics-collector/internal/shared/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"
)

func TestCreateAccount(t *testing.T) {
	baseAccount := models.User{
		Name:     "Test",
		Email:    "test@test.com",
		Password: "123456",
	}

	baseAccountWithID := models.User{
		ID:       uuid.NewString(),
		Name:     "Test",
		Email:    "test@test.com",
		Password: "123456",
	}

	var tests = []struct {
		name         string
		setup        func(ctrl *gomock.Controller) UserLogic
		givenAccount models.User
		assert       func(t *testing.T, err error)
	}{
		{
			name: "create account with success",
			setup: func(ctrl *gomock.Controller) UserLogic {
				roleID := uuid.NewString()
				userService := services.NewMockUserService(ctrl)
				userService.EXPECT().CreateAccount(gomock.Any(), baseAccount).Return(nil).Times(1)
				userService.EXPECT().GetUser(gomock.Any(), baseAccount.Email).Return(baseAccountWithID, nil).Times(1)
				userService.EXPECT().AssignRoleToUser(gomock.Any(), roleID, baseAccountWithID.ID).Return(nil).Times(1)
				authService := services.NewMockAuthenticator(ctrl)
				return NewUserLogic(userService, authService, roleID)
			},
			givenAccount: baseAccount,
			assert: func(t *testing.T, err error) {
				assert.Nil(t, err)
			},
		},
		{
			name: "should return an error if CreateAccount call fails",
			setup: func(ctrl *gomock.Controller) UserLogic {
				roleID := uuid.NewString()
				userService := services.NewMockUserService(ctrl)
				userService.EXPECT().CreateAccount(gomock.Any(), baseAccount).Return(errors.New("random error")).Times(1)
				authService := services.NewMockAuthenticator(ctrl)
				return NewUserLogic(userService, authService, roleID)
			},
			givenAccount: baseAccount,
			assert: func(t *testing.T, err error) {
				assert.NotNil(t, err)
			},
		},
		{
			name: "should return an error if failed to retrieve the user after creation",
			setup: func(ctrl *gomock.Controller) UserLogic {
				roleID := uuid.NewString()
				userService := services.NewMockUserService(ctrl)
				userService.EXPECT().CreateAccount(gomock.Any(), baseAccount).Return(nil).Times(1)
				userService.EXPECT().GetUser(gomock.Any(), baseAccount.Email).Return(baseAccountWithID, errors.New("random error")).Times(1)
				authService := services.NewMockAuthenticator(ctrl)
				return NewUserLogic(userService, authService, roleID)
			},
			givenAccount: baseAccount,
			assert: func(t *testing.T, err error) {
				assert.NotNil(t, err)
			},
		},
		{
			name: "should return an error if failed to assign role to user",
			setup: func(ctrl *gomock.Controller) UserLogic {
				roleID := uuid.NewString()
				userService := services.NewMockUserService(ctrl)
				userService.EXPECT().CreateAccount(gomock.Any(), baseAccount).Return(nil).Times(1)
				userService.EXPECT().GetUser(gomock.Any(), baseAccount.Email).Return(baseAccountWithID, nil).Times(1)
				userService.EXPECT().AssignRoleToUser(gomock.Any(), roleID, baseAccountWithID.ID).Return(errors.New("random error")).Times(1)
				authService := services.NewMockAuthenticator(ctrl)
				return NewUserLogic(userService, authService, roleID)
			},
			givenAccount: baseAccount,
			assert: func(t *testing.T, err error) {
				assert.NotNil(t, err)
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			logic := tt.setup(ctrl)
			err := logic.CreateAccount(context.Background(), tt.givenAccount)
			tt.assert(t, err)
		})
	}
}

func TestLogin(t *testing.T) {
	basicCredentials := models.Credentials{
		Email:    "test@test.com",
		Password: "123456",
	}
	credentialsWithScope := basicCredentials
	scope := "w:random_scope"
	credentialsWithScope.Scope = scope
	baseAccount := models.User{
		ID:            uuid.NewString(),
		Name:          "Test",
		Email:         "test@test.com",
		Password:      "123456",
		Role:          "randomRole",
		EmailVerified: true,
	}
	accountWithoutEmailVerified := baseAccount
	accountWithoutEmailVerified.EmailVerified = false
	baseToken := models.Token{
		IDToken:      uuid.NewString(),
		AccessToken:  uuid.NewString(),
		RefreshToken: uuid.NewString(),
		ExpiresIn:    10,
	}
	var tests = []struct {
		name             string
		setup            func(ctrl *gomock.Controller) UserLogic
		givenCredentials models.Credentials
		assert           func(t *testing.T, token models.Token, err error)
	}{
		{
			name: "login with success",
			setup: func(ctrl *gomock.Controller) UserLogic {
				roleID := uuid.NewString()
				userService := services.NewMockUserService(ctrl)
				userService.EXPECT().GetUser(gomock.Any(), basicCredentials.Email).Return(baseAccount, nil).Times(1)
				userService.EXPECT().GetRolePermissions(gomock.Any(), baseAccount.Role).Return(scope, nil).Times(1)
				authService := services.NewMockAuthenticator(ctrl)
				authService.EXPECT().SignIn(gomock.Any(), credentialsWithScope).Return(baseToken, nil).Times(1)
				return NewUserLogic(userService, authService, roleID)
			},
			givenCredentials: basicCredentials,
			assert: func(t *testing.T, token models.Token, err error) {
				assert.Nil(t, err)
				assert.Equal(t, baseToken, token)
			},
		},
		{
			name: "login failed because the email wasn't verified",
			setup: func(ctrl *gomock.Controller) UserLogic {
				roleID := uuid.NewString()
				userService := services.NewMockUserService(ctrl)
				userService.EXPECT().GetUser(gomock.Any(), basicCredentials.Email).Return(accountWithoutEmailVerified, nil).Times(1)
				authService := services.NewMockAuthenticator(ctrl)
				return NewUserLogic(userService, authService, roleID)
			},
			givenCredentials: basicCredentials,
			assert: func(t *testing.T, token models.Token, err error) {
				assert.NotNil(t, err)
				assert.ErrorIs(t, err, localErrs.ForbiddenErr)
				assert.Equal(t, models.Token{}, token)
			},
		},
		{
			name: "login failed because GetUser returned an error",
			setup: func(ctrl *gomock.Controller) UserLogic {
				roleID := uuid.NewString()
				userService := services.NewMockUserService(ctrl)
				userService.EXPECT().GetUser(gomock.Any(), basicCredentials.Email).Return(baseAccount, errors.New("random error")).Times(1)
				authService := services.NewMockAuthenticator(ctrl)
				return NewUserLogic(userService, authService, roleID)
			},
			givenCredentials: basicCredentials,
			assert: func(t *testing.T, token models.Token, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, models.Token{}, token)
			},
		},
		{
			name: "login failed because the we couldn't retrieve role permissions",
			setup: func(ctrl *gomock.Controller) UserLogic {
				roleID := uuid.NewString()
				userService := services.NewMockUserService(ctrl)
				userService.EXPECT().GetUser(gomock.Any(), basicCredentials.Email).Return(baseAccount, nil).Times(1)
				userService.EXPECT().GetRolePermissions(gomock.Any(), baseAccount.Role).Return("", errors.New("random error")).Times(1)
				authService := services.NewMockAuthenticator(ctrl)
				return NewUserLogic(userService, authService, roleID)
			},
			givenCredentials: basicCredentials,
			assert: func(t *testing.T, token models.Token, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, models.Token{}, token)
			},
		},
		{
			name: "login failed because signin failed",
			setup: func(ctrl *gomock.Controller) UserLogic {
				roleID := uuid.NewString()
				userService := services.NewMockUserService(ctrl)
				userService.EXPECT().GetUser(gomock.Any(), basicCredentials.Email).Return(baseAccount, nil).Times(1)
				userService.EXPECT().GetRolePermissions(gomock.Any(), baseAccount.Role).Return(scope, nil).Times(1)
				authService := services.NewMockAuthenticator(ctrl)
				authService.EXPECT().SignIn(gomock.Any(), credentialsWithScope).Return(models.Token{}, errors.New("random error")).Times(1)
				return NewUserLogic(userService, authService, roleID)
			},
			givenCredentials: basicCredentials,
			assert: func(t *testing.T, token models.Token, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, models.Token{}, token)
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			logic := tt.setup(ctrl)
			token, err := logic.Login(context.Background(), tt.givenCredentials)
			tt.assert(t, token, err)
		})
	}
}
