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
	App     *gin.Engine
	Auth    AuthHandler
	User    UserHandler
	Address AddressHandler
	Seller  SellerHandler
	Product ProductHandler
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

	protected := v1.Group("/user", middleware.AuthMiddleware(jwt))
	{
		protected.GET("/test", testUserInline)
		protected.PUT("", routeConfig.User.UpdateUserProfile)
		protected.GET("", routeConfig.User.GetUserProfile)
		protected.POST("/avatar", routeConfig.User.UpdateUserAvatar)
		protected.GET("/address", routeConfig.Address.GetAll)
		protected.POST("/address", routeConfig.Address.Create)
		protected.PUT("/address/:id", routeConfig.Address.Update)
		protected.DELETE("/address/:id", routeConfig.Address.Delete)
	}

	protectedSeller := v1.Group("/seller", middleware.AuthMiddleware(jwt), middleware.RoleMiddleware("seller"))
	{
		protectedSeller.POST("", routeConfig.Seller.RegisterSeller)
		protectedSeller.PUT("", routeConfig.Seller.UpdateSeller)
		protectedSeller.GET("", routeConfig.Seller.GetSeller)

		// Product routes (seller only)
		protectedSeller.POST("/products", routeConfig.Product.CreateProduct)
		protectedSeller.GET("/products", routeConfig.Product.GetAllProductsForSeller)
		protectedSeller.GET("/products/:id", routeConfig.Product.GetProductForSeller)
		protectedSeller.PUT("/products/:id", routeConfig.Product.UpdateProduct)
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
