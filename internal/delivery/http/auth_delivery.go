package http

import (
	"github.com/febry3/gamingin/internal/dto"
	"github.com/febry3/gamingin/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

type AuthHandler struct {
	uc  usecase.AuthUsecaseContract
	log *logrus.Logger
}

func NewAuthHandler(router *gin.Engine, uc usecase.AuthUsecaseContract, log *logrus.Logger) *AuthHandler {
	return &AuthHandler{uc: uc, log: log}
}

func (a *AuthHandler) Login(c *gin.Context) {
	deviceInfo := c.Request.Header.Get("User-Agent")

	var loginRequest dto.LoginRequest

	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		a.log.Errorf("[AuthDelivery] Bind Error: %s", err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "invalid request",
			"error":   err.Error()},
		)
		return
	}
	loginRequest.DeviceInfo = deviceInfo
	userResponse, refreshToken, err := a.uc.Login(c.Request.Context(), loginRequest)
	if err != nil {
		a.log.Errorf("[AuthDelivery] Login Error: %s", err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "invalid request",
			"error":   err.Error(),
		})
		return
	}

	c.SetCookie("refresh_token", refreshToken, 7*24*60*60, "/api/v1/auth", "localhost", false, true)

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "login success",
		"data":    userResponse,
	})
}
func (a *AuthHandler) Register(c *gin.Context) {
	var registerRequest dto.RegisterRequest
	if err := c.ShouldBindJSON(&registerRequest); err != nil {
		a.log.Errorf("[AuthDelivery] Bind Error: %s", err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "invalid request",
			"error":   err.Error(),
		})
		return
	}
	userResponse, err := a.uc.Register(c.Request.Context(), registerRequest)
	if err != nil {
		a.log.Errorf("[AuthDelivery] Register Error: %s", err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "invalid credentials",
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "register success",
		"data":    userResponse,
	})
}
func (a *AuthHandler) Logout(c *gin.Context) {}
