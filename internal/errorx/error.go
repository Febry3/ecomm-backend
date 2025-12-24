package errorx

import "errors"

var (
	ErrEmailTaken         = errors.New("email already registered")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidLogin       = errors.New("invalid email or password")
	ErrUserNotFound       = errors.New("user not found")

	ErrTokenEmpty   = errors.New("token is empty")
	ErrTokenInvalid = errors.New("invalid refresh token")
	ErrTokenRevoked = errors.New("refresh token revoked")
	ErrTokenExpired = errors.New("refresh token expired")
	ErrParseToken   = errors.New("failed to parse token")

	ErrSellerAlreadyExists = errors.New("seller already exists")

	ErrInsufficientStock = errors.New("product variant stock is not enough")

	// related to group buying feature
	ErrConflict      = errors.New("failed to purcase")
	ErrNoStock       = errors.New("no stock available")
	ErrSessionClosed = errors.New("session already closed")
)
