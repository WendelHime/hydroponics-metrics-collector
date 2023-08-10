package models

import (
	"net/http"

	localErrs "github.com/WendelHime/hydroponics-metrics-collector/internal/shared/errors"

	"github.com/go-playground/validator/v10"
)

// User account fields
type User struct {
	ID            string `json:"-"`
	Name          string `json:"name" validate:"required"`
	Email         string `json:"email" validate:"required,email"`
	Password      string `json:"password" validate:"required,min=8"`
	Role          string `json:"-"`
	EmailVerified bool   `json:"-"`
}

func (u *User) Bind(r *http.Request) error {
	validate := validator.New()
	err := validate.Struct(u)
	if err != nil {
		return localErrs.BadRequestErr.WithErr(err).WithMsg("failed to validate request")
	}
	return nil
}

// Credentials for login
type Credentials struct {
	Email    string
	Password string
	Scope    string
}

type Token struct {
	IDToken      string `json:"id_token,omitempty"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in,omitempty"`
}

func (Token) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
