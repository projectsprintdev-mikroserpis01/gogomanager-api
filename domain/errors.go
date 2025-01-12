package domain

import (
	"errors"
	"net/http"
)

type SerializableError interface {
	Serialize() any
}

type RequestError struct {
	StatusCode int
	Err        error
}

func (r *RequestError) Error() string {
	return r.Err.Error()
}

var ErrNotFound = &RequestError{
	StatusCode: http.StatusNotFound,
	Err:        errors.New("something not found"),
}

var ErrNoAPIKey = &RequestError{
	StatusCode: http.StatusUnauthorized,
	Err:        errors.New("no api key provided"),
}

var ErrInvalidAPIKey = &RequestError{
	StatusCode: http.StatusUnauthorized,
	Err:        errors.New("invalid api key"),
}

var ErrUserNotFound = &RequestError{
	StatusCode: http.StatusNotFound,
	Err:        errors.New("user not found"),
}

var ErrUserEmailAlreadyExists = &RequestError{
	StatusCode: http.StatusConflict,
	Err:        errors.New("user email already exists"),
}

var ErrNoBearerToken = &RequestError{
	StatusCode: http.StatusUnauthorized,
	Err:        errors.New("no bearer token provided"),
}

var ErrInvalidBearerToken = &RequestError{
	StatusCode: http.StatusUnauthorized,
	Err:        errors.New("invalid bearer token"),
}

var ErrExpiredBearerToken = &RequestError{
	StatusCode: http.StatusUnauthorized,
	Err:        errors.New("expired bearer token"),
}

var ErrBearerTokenNotActive = &RequestError{
	StatusCode: http.StatusUnauthorized,
	Err:        errors.New("bearer token not active"),
}

var ErrEmailNotFound = &RequestError{
	StatusCode: http.StatusNotFound,
	Err:        errors.New("email not found"),
}

var ErrCredentialsNotMatch = &RequestError{
	StatusCode: http.StatusUnauthorized,
	Err:        errors.New("credentials do not match"),
}

var ErrRoleCantAccessResource = &RequestError{
	StatusCode: http.StatusForbidden,
	Err:        errors.New("role can't access resource"),
}

var ErrFileSizeLimitExceeded = &RequestError{
	StatusCode: http.StatusBadRequest,
	Err:        errors.New("file size limit exceeded"),
}

var ErrInvalidFileExtension = &RequestError{
	StatusCode: http.StatusBadRequest,
	Err:        errors.New("invalid file extension"),
}

var ErrFileNotFound = &RequestError{
	StatusCode: http.StatusBadRequest,
	Err:        errors.New("file not found"),
}

var ErrInvalidMimeType = &RequestError{
	StatusCode: http.StatusBadRequest,
	Err:        errors.New("invalid mime type"),
}
