package main

import (
	"github.com/febry3/gamingin/internal/config"
	"github.com/febry3/gamingin/internal/entity"
)

func main() {
	log := config.NewLogrus()
	viperConfig := config.NewViper(log)
	db, _ := config.NewGorm(viperConfig, log)
	_ = db.Migrator().DropTable(&entity.User{}, &entity.AuthProvider{}, &entity.RefreshToken{})
	_ = db.AutoMigrate(&entity.User{}, &entity.AuthProvider{}, &entity.RefreshToken{})
}
