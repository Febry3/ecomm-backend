package http

import (
	"net/http"

	"github.com/febry3/gamingin/internal/dto"
	"github.com/febry3/gamingin/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type SellerHandler struct {
	sc  usecase.SellerUsecaseContract
	log *logrus.Logger
}

func NewSellerHandler(sc usecase.SellerUsecaseContract, log *logrus.Logger) *SellerHandler {
	return &SellerHandler{
		sc:  sc,
		log: log,
	}
}

func (sd *SellerHandler) RegisterSeller(c *gin.Context) {
	v, ok := c.Get("user")
	if !ok {
		sd.log.Error("[UserDelivery] No User in Context")
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "no user in context",
			"error":   ok,
		})
		return
	}
	jwt := v.(*dto.JwtPayload)

	var req dto.SellerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		sd.log.Errorf("[SellerDelivery] Bind JSON Error: %v", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	seller, err := sd.sc.RegisterSeller(c.Request.Context(), req, jwt.ID)
	if err != nil {
		sd.log.Errorf("[SellerDelivery] Register Seller Error: %v", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "seller registered successfully",
		"data":    seller,
	})
}

func (sd *SellerHandler) UpdateSeller(c *gin.Context) {
	v, ok := c.Get("user")
	if !ok {
		sd.log.Error("[UserDelivery] No User in Context")
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "no user in context",
			"error":   ok,
		})
		return
	}
	jwt := v.(*dto.JwtPayload)

	var req dto.SellerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		sd.log.Errorf("[SellerDelivery] Bind JSON Error: %v", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	seller, err := sd.sc.UpdateSeller(c.Request.Context(), req, jwt.ID)
	if err != nil {
		sd.log.Errorf("[SellerDelivery] Update Seller Error: %v", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "seller updated successfully",
		"data":    seller,
	})
}

func (sd *SellerHandler) GetSeller(c *gin.Context) {
	v, ok := c.Get("user")
	if !ok {
		sd.log.Error("[UserDelivery] No User in Context")
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "no user in context",
			"error":   ok,
		})
		return
	}
	jwt := v.(*dto.JwtPayload)

	seller, err := sd.sc.GetSeller(c.Request.Context(), jwt.ID)
	if err != nil {
		sd.log.Errorf("[SellerDelivery] Get Seller Error: %v", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "seller fetched successfully",
		"data":    seller,
	})
}
