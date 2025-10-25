package http

import "github.com/gin-gonic/gin"

type RouteConfig struct {
	App  *gin.Engine
	Auth AuthHandler
}

func (routeConfig *RouteConfig) Init() {
	v1 := routeConfig.App.Group("/v1/api")

	v1.POST("/login", routeConfig.Auth.Login)
	v1.POST("/register", routeConfig.Auth.Register)
	v1.POST("/logout", routeConfig.Auth.Logout)
}
