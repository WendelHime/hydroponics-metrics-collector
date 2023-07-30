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
	Email         string `json:"email" validate:"required"`
	Password      string `json:"password" validate:"required"`
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
