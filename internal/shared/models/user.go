package models

// User account fields
type User struct {
	ID            string `json:"-"`
	Name          string `json:"name" validate:"required"`
	Email         string `json:"email" validate:"required"`
	Password      string `json:"password" validate:"required"`
	Role          string `json:"-"`
	EmailVerified bool   `json:"-"`
}

// Credentials for login
type Credentials struct {
	Email    string
	Password string
	Scope    string
}
