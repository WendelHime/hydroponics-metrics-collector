package errors

import "fmt"

type Error struct {
	StatusCode        int
	StatusDescription string
	Msg               string
	Details           map[string]any
	Err               error
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
