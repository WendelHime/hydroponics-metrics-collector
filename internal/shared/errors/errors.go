package errors

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/render"
	"github.com/rs/zerolog/log"
)

type Error struct {
	StatusCode        int            `json:"status"`
	StatusDescription string         `json:"description"`
	Msg               string         `json:"-"`
	Details           map[string]any `json:"-"`
	Err               error          `json:"-"`
}

func newError(statusCode int, statusDescription string) *Error {
	return &Error{StatusCode: statusCode, StatusDescription: statusDescription, Details: make(map[string]any)}
}

func (e *Error) WithMsg(message string) *Error {
	e.Msg = message
	return e
}

func (e *Error) WithDetails(key string, value any) *Error {
	e.Details[key] = value
	return e
}

func (e *Error) WithErr(err error) *Error {
	e.Err = err
	return e
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s: %s, %+v", e.StatusDescription, e.Msg, e.Details)
}

func (e *Error) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.StatusCode)
	return nil
}

func RenderErr(w http.ResponseWriter, r *http.Request, err error) {
	var apiErr *Error
	if errors.As(err, &apiErr) {
		log.Warn().Err(err).Msg("request failed")
		render.Render(w, r, apiErr)
		return
	}

	log.Warn().Err(err).Msg("internal server error")
	render.Status(r, http.StatusInternalServerError)
}

// InternalServerErr used for random/unexpected internal errors
var InternalServerErr *Error = newError(500, "internal server error")

// NotFoundErr when the resource was not found
var NotFoundErr *Error = newError(404, "not found")

// AlreadyExistsErr when the resource already exists
var AlreadyExistsErr *Error = newError(409, "resource already exists")

// BadRequestErr for malformed requests
var BadRequestErr *Error = newError(400, "bad request")

// ForbiddenErr used when the user is inactive or doesn't the expected permissions
var ForbiddenErr *Error = newError(403, "forbidden")

// UnauthorizedErr used when the provided token is invalid
var UnauthorizedErr *Error = newError(401, "unauthorized")
