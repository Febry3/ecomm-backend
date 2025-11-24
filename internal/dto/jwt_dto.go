package dto

import "github.com/golang-jwt/jwt/v5"

type JwtPayload struct {
	ID       int64  `json:"user_id" `
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}
