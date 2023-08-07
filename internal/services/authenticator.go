package services

import (
	"context"
	"errors"
	"time"

	localErrs "github.com/WendelHime/hydroponics-metrics-collector/internal/shared/errors"
	"github.com/WendelHime/hydroponics-metrics-collector/internal/shared/models"
	"github.com/auth0/go-auth0/authentication"
	"github.com/auth0/go-auth0/authentication/oauth"
	"github.com/auth0/go-auth0/management"
	"github.com/google/uuid"
)

// Authenticator holds functions for authentication methods
//
//go:generate mockgen -destination authenticator_mock.go -package services github.com/WendelHime/hydroponics-metrics-collector/internal/services Authenticator,OAuth
type Authenticator interface {
	SignIn(ctx context.Context, credentials models.Credentials) (models.Token, error)
}

// OAuth interafces oauth functionalities from auth0
type OAuth interface {
	LoginWithPassword(ctx context.Context, body oauth.LoginWithPasswordRequest, validationOptions oauth.IDTokenValidationOptions, opts ...authentication.RequestOption) (t *oauth.TokenSet, err error)
	LoginWithAuthCodeWithPKCE(ctx context.Context, body oauth.LoginWithAuthCodeWithPKCERequest, validationOptions oauth.IDTokenValidationOptions, opts ...authentication.RequestOption) (t *oauth.TokenSet, err error)
}

type authService struct {
	oauth       OAuth
	environment string
	connection  string
	audience    string
}

// NewAuthService builds an authenticator service
func NewAuthService(oauth OAuth, env, audience string) Authenticator {
	return &authService{
		oauth:       oauth,
		environment: env,
		audience:    audience,
		connection:  "Username-Password-Authentication",
	}
}

func (u *authService) SignIn(ctx context.Context, credentials models.Credentials) (models.Token, error) {
	token, err := u.oauth.LoginWithPassword(ctx, oauth.LoginWithPasswordRequest{
		Username: credentials.Email,
		Password: credentials.Password,
		Scope:    credentials.Scope,
		Audience: u.audience,
		Realm:    u.environment,
	}, oauth.IDTokenValidationOptions{
		MaxAge: 10 * time.Minute,
		Nonce:  uuid.NewString(),
	})
	if err != nil {
		var mngmtErr management.Error
		if errors.As(err, &mngmtErr) {
			if mngmtErr.Status() == 400 {
				return models.Token{}, localErrs.BadRequestErr
			}
			if mngmtErr.Status() == 403 {
				return models.Token{}, localErrs.ForbiddenErr
			}
		}
		return models.Token{}, localErrs.InternalServerErr.WithMsg("failed to login").WithDetails("err", err)
	}

	return models.Token{
		IDToken:      token.IDToken,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		ExpiresIn:    token.ExpiresIn,
	}, nil
}
