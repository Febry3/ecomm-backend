package http

import (
	"net/http"
	"time"

	"github.com/febry3/gamingin/internal/delivery/http/middleware"
	"github.com/febry3/gamingin/internal/dto"
	"github.com/febry3/gamingin/internal/helpers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type RouteConfig struct {
	App      *gin.Engine
	Auth     AuthHandler
	User     UserHandler
	Address  AddressHandler
	Seller   SellerHandler
	Product  ProductHandler
	GroupBuy GroupBuyHandler
}

func (routeConfig *RouteConfig) Init(jwt *helpers.JwtService) {
	corsConf := cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization", "Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	routeConfig.App.Use(cors.New(corsConf))

	v1 := routeConfig.App.Group("/v1/api")

	auth := v1.Group("/auth")
	auth.POST("/login", routeConfig.Auth.Login)
	auth.POST("/register", routeConfig.Auth.Register)
	auth.POST("/logout", routeConfig.Auth.Logout)
	auth.POST("/refresh", routeConfig.Auth.RefreshToken)
	auth.POST("/google", routeConfig.Auth.LoginOrRegisterWithGoogle)

	product := v1.Group("/product")
	{
		product.GET("/categories", routeConfig.Product.GetAllCategories)
		product.GET("", routeConfig.Product.GetAllProductsForBuyer)
		product.GET("/:id", routeConfig.Product.GetProductByIDForBuyer)
		product.GET("/variants/:id", routeConfig.Product.GetProductVariantByID)
	}

	protected := v1.Group("", middleware.AuthMiddleware(jwt))
	{
		protected.POST("/group-buy", routeConfig.GroupBuy.CreateBuyerSession)
		protected.GET("/group-buy/:sessionId", routeConfig.GroupBuy.GetSessionForBuyerByCode)
		protected.POST("/group-buy/:sessionId/join", routeConfig.GroupBuy.JoinSession)
	}

	protectedUser := v1.Group("/user", middleware.AuthMiddleware(jwt))
	{
		protectedUser.GET("/test", testUserInline)
		protectedUser.PUT("", routeConfig.User.UpdateUserProfile)
		protectedUser.GET("", routeConfig.User.GetUserProfile)
		protectedUser.POST("/avatar", routeConfig.User.UpdateUserAvatar)
		protectedUser.GET("/address", routeConfig.Address.GetAll)
		protectedUser.POST("/address", routeConfig.Address.Create)
		protectedUser.PUT("/address/:id", routeConfig.Address.Update)
		protectedUser.DELETE("/address/:id", routeConfig.Address.Delete)
	}

	protectedSeller := v1.Group("/seller", middleware.AuthMiddleware(jwt))
	{
		protectedSeller.POST("", routeConfig.Seller.RegisterSeller)

		sellerRole := protectedSeller.Group("", middleware.RoleMiddleware("seller"))
		{
			sellerRole.PUT("", routeConfig.Seller.UpdateSeller)
			sellerRole.GET("", routeConfig.Seller.GetSeller)

			// Product routes (seller only)
			sellerRole.POST("/products", routeConfig.Product.CreateProduct)
			sellerRole.GET("/products", routeConfig.Product.GetAllProductsForSeller)
			sellerRole.GET("/products/:id", routeConfig.Product.GetProductForSeller)
			sellerRole.PUT("/products/:id", routeConfig.Product.UpdateProduct)

			// Product Variant
			sellerRole.DELETE("/products/variants/:id", routeConfig.Product.DeleteProductVariant)

			// Group Buy
			sellerRole.POST("/group-buy", routeConfig.GroupBuy.CreateGroupBuySession)
			sellerRole.GET("/group-buy", routeConfig.GroupBuy.GetAllGroupBuySessionForSeller)
			sellerRole.PATCH("/group-buy/status", routeConfig.GroupBuy.ChangeGroupBuySessionStatus)

		}

	}
}

func testUserInline(c *gin.Context) {
	v, ok := c.Get("user")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no user in context"})
		return
	}

	claims := v.(*dto.JwtPayload)
	c.JSON(http.StatusOK, gin.H{
		"user_id":  claims.ID,
		"username": claims.Username,
		"email":    claims.Email,
		"role":     claims.Role,
		"exp":      claims.ExpiresAt,
		"iat":      claims.IssuedAt,
	})
}
