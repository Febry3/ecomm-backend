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

func (gh *GroupBuyHandler) GetAllGroupBuySessionForSeller(c *gin.Context) {
	v, ok := c.Get("user")
	if !ok {
		gh.log.Error("[ProductDelivery] No User in Context")
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  false,
			"message": "unauthorized",
		})
		return
	}
	jwt := v.(*dto.JwtPayload)

	groupBuySessions, err := gh.pu.GetAllGroupBuySessionForSeller(c, jwt.SellerID)
	if err != nil {
		gh.log.Error("[ProductDelivery] GetAllGroupBuySessionForSeller failed: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": err.Error(),
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
			"status":  false,
			"message": "internal server error",
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
			"status":  false,
			"message": "unauthorized",
		})
		return
	}
	jwt := v.(*dto.JwtPayload)

	var req *dto.GroupBuySessionRequest
	if err := c.ShouldBind(&req); err != nil {
		gh.log.Error("[ProductDelivery] CreateGroupBuySession failed: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "internal server error",
		})
		return
	}

	response, err := gh.pu.CreateGroupBuySession(c.Request.Context(), req, jwt.SellerID)
	if err != nil {
		gh.log.Error("[ProductDelivery] CreateGroupBuySession failed: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "success",
		"data":    response,
	})
}
