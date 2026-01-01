package http

import (
	"net/http"
	"strconv"

	"github.com/febry3/gamingin/internal/dto"
	"github.com/febry3/gamingin/internal/errorx"
	"github.com/febry3/gamingin/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type OrderHandler struct {
	orderUsecase usecase.OrderUsecaseContract
	log          *logrus.Logger
}

func NewOrderHandler(orderUsecase usecase.OrderUsecaseContract, log *logrus.Logger) *OrderHandler {
	return &OrderHandler{
		orderUsecase: orderUsecase,
		log:          log,
	}
}

// CreateDirectOrder handles POST /user/orders - direct buy flow
func (h *OrderHandler) CreateDirectOrder(c *gin.Context) {
	claims, err := getUserClaims(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, errorResponse("Unauthorized"))
		return
	}

	var request dto.CreateOrderRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err.Error()))
		return
	}

	order, err := h.orderUsecase.CreateDirectOrder(c.Request.Context(), claims.ID, &request)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Order created successfully",
		"data":    order,
	})
}

// CreateGroupBuyOrder handles POST /user/orders/group-buy
func (h *OrderHandler) CreateGroupBuyOrder(c *gin.Context) {
	claims, err := getUserClaims(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, errorResponse("Unauthorized"))
		return
	}

	var request dto.CreateGroupBuyOrderRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err.Error()))
		return
	}

	order, err := h.orderUsecase.CreateGroupBuyOrder(c.Request.Context(), claims.ID, &request)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Group buy order created successfully",
		"data":    order,
	})
}

// GetOrders handles GET /user/orders - list user's orders
func (h *OrderHandler) GetOrders(c *gin.Context) {
	claims, err := getUserClaims(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, errorResponse("Unauthorized"))
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	orders, err := h.orderUsecase.GetOrders(c.Request.Context(), claims.ID, page, limit)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Orders retrieved successfully",
		"data":    orders,
	})
}

// GetOrderByID handles GET /user/orders/:id
func (h *OrderHandler) GetOrderByID(c *gin.Context) {
	claims, err := getUserClaims(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, errorResponse("Unauthorized"))
		return
	}

	orderID := c.Param("id")
	if orderID == "" {
		c.JSON(http.StatusBadRequest, errorResponse("Order ID is required"))
		return
	}

	order, err := h.orderUsecase.GetOrderByID(c.Request.Context(), claims.ID, orderID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Order retrieved successfully",
		"data":    order,
	})
}

// HandlePaymentNotification handles POST /payments/webhook - Midtrans callback
func (h *OrderHandler) HandlePaymentNotification(c *gin.Context) {
	var notification dto.MidtransNotification
	if err := c.ShouldBindJSON(&notification); err != nil {
		h.log.Warnf("Invalid webhook payload: %v", err)
		c.JSON(http.StatusBadRequest, errorResponse("Invalid payload"))
		return
	}

	h.log.Infof("Received payment notification: OrderID=%s, Status=%s", notification.OrderID, notification.TransactionStatus)

	err := h.orderUsecase.HandlePaymentNotification(c.Request.Context(), &notification)
	if err != nil {
		handleError(c, err)
		return
	}

	// Midtrans expects 200 OK to acknowledge receipt
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// Helper functions

func errorResponse(message string) gin.H {
	return gin.H{"error": message}
}

func handleError(c *gin.Context, err error) {
	switch e := err.(type) {
	case *errorx.NotFoundError:
		c.JSON(http.StatusNotFound, errorResponse(e.Error()))
	case *errorx.BadRequestError:
		c.JSON(http.StatusBadRequest, errorResponse(e.Error()))
	case *errorx.ForbiddenError:
		c.JSON(http.StatusForbidden, errorResponse(e.Error()))
	case *errorx.InternalError:
		c.JSON(http.StatusInternalServerError, errorResponse(e.Error()))
	default:
		c.JSON(http.StatusInternalServerError, errorResponse("Internal server error"))
	}
}

func getUserClaims(c *gin.Context) (*dto.JwtPayload, error) {
	v, exists := c.Get("user")
	if !exists {
		return nil, errorx.NewUnauthorizedError("User not found in context")
	}
	claims, ok := v.(*dto.JwtPayload)
	if !ok {
		return nil, errorx.NewUnauthorizedError("Invalid user claims")
	}
	return claims, nil
}
