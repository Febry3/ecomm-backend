package http

import (
	"net/http"

	"github.com/febry3/gamingin/internal/dto"
	"github.com/febry3/gamingin/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type ProductHandler struct {
	pr  usecase.ProductUsecaseContract
	log *logrus.Logger
}

func NewProductHandler(pr usecase.ProductUsecaseContract, log *logrus.Logger) *ProductHandler {
	return &ProductHandler{
		pr:  pr,
		log: log,
	}
}

func (ph *ProductHandler) CreateProduct(c *gin.Context) {
	v, ok := c.Get("user")
	if !ok {
		ph.log.Error("[ProductDelivery] No User in Context")
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  false,
			"message": "unauthorized",
		})
		return
	}
	jwt := v.(*dto.JwtPayload)

	// Validate user is a seller
	if jwt.Role != "seller" || jwt.SellerID == 0 {
		ph.log.Error("[ProductDelivery] User is not a seller")
		c.JSON(http.StatusForbidden, gin.H{
			"status":  false,
			"message": "only sellers can create products",
		})
		return
	}

	var req dto.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ph.log.Errorf("[ProductDelivery] Bind JSON Error: %v", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"error":  err.Error(),
		})
		return
	}

	// Use SellerID from JWT, not user ID
	product, err := ph.pr.CreateProduct(c.Request.Context(), req, jwt.SellerID)
	if err != nil {
		ph.log.Errorf("[ProductDelivery] Create Product Error: %v", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": false,
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  true,
		"message": "product created successfully",
		"data":    product,
	})
}

func (ph *ProductHandler) GetAllProductsForSeller(c *gin.Context) {
	v, ok := c.Get("user")
	if !ok {
		ph.log.Error("[ProductDelivery] No User in Context")
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  false,
			"message": "unauthorized",
		})
		return
	}
	jwt := v.(*dto.JwtPayload)

	products, err := ph.pr.GetAllProductsForSeller(c.Request.Context(), jwt.SellerID)
	if err != nil {
		ph.log.Errorf("[ProductDelivery] Get All Products Error: %v", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": false,
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "products retrieved successfully",
		"data":    products,
	})
}

func (ph *ProductHandler) GetProductForSeller(c *gin.Context) {
	v, ok := c.Get("user")
	if !ok {
		ph.log.Error("[ProductDelivery] No User in Context")
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  false,
			"message": "unauthorized",
		})
		return
	}
	jwt := v.(*dto.JwtPayload)

	product, err := ph.pr.GetProductForSeller(c.Request.Context(), c.Param("id"), jwt.SellerID)
	if err != nil {
		ph.log.Errorf("[ProductDelivery] Get Product Error: %v", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": false,
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "product retrieved successfully",
		"data":    product,
	})
}
