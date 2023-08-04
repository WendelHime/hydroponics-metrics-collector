package services

import (
	"context"
	"errors"
	"testing"

	localErrs "github.com/WendelHime/hydroponics-metrics-collector/internal/shared/errors"
	"github.com/WendelHime/hydroponics-metrics-collector/internal/shared/models"
	"github.com/auth0/go-auth0/authentication/oauth"
	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"
)

type oauthError struct {
	StatusCode int    `json:"statusCode"`
	Err        string `json:"error"`
	Message    string `json:"message"`
}

func (o oauthError) Status() int {
	return o.StatusCode
}

func (oauthError) Error() string {
	return ""
}

func TestSignIn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	basicCredentials := models.Credentials{
		Email:    "test@test.com",
		Password: "^SuperStr0ngRandomPassword$",
	}
	emptyCredentials := models.Credentials{}

	badRequestErr := oauthError{
		StatusCode: 400,
		Err:        "bad request",
		Message:    "invalid credentials provided",
	}

	forbiddenRequestErr := oauthError{
		StatusCode: 403,
		Err:        "forbidden",
		Message:    "forbidden",
	}

	var tests = []struct {
		name             string
		assert           func(t *testing.T, accessToken string, err error)
		authService      func() Authenticator
		givenCredentials models.Credentials
	}{
		{
			name: "creating user with success",
			assert: func(t *testing.T, accessToken string, err error) {
				assert.Nil(t, err)
				assert.NotEmpty(t, accessToken)
			},
			authService: func() Authenticator {
				oAuth := NewMockOAuth(ctrl)
				audience := "test"
				environment := "local"

				oAuth.EXPECT().LoginWithPassword(
					gomock.Any(),
					oauth.LoginWithPasswordRequest{
						Username: basicCredentials.Email,
						Password: basicCredentials.Password,
						Scope:    basicCredentials.Scope,
						Audience: audience,
						Realm:    environment,
					},
					gomock.Any(),
				).Return(&oauth.TokenSet{AccessToken: "random token"}, nil)
				return NewAuthService(oAuth, environment, audience)
			},
			givenCredentials: basicCredentials,
		},
		{
			name: "should return bad request when email and password are empty or invalid",
			assert: func(t *testing.T, accessToken string, err error) {
				assert.Empty(t, accessToken)
				assert.NotNil(t, err)
				assert.ErrorIs(t, err, localErrs.BadRequestErr)
			},
			authService: func() Authenticator {
				oAuth := NewMockOAuth(ctrl)
				audience := "test"
				environment := "local"

				oAuth.EXPECT().LoginWithPassword(
					gomock.Any(),
					oauth.LoginWithPasswordRequest{
						Username: emptyCredentials.Email,
						Password: emptyCredentials.Password,
						Scope:    emptyCredentials.Scope,
						Audience: audience,
						Realm:    environment,
					},
					gomock.Any(),
				).Return(nil, badRequestErr)
				return NewAuthService(oAuth, environment, audience)
			},
			givenCredentials: emptyCredentials,
		},
		{
			name: "should return forbidden request when email or password are wrong",
			assert: func(t *testing.T, accessToken string, err error) {
				assert.Empty(t, accessToken)
				assert.NotNil(t, err)
				assert.ErrorIs(t, err, localErrs.ForbiddenErr)
			},
			authService: func() Authenticator {
				oAuth := NewMockOAuth(ctrl)
				audience := "test"
				environment := "local"

				oAuth.EXPECT().LoginWithPassword(
					gomock.Any(),
					oauth.LoginWithPasswordRequest{
						Username: emptyCredentials.Email,
						Password: emptyCredentials.Password,
						Scope:    emptyCredentials.Scope,
						Audience: audience,
						Realm:    environment,
					},
					gomock.Any(),
				).Return(nil, forbiddenRequestErr)
				return NewAuthService(oAuth, environment, audience)
			},
			givenCredentials: emptyCredentials,
		},
		{
			name: "should return internal server error when any random error happens",
			assert: func(t *testing.T, accessToken string, err error) {
				assert.Empty(t, accessToken)
				assert.NotNil(t, err)
				assert.ErrorIs(t, err, localErrs.InternalServerErr)
			},
			authService: func() Authenticator {
				oAuth := NewMockOAuth(ctrl)
				audience := "test"
				environment := "local"

				oAuth.EXPECT().LoginWithPassword(
					gomock.Any(),
					oauth.LoginWithPasswordRequest{
						Username: emptyCredentials.Email,
						Password: emptyCredentials.Password,
						Scope:    emptyCredentials.Scope,
						Audience: audience,
						Realm:    environment,
					},
					gomock.Any(),
				).Return(nil, errors.New("random error"))
				return NewAuthService(oAuth, environment, audience)
			},
			givenCredentials: emptyCredentials,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			accessToken, err := tt.authService().SignIn(context.Background(), tt.givenCredentials)
			tt.assert(t, accessToken, err)
		})
	}
}
