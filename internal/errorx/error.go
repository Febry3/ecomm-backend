package errors

import "errors"

var (
	ErrEmailTaken   = errors.New("email already registered")
	ErrInvalidLogin = errors.New("invalid email or password")
	ErrTokenInvalid = errors.New("invalid refresh token")
	ErrTokenRevoked = errors.New("refresh token revoked")
	ErrTokenExpired = errors.New("refresh token expired")
)
