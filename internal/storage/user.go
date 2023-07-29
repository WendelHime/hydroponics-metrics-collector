package storage

import (
	"context"
	"errors"
	"time"

	localErrs "github.com/WendelHime/hydroponics-metrics-collector/internal/shared/errors"
	"github.com/WendelHime/hydroponics-metrics-collector/internal/shared/models"

	"github.com/auth0/go-auth0"
	"github.com/auth0/go-auth0/authentication"
	"github.com/auth0/go-auth0/authentication/oauth"
	"github.com/auth0/go-auth0/management"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type UserRepository interface {
	CreateAccount(ctx context.Context, account models.User) error
	GetUser(ctx context.Context, email string) (models.User, error)
	AssignRoleToUser(ctx context.Context, roleID, userID string) error
	GetRolePermissions(ctx context.Context, roleID string) (string, error)
	SignIn(ctx context.Context, credentials models.Credentials) (string, error)
}

type userRepo struct {
	authCli       *authentication.Authentication
	managementCli *management.Management
	environment   string
	connection    string
	audience      string
}

func NewUserRepository(ctx context.Context, domain, clientID, clientSecret, env, audience string) (UserRepository, error) {
	authCli, err := authentication.New(
		ctx,
		domain,
		authentication.WithClientID(clientID),
		authentication.WithClientSecret(clientSecret))
	if err != nil {
		return nil, localErrs.InternalServerErr.WithMsg("failed to create auth0 authentication client").WithDetails("err", err.Error())
	}

	managementCli, err := management.New(domain, management.WithClientCredentials(ctx, clientID, clientSecret))
	if err != nil {
		return nil, localErrs.InternalServerErr.WithMsg("failed to create auth0 management client").WithDetails("err", err.Error())
	}
	return &userRepo{
		authCli:       authCli,
		managementCli: managementCli,
		environment:   env,
		audience:      audience,
		connection:    "Username-Password-Authentication",
	}, nil
}

func (u *userRepo) CreateAccount(ctx context.Context, account models.User) error {
	err := u.managementCli.User.Create(ctx, &management.User{
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
			if status == 401 {
				return localErrs.BadRequestErr
			}
		}
		return localErrs.InternalServerErr.WithMsg("failed to create account").WithErr(err)
	}

	return nil
}

func (u *userRepo) SignIn(ctx context.Context, credentials models.Credentials) (string, error) {
	token, err := u.authCli.OAuth.LoginWithPassword(ctx, oauth.LoginWithPasswordRequest{
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
			if mngmtErr.Status() == 401 {
				return "", localErrs.BadRequestErr
			}
		}
		return "", localErrs.InternalServerErr.WithMsg("failed to login").WithDetails("err", err)
	}

	return token.AccessToken, nil
}

func (u *userRepo) AssignRoleToUser(ctx context.Context, roleID, userID string) error {
	err := u.managementCli.Role.AssignUsers(ctx, roleID, []*management.User{{Connection: &u.connection, ID: &userID}})
	if err != nil {
		log.Warn().Err(err).Msg("failed to assign role to user")
		return localErrs.InternalServerErr.WithErr(err).WithMsg("failed to assign role to user").WithDetails("roleID", roleID).WithDetails("userID", userID)
	}
	return nil
}

func (u *userRepo) GetRolePermissions(ctx context.Context, roleID string) (string, error) {
	permissions, err := u.managementCli.Role.Permissions(ctx, roleID)
	if err != nil {
		log.Warn().Err(err).Msg("failed to get role permissions")
		return "", localErrs.InternalServerErr.WithErr(err).WithMsg("failed to retrieve role permissions").WithDetails("roleID", roleID)
	}
	return permissions.String(), nil
}

func (u *userRepo) GetUser(ctx context.Context, email string) (models.User, error) {
	users, err := u.managementCli.User.ListByEmail(ctx, email)
	if err != nil {
		var mngmtErr management.Error
		if errors.As(err, &mngmtErr) {
			if mngmtErr.Status() == 404 {
				return models.User{}, localErrs.NotFoundErr
			}
		}
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
