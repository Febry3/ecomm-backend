package helpers

import (
	"errors"
	"github.com/febry3/gamingin/internal/dto"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"log"
	"time"
)

type JwtConfig struct {
	Secret     string
	AccessTTL  time.Duration
	RefreshTTL time.Duration
}

type JwtService struct {
	config JwtConfig
	log    *logrus.Logger
}

func NewJwtService(config JwtConfig, log *logrus.Logger) *JwtService {
	return &JwtService{config: config, log: log}
}

func (j *JwtService) IssueJwt(payload dto.JwtPayload) string {
	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": payload.Username,
		"email":    payload.Email,
		"user_id":  payload.ID,
		"role":     payload.Role,
		"exp":      jwt.NewNumericDate(now.Add(j.config.AccessTTL)),
		"iat":      jwt.NewNumericDate(now),
	})

	signedToken, err := token.SignedString([]byte(j.config.Secret))
	if err != nil {
		log.Fatalf("signing token error: %v", err)
	}
	return signedToken
}

func (j *JwtService) VerifyJwt(tokenString string) (*dto.JwtPayload, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return j.config.Secret, nil
	})

	if err != nil {
		return nil, errors.New("failed to parse token")
	}

	claims, ok := token.Claims.(*dto.JwtPayload)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
