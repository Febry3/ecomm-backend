package main

import (
	"fmt"

	"github.com/febry3/gamingin/internal/config"
)

func main() {
	log := config.NewLogrus()
	viperConfig := config.NewViper(log)
	app := config.NewGin(viperConfig)
	db, err := config.NewGorm(viperConfig, log)
	// db.AutoMigrate(&entity.User{}, &entity.AuthProvider{}, &entity.RefreshToken{})

	if err != nil {
		log.Errorf("unable to connect database: %v", err.Error())
	}

	config.Bootstrap(&config.BootstrapConfig{
		Log:    log,
		App:    app,
		Config: viperConfig,
		DB:     db,
	})

	if err := app.Run(fmt.Sprintf(":%d", viperConfig.GetInt("app.port"))); err != nil {
		log.Fatalf("could not start server: %v", err.Error())
	}
}
