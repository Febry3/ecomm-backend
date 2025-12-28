package http

import (
	"net/http"

	"github.com/febry3/gamingin/internal/dto"
	"github.com/febry3/gamingin/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type GroupBuyHandler struct {
	pu  usecase.GroupBuyUsecaseContract
	log *logrus.Logger
}

func NewGroupBuyHandler(pu usecase.GroupBuyUsecaseContract, log *logrus.Logger) *GroupBuyHandler {
	return &GroupBuyHandler{pu: pu, log: log}
}

func (gh *GroupBuyHandler) ChangeGroupBuySessionStatus(c *gin.Context) {
	v, ok := c.Get("user")
	if !ok {
		gh.log.Error("[ProductDelivery] No User in Context")
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "unauthorized",
			"error":   "unauthorized user",
		})
		return
	}
	jwt := v.(*dto.JwtPayload)

	var req dto.ChangeStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		gh.log.Error("[ProductDelivery] ChangeGroupBuySessionStatus failed: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to change group buy session status",
			"error":   err.Error(),
		})
		return
	}

	err := gh.pu.ChangeGroupBuySessionStatus(c.Request.Context(), req.SessionID, req.Status, jwt.SellerID)
	if err != nil {
		gh.log.Error("[ProductDelivery] ChangeGroupBuySessionStatus failed: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to change group buy session status",
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "success to change group buy status",
	})
}

func (gh *GroupBuyHandler) GetAllGroupBuySessionForSeller(c *gin.Context) {
	v, ok := c.Get("user")
	if !ok {
		gh.log.Error("[ProductDelivery] No User in Context")
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "unauthorized",
			"error":   "unauthorized user",
		})
		return
	}
	jwt := v.(*dto.JwtPayload)

	groupBuySessions, err := gh.pu.GetAllGroupBuySessionForSeller(c, jwt.SellerID)
	if err != nil {
		gh.log.Error("[ProductDelivery] GetAllGroupBuySessionForSeller failed: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to get group buy session for seller",
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "success",
		"data":    groupBuySessions,
	})
}

func (gh *GroupBuyHandler) GetAllGroupBuySessionForBuyer(c *gin.Context) {
	groupBuySessions, err := gh.pu.GetAllGroupBuySessionForBuyer(c)
	if err != nil {
		gh.log.Error("[ProductDelivery] GetAllGroupBuySessionForBuyer failed: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to get group buy session for buyer",
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "success",
		"data":    groupBuySessions,
	})
}

func (gh *GroupBuyHandler) CreateGroupBuySession(c *gin.Context) {
	v, ok := c.Get("user")
	if !ok {
		gh.log.Error("[ProductDelivery] No User in Context")
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "unauthorized",
			"error":   "unauthorized user",
		})
		return
	}
	jwt := v.(*dto.JwtPayload)

	var req *dto.GroupBuySessionRequest
	if err := c.ShouldBind(&req); err != nil {
		gh.log.Error("[ProductDelivery] CreateGroupBuySession failed: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to create group buy session",
			"error":   err.Error(),
		})
		return
	}

	response, err := gh.pu.CreateGroupBuySession(c.Request.Context(), req, jwt.SellerID)
	if err != nil {
		gh.log.Error("[ProductDelivery] CreateGroupBuySession failed: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to create group buy session",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "success",
		"data":    response,
	})
}

func (gh *GroupBuyHandler) CreateBuyerSession(c *gin.Context) {
	v, ok := c.Get("user")
	if !ok {
		gh.log.Error("[ProductDelivery] No User in Context")
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "unauthorized",
			"error":   "unauthorized user",
		})
		return
	}
	jwt := v.(*dto.JwtPayload)

	var req dto.CreateBuyerGroupSessionRequest
	req.OrganizerUserID = jwt.ID
	if err := c.ShouldBind(&req); err != nil {
		gh.log.Error("[ProductDelivery] CreateBuyerSession failed: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to create buyer session",
			"error":   err.Error(),
		})
		return
	}

	err := gh.pu.CreateBuyerSession(c.Request.Context(), &req)
	if err != nil {
		gh.log.Error("[ProductDelivery] CreateBuyerSession failed: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to create buyer session",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "success",
	})
}
