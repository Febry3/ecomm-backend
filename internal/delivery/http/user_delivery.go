package http

import (
	"net/http"

	"github.com/febry3/gamingin/internal/dto"
	"github.com/febry3/gamingin/internal/helpers"
	"github.com/febry3/gamingin/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type UserHandler struct {
	uc  usecase.UserUsecaseContract
	log *logrus.Logger
}

func NewUserHandler(uc usecase.UserUsecaseContract, log *logrus.Logger) *UserHandler {
	return &UserHandler{uc: uc, log: log}
}

func (uh *UserHandler) GetUserProfile(c *gin.Context) {
	v, ok := c.Get("user")
	if !ok {
		uh.log.Error("[UserDelivery] No User in Context")
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "no user in context",
			"error":   ok,
		})
		return
	}
	jwt := v.(*dto.JwtPayload)

	user, err := uh.uc.GetProfile(c, jwt.ID)
	if err != nil {
		uh.log.Error("[UserDelivery] Error in Getting User: ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "no user in context",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "successfully get user data",
		"data":    user,
	})

}

func (uh *UserHandler) UpdateUserProfile(c *gin.Context) {
	v, ok := c.Get("user")
	if !ok {
		uh.log.Error("[UserDelivery] No User in Context")
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "no user in context",
			"error":   ok,
		})
		return
	}
	jwt := v.(*dto.JwtPayload)

	var userRequest dto.UserRequest
	if err := c.ShouldBindJSON(&userRequest); err != nil {
		uh.log.Errorf("[UserDelivery] Bind Error: %s", err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "invalid request",
			"error":   err.Error(),
		})
		return
	}

	userRequest.UserID = jwt.ID

	updatedUser, err := uh.uc.UpdateProfile(c.Request.Context(), userRequest)
	if err != nil {
		uh.log.Error("[UserDelivery] Update User Error: " + err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "failed to update user data",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "successfully update user data",
		"data":    updatedUser,
	})
}

func (uh *UserHandler) UpdateUserAvatar(c *gin.Context) {
	v, ok := c.Get("user")
	if !ok {
		uh.log.Error("[UserDelivery] No User in Context")
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "no user in context",
			"error":   ok,
		})
		return
	}
	jwt := v.(*dto.JwtPayload)

	fileBytes, err := helpers.GetFileFromContext(c, "file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"error":   err.Error(),
		})
		return
	}

	updatedUser, err := uh.uc.UpdateAvatar(c.Request.Context(), fileBytes, jwt.ID)
	if err != nil {
		uh.log.Printf("Upload error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to upload to Supabase",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "Upload successful",
		"data":    updatedUser,
	})
}
