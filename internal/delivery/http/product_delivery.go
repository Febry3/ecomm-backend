package http

import (
	"encoding/json"
	"mime/multipart"
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

func (ph *ProductHandler) GetAllCategories(c *gin.Context) {
	categories, err := ph.pr.GetAllCategories(c.Request.Context())
	if err != nil {
		ph.log.Errorf("[ProductDelivery] Get All Categories Error: %v", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": false,
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "categories retrieved successfully",
		"data":    categories,
	})
}

func (ph *ProductHandler) DeleteProductVariant(c *gin.Context) {
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

	err := ph.pr.DeleteProductVariant(c.Request.Context(), c.Param("id"), jwt.SellerID)
	if err != nil {
		ph.log.Errorf("[ProductDelivery] Delete Product Variant Error: %v", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": false,
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "product variant deleted successfully",
	})
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

	if err := c.Request.ParseMultipartForm(32 << 20); err != nil {
		ph.log.Errorf("[ProductDelivery] Parse Form Error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"error":  "failed to parse form data",
		})
		return
	}

	reqData := c.PostForm("data")
	var req dto.CreateProductRequest
	if err := json.Unmarshal([]byte(reqData), &req); err != nil {
		ph.log.Errorf("[ProductDelivery] Invalid JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"error":  "invalid JSON format",
		})
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		ph.log.Errorf("[ProductDelivery] Parse Form Error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"error":  "failed to parse form data",
		})
		return
	}

	var files []*multipart.FileHeader
	if form != nil && form.File != nil && form.File["images"] != nil {
		files = form.File["images"]
	}

	product, err := ph.pr.CreateProduct(c.Request.Context(), req, jwt.SellerID, files)
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

	products, countVariant, totalStock, totalInventoryValue, totalStockAlert, err := ph.pr.GetAllProductsForSeller(c.Request.Context(), jwt.SellerID)
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
		"data": gin.H{
			"total_stock_alert":     totalStockAlert,
			"total_inventory_value": totalInventoryValue,
			"total_stock":           totalStock,
			"count_variant":         countVariant,
			"products":              products,
		},
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

func (ph *ProductHandler) UpdateProduct(c *gin.Context) {
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

	var req dto.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ph.log.Errorf("[ProductDelivery] Bind JSON Error: %v", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"error":  err.Error(),
		})
		return
	}

	product, err := ph.pr.UpdateProduct(c.Request.Context(), req, c.Param("id"), jwt.SellerID)
	if err != nil {
		ph.log.Errorf("[ProductDelivery] Update Product Error: %v", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": false,
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "product updated successfully",
		"data":    product,
	})
}
