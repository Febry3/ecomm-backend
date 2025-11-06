package config

import (
	"time"

	"github.com/febry3/gamingin/internal/delivery/http"
	"github.com/febry3/gamingin/internal/helpers"
	"github.com/febry3/gamingin/internal/repository/pg"
	"github.com/febry3/gamingin/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type BootstrapConfig struct {
	DB     *gorm.DB
	App    *gin.Engine
	Log    *logrus.Logger
	Config *viper.Viper
}

func Bootstrap(config *BootstrapConfig) {
	accessTtl, err := time.ParseDuration(config.Config.GetString("jwt.access_ttl"))
	if err != nil {
		config.Log.Fatalf("unable to parse access_ttl: %v", err.Error())
	}
	refreshTtl, err := time.ParseDuration(config.Config.GetString("jwt.refresh_ttl"))
	if err != nil {
		config.Log.Fatalf("unable to parse refresh_ttl: %v", err.Error())
	}

	jwt := helpers.NewJwtService(helpers.JwtConfig{
		Secret:     config.Config.GetString("jwt.secret_key"),
		AccessTTL:  accessTtl,
		RefreshTTL: refreshTtl,
	}, config.Log)

	gauth := NewGoogleAuth(config.Config)

	// setup repo
	userRepository := pg.NewUserRepositoryPg(config.DB, config.Log)
	tokenRepository := pg.NewTokenRepositoryPg(config.DB, config.Log)
	authProviderRepository := pg.NewAuthProvider(config.DB)

	// setup usecase
	authUsecase := usecase.NewAuthUsecase(userRepository, config.Log, *jwt, tokenRepository, authProviderRepository)

	// setup handler
	authHandler := http.NewAuthHandler(config.App, authUsecase, config.Log, gauth)

	routeConfig := http.RouteConfig{
		App:  config.App,
		Auth: *authHandler,
	}

	routeConfig.Init(jwt)
}
