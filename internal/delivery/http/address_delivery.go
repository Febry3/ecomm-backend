package http

import (
	"net/http"

	"github.com/febry3/gamingin/internal/dto"
	"github.com/febry3/gamingin/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type AddressHandler struct {
	address usecase.AddressUsecaseContract
	user    usecase.UserUsecaseContract
	log     *logrus.Logger
}

func NewAddressHandler(address usecase.AddressUsecaseContract, user usecase.UserUsecaseContract, log *logrus.Logger) *AddressHandler {
	return &AddressHandler{
		address: address,
		user:    user,
		log:     log,
	}
}

func (ah *AddressHandler) GetAll(c *gin.Context) {
	v, ok := c.Get("user")
	if !ok {
		ah.log.Error("[UserDelivery] No User in Context")
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "no user in context",
			"error":   ok,
		})
		return
	}
	jwt := v.(*dto.JwtPayload)

	addresses, err := ah.address.GetAll(c.Request.Context(), jwt.ID)
	if err != nil {
		ah.log.Errorf("[AddressHandler] Get All Address Error: %s", err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "invalid request",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "get all address success",
		"data":    addresses,
	})
}

func (ah *AddressHandler) Create(c *gin.Context) {
	v, ok := c.Get("user")
	if !ok {
		ah.log.Error("[UserDelivery] No User in Context")
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "no user in context",
			"error":   ok,
		})
		return
	}
	jwt := v.(*dto.JwtPayload)

	var addressRequest dto.AddressRequest
	if err := c.ShouldBindJSON(&addressRequest); err != nil {
		ah.log.Errorf("[UserDelivery] Bind Error: %s", err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "invalid request",
			"error":   err.Error(),
		})
		return
	}

	address, err := ah.address.Create(c.Request.Context(), addressRequest, jwt.ID)
	if err != nil {
		ah.log.Errorf("[AddressHandler] Create Address Error: %s", err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "invalid request",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "create address success",
		"data":    address,
	})
}

func (ah *AddressHandler) Update(c *gin.Context) {
	id := c.Param("id")
	v, ok := c.Get("user")
	if !ok {
		ah.log.Error("[UserDelivery] No User in Context")
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "no user in context",
			"error":   ok,
		})
		return
	}
	jwt := v.(*dto.JwtPayload)

	var addressRequest dto.AddressRequest
	if err := c.ShouldBindJSON(&addressRequest); err != nil {
		ah.log.Errorf("[UserDelivery] Bind Error: %s", err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "invalid request",
			"error":   err.Error(),
		})
		return
	}

	address, err := ah.address.Update(c.Request.Context(), addressRequest, id, jwt.ID)
	if err != nil {
		ah.log.Errorf("[AddressHandler] Update Address Error: %s", err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "invalid request",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "update address success",
		"data":    address,
	})
}

func (ah *AddressHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	v, ok := c.Get("user")
	if !ok {
		ah.log.Error("[UserDelivery] No User in Context")
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "no user in context",
			"error":   ok,
		})
		return
	}
	jwt := v.(*dto.JwtPayload)

	err := ah.address.Delete(c.Request.Context(), id, jwt.ID)
	if err != nil {
		ah.log.Errorf("[AddressHandler] Delete Address Error: %s", err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "invalid request",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "delete address success",
	})
}
