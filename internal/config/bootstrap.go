package config

import (
	"time"

	"github.com/febry3/gamingin/internal/delivery/http"
	"github.com/febry3/gamingin/internal/helpers"
	"github.com/febry3/gamingin/internal/infra/storage"
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

	// service
	gauth := NewGoogleAuth(config.Config)
	supabaseConfig := NewSupabaseConfig(config.Config)
	storage := storage.NewSupabaseHttpRepo(supabaseConfig)

	// setup repo
	userRepository := pg.NewUserRepositoryPg(config.DB, config.Log)
	tokenRepository := pg.NewTokenRepositoryPg(config.DB, config.Log)
	authProviderRepository := pg.NewAuthProvider(config.DB)
	addressRepository := pg.NewAddressRepositoryPg(config.DB)
	sellerRepository := pg.NewSellerRepositoryPg(config.DB, config.Log)
	productRepository := pg.NewProductRepositoryPg(config.DB)
	variantRepository := pg.NewProductVariantRepositoryPg(config.DB)
	stockRepository := pg.NewProductVariantStockRepositoryPg(config.DB)
	categoryRepository := pg.NewCategoryRepositoryPg(config.DB)
	productImageRepository := pg.NewProductImageRepositoryPg(config.DB)
	txManager := pg.NewTxManager(config.DB)

	// setup usecase
	authUsecase := usecase.NewAuthUsecase(userRepository, config.Log, *jwt, tokenRepository, authProviderRepository, sellerRepository)
	userUsecase := usecase.NewUserUsecase(userRepository, config.Log, storage, sellerRepository)
	addressUsecase := usecase.NewAddressUsecase(addressRepository, userRepository, config.Log)
	sellerUsecase := usecase.NewSellerUsecase(sellerRepository, userRepository, txManager, config.Log, storage)
	productUsecase := usecase.NewProductUsecase(productRepository, variantRepository, stockRepository, sellerRepository, categoryRepository, productImageRepository, storage, txManager, config.Log)

	// setup handler
	authHandler := http.NewAuthHandler(authUsecase, config.Log, gauth)
	userHandler := http.NewUserHandler(userUsecase, config.Log)
	addressHandler := http.NewAddressHandler(addressUsecase, userUsecase, config.Log)
	sellerHandler := http.NewSellerHandler(sellerUsecase, config.Log)
	productHandler := http.NewProductHandler(productUsecase, config.Log)

	routeConfig := http.RouteConfig{
		App:     config.App,
		Auth:    *authHandler,
		User:    *userHandler,
		Address: *addressHandler,
		Seller:  *sellerHandler,
		Product: *productHandler,
	}

	routeConfig.Init(jwt)
}
