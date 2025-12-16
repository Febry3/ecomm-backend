package main

import (
	"github.com/febry3/gamingin/internal/config"
	"github.com/febry3/gamingin/internal/entity"
)

func main() {
	log := config.NewLogrus()
	viperConfig := config.NewViper(log)
	db, _ := config.NewGorm(viperConfig, log)
	_ = db.Migrator().DropTable(&entity.User{}, &entity.AuthProvider{}, &entity.RefreshToken{}, &entity.Address{}, &entity.Seller{}, &entity.Product{}, &entity.ProductVariant{}, &entity.ProductVariantStock{}, &entity.Category{}, &entity.ProductImage{})
	_ = db.AutoMigrate(&entity.User{}, &entity.AuthProvider{}, &entity.RefreshToken{}, &entity.Address{}, &entity.Seller{}, &entity.Product{}, &entity.ProductVariant{}, &entity.ProductVariantStock{}, &entity.Category{}, &entity.ProductImage{})

	categories := []entity.Category{
		{
			Name: "Keyboard",
			Slug: "keyboard",
		},
		{
			Name: "Mouse",
			Slug: "mouse",
		},
		{
			Name: "Headset",
			Slug: "headset",
		},
		{
			Name: "Monitor",
			Slug: "monitor",
		},
		{
			Name: "Mousepad",
			Slug: "mousepad",
		},
		{
			Name: "Controller",
			Slug: "controller", // e.g., Gamepads, Joysticks
		},
		{
			Name: "Microphone",
			Slug: "microphone", // For streamers
		},
		{
			Name: "Gaming Chair",
			Slug: "gaming-chair",
		},
		{
			Name: "Webcam",
			Slug: "webcam",
		},
		{
			Name: "Components",
			Slug: "components", // GPU, RAM, etc.
		},
		{
			Name: "Console",
			Slug: "console",
		},
	}

	for _, category := range categories {
		_ = db.Create(&category)
	}
}
