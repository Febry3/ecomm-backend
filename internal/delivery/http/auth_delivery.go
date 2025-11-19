package http

import (
	"net/http"

	"github.com/febry3/gamingin/internal/helpers"
	"golang.org/x/oauth2"

	"github.com/febry3/gamingin/internal/dto"
	"github.com/febry3/gamingin/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type AuthHandler struct {
	uc    usecase.AuthUsecaseContract
	log   *logrus.Logger
	gauth *oauth2.Config
}

func NewAuthHandler(uc usecase.AuthUsecaseContract, log *logrus.Logger, gauth *oauth2.Config) *AuthHandler {
	return &AuthHandler{uc: uc, log: log, gauth: gauth}
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

	c.SetCookie("refresh_token", refreshToken, 7*24*60*60, "*", "localhost", false, true)

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

func (a *AuthHandler) RefreshToken(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	a.log.Debug("check refresh token", refreshToken)
	if err != nil {
		a.log.Errorf("[AuthDelivery] Get Cookie Error: %s", err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "invalid request",
			"error":   err.Error(),
		})
		return
	}

	newAccessToken, err := a.uc.RefreshAccessToken(c.Request.Context(), refreshToken)
	if err != nil {
		a.log.Errorf("[AuthDelivery] Refresh Token Error: %s", err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "invalid request",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":       true,
		"message":      "token refreshed successfully",
		"access_token": newAccessToken,
	})
}

func (a *AuthHandler) LoginOrRegisterWithGoogle(c *gin.Context) {
	var googleLoginRequest dto.LoginWithGoogleRequest
	if err := c.ShouldBindJSON(&googleLoginRequest); err != nil {
		a.log.Errorf("[AuthDelivery] Bind Error: %s", err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "invalid request",
			"error":   err.Error(),
		})
		return
	}

	token, err := a.gauth.Exchange(c.Request.Context(), googleLoginRequest.Code)
	if err != nil {
		a.log.Errorf("[AuthDelivery] Exchange Error: %s", err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "invalid request",
			"error":   err.Error(),
		})
		return
	}

	userInfo, err := helpers.GetGoogleUserInfo(c.Request.Context(), token, a.gauth)
	if err != nil {
		a.log.Errorf("[AuthDelivery] Get Google User Info error: %s", err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "failed to get google user info",
			"error":   err.Error(),
		})
	}

	userInfo.DeviceInfo = c.Request.Header.Get("User-Agent")

	userResponse, refreshToken, err := a.uc.LoginOrRegisterWithGoogle(c.Request.Context(), userInfo)
	if err != nil {
		a.log.Errorf("[AuthDelivery] Login OrRegisterWithGoogle Error: %s", err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "failed to login with google",
			"error":   err.Error(),
		})
		return
	}

	c.SetCookie("refresh_token", refreshToken, 7*24*60*60, "*", "localhost", false, true)

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "login success",
		"data":    userResponse,
	})
}

func (a *AuthHandler) Logout(c *gin.Context) {
	accessToken, err := c.Cookie("refresh_token")

	if err != nil {
		a.log.Errorf("[AuthDelivery] Get Cookie Error: %s", err.Error())
		c.SetCookie("refresh_token", "", -1, "*", "localhost", false, true)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status":  true,
			"message": "logout successfully",
		})
		return
	}

	err = a.uc.Logout(c.Request.Context(), accessToken)
	if err != nil {
		a.log.Errorf("[AuthDelivery] Logout Error: %s", err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "failed to logout",
			"error":   err.Error(),
		})
		return
	}

	c.SetCookie("refresh_token", "", -1, "*", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "logout success",
	})
}
