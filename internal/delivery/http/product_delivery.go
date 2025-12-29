package http

import (
	"encoding/json"
	"mime/multipart"
	"net/http"
	"strconv"

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

func (ph *ProductHandler) GetProductByIDForBuyer(c *gin.Context) {
	product, err := ph.pr.GetProductForBuyer(c.Request.Context(), c.Param("id"))
	if err != nil {
		ph.log.Errorf("[ProductDelivery] Get Product Error: %v", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to get product",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "product retrieved successfully",
		"data":    product,
	})
}

func (ph *ProductHandler) GetAllProductsForBuyer(c *gin.Context) {
	cursor := c.Query("cursor")
	limit, _ := strconv.Atoi(c.Query("limit"))

	if limit == 0 {
		limit = 10
	}

	products, err := ph.pr.GetAllProductsForBuyer(c.Request.Context(), limit, cursor)
	if err != nil {
		ph.log.Errorf("[ProductDelivery] Get All Products Error: %v", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to get all products",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "products retrieved successfully",
		"data":    dto.ToGetProductResponse(products, 2),
	})
}

func (ph *ProductHandler) GetAllCategories(c *gin.Context) {
	categories, err := ph.pr.GetAllCategories(c.Request.Context())
	if err != nil {
		ph.log.Errorf("[ProductDelivery] Get All Categories Error: %v", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to get all categories",
			"error":   err.Error(),
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
			"message": "unauthorized",
			"error":   "unauthorized user",
		})
		return
	}
	jwt := v.(*dto.JwtPayload)

	err := ph.pr.DeleteProductVariant(c.Request.Context(), c.Param("id"), jwt.SellerID)
	if err != nil {
		ph.log.Errorf("[ProductDelivery] Delete Product Variant Error: %v", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to delete product variant",
			"error":   err.Error(),
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
			"message": "unauthorized",
			"error":   "unauthorized user",
		})
		return
	}
	jwt := v.(*dto.JwtPayload)

	if err := c.Request.ParseMultipartForm(32 << 20); err != nil {
		ph.log.Errorf("[ProductDelivery] Parse Form Error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "failed to parse form data",
			"error":   err.Error(),
		})
		return
	}

	reqData := c.PostForm("data")
	var req dto.CreateProductRequest
	if err := json.Unmarshal([]byte(reqData), &req); err != nil {
		ph.log.Errorf("[ProductDelivery] Invalid JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid JSON format",
			"error":   err.Error(),
		})
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		ph.log.Errorf("[ProductDelivery] Parse Form Error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "failed to parse form data",
			"error":   err.Error(),
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
			"message": "failed to create product",
			"error":   err.Error(),
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
			"message": "unauthorized",
			"error":   "unauthorized user",
		})
		return
	}
	jwt := v.(*dto.JwtPayload)

	products, countVariant, totalStock, totalInventoryValue, totalStockAlert, err := ph.pr.GetAllProductsForSeller(c.Request.Context(), jwt.SellerID)
	if err != nil {
		ph.log.Errorf("[ProductDelivery] Get All Products Error: %v", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to get all products",
			"error":   err.Error(),
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
			"message": "unauthorized",
			"error":   "unauthorized user",
		})
		return
	}
	jwt := v.(*dto.JwtPayload)

	product, err := ph.pr.GetProductForSeller(c.Request.Context(), c.Param("id"), jwt.SellerID)
	if err != nil {
		ph.log.Errorf("[ProductDelivery] Get Product Error: %v", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to get product",
			"error":   err.Error(),
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
			"message": "unauthorized",
			"error":   "unauthorized user",
		})
		return
	}
	jwt := v.(*dto.JwtPayload)

	if err := c.Request.ParseMultipartForm(32 << 20); err != nil {
		ph.log.Errorf("[ProductDelivery] Parse Form Error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "failed to parse form data",
			"error":   err.Error(),
		})
		return
	}

	reqData := c.PostForm("data")
	var req dto.UpdateProductRequest
	if err := json.Unmarshal([]byte(reqData), &req); err != nil {
		ph.log.Errorf("[ProductDelivery] Invalid JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid JSON format",
			"error":   err.Error(),
		})
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		ph.log.Errorf("[ProductDelivery] Parse Form Error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "failed to parse form data",
			"error":   err.Error(),
		})
		return
	}

	var files []*multipart.FileHeader
	if form != nil && form.File != nil && form.File["images"] != nil {
		files = form.File["images"]
	}

	product, err := ph.pr.UpdateProduct(c.Request.Context(), req, c.Param("id"), jwt.SellerID, files)
	if err != nil {
		ph.log.Errorf("[ProductDelivery] Update Product Error: %v", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to update product",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "product updated successfully",
		"data":    product,
	})
}

func (ph *ProductHandler) GetProductVariantByID(c *gin.Context) {
	productVariant, err := ph.pr.GetProductVariantByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		ph.log.Errorf("[ProductDelivery] Get Product Variant Error: %v", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to get product variant",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "product variant retrieved successfully",
		"data":    productVariant,
	})
}
