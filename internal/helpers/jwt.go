package helpers

import (
	"errors"
	"log"
	"time"

	"github.com/febry3/gamingin/internal/dto"
	"github.com/febry3/gamingin/internal/errorx"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
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

func (j *JwtService) IssueAccessToken(payload dto.JwtPayload) string {
	now := time.Now().UTC()
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

func (j *JwtService) VerifyToken(tokenString string) (*dto.JwtPayload, error) {
	token, err := jwt.ParseWithClaims(tokenString, &dto.JwtPayload{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(j.config.Secret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, errorx.ErrTokenExpired
		}
		return nil, errorx.ErrParseToken
	}

	claims, ok := token.Claims.(*dto.JwtPayload)
	if !ok || !token.Valid {
		return nil, errorx.ErrTokenInvalid
	}

	return claims, nil
}
