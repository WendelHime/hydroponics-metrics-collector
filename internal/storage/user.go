package storage

import (
	"context"
	"errors"
	"time"

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
	Login(ctx context.Context, credentials models.Credentials) (string, error)
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
		log.Error().Err(err).Msg("failed to create auth0 authentication client")
		return nil, err
	}

	managementCli, err := management.New(domain, management.WithClientCredentials(ctx, clientID, clientSecret))
	if err != nil {
		log.Error().Err(err).Msg("failed to create auth0 management client")
		return nil, err
	}
	return &userRepo{
		authCli:       authCli,
		managementCli: managementCli,
		environment:   env,
		// "https://hydroponics.capivarasousadas.com/"
		audience:   audience,
		connection: "Username-Password-Authentication",
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
		log.Warn().Err(err).Msg("failed to create account")
		return err
	}

	return nil
}

func (u *userRepo) Login(ctx context.Context, credentials models.Credentials) (string, error) {
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
		log.Warn().Err(err).Msg("failed to login")
		return "", err
	}

	return token.AccessToken, nil
}

func (u *userRepo) AssignRoleToUser(ctx context.Context, roleID, userID string) error {
	err := u.managementCli.Role.AssignUsers(ctx, roleID, []*management.User{{Connection: &u.connection, ID: &userID}})
	if err != nil {
		log.Warn().Err(err).Msg("failed to assign role to user")
		return err
	}
	return nil
}

func (u *userRepo) GetRolePermissions(ctx context.Context, roleID string) (string, error) {
	permissions, err := u.managementCli.Role.Permissions(ctx, roleID)
	if err != nil {
		log.Warn().Err(err).Msg("failed to get role permissions")
		return "", err
	}
	return permissions.String(), nil
}

func (u *userRepo) GetUser(ctx context.Context, email string) (models.User, error) {
	users, err := u.managementCli.User.ListByEmail(ctx, email)
	if err != nil {
		log.Warn().Err(err).Msg("failed when retrieving user")
		return models.User{}, err
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
	// add error not foind
	return models.User{}, errors.New("not found")
}
