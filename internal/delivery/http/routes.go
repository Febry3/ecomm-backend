package http

import (
	"github.com/febry3/gamingin/internal/delivery/http/middleware"
	"github.com/gin-gonic/gin"
)

type RouteConfig struct {
	App  *gin.Engine
	Auth AuthHandler
}

func (routeConfig *RouteConfig) Init() {
	routeConfig.App.Use(middleware.CORSMiddleware())
	v1 := routeConfig.App.Group("/v1/api")

	auth := v1.Group("/auth")
	auth.POST("/login", routeConfig.Auth.Login)
	auth.POST("/register", routeConfig.Auth.Register)
	auth.POST("/logout", routeConfig.Auth.Logout)
}
