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
	ErrConflict                = errors.New("failed to purcase")
	ErrNoStock                 = errors.New("no stock available")
	ErrSessionClosed           = errors.New("session already closed")
	ErrSessionAlreadyStarted   = errors.New("you already started a session")
	ErrGroupBuySessionNotFound = errors.New("group buy session not found")
	ErrSessionFull             = errors.New("session is already full")
)

// Custom error types for HTTP-semantic errors

type NotFoundError struct {
	Message string
}

func (e *NotFoundError) Error() string {
	return e.Message
}

func NewNotFoundError(message string) *NotFoundError {
	return &NotFoundError{Message: message}
}

type BadRequestError struct {
	Message string
}

func (e *BadRequestError) Error() string {
	return e.Message
}

func NewBadRequestError(message string) *BadRequestError {
	return &BadRequestError{Message: message}
}

type ForbiddenError struct {
	Message string
}

func (e *ForbiddenError) Error() string {
	return e.Message
}

func NewForbiddenError(message string) *ForbiddenError {
	return &ForbiddenError{Message: message}
}

type InternalError struct {
	Message string
}

func (e *InternalError) Error() string {
	return e.Message
}

func NewInternalError(message string) *InternalError {
	return &InternalError{Message: message}
}

type UnauthorizedError struct {
	Message string
}

func (e *UnauthorizedError) Error() string {
	return e.Message
}

func NewUnauthorizedError(message string) *UnauthorizedError {
	return &UnauthorizedError{Message: message}
}
