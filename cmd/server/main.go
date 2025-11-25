package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/febry3/gamingin/internal/config"
)

func main() {
	log := config.NewLogrus()
	viperConfig := config.NewViper(log)
	app := config.NewGin(viperConfig)

	db, err := config.NewGorm(viperConfig, log)
	if err != nil {
		log.Errorf("unable to connect database: %v", err.Error())
	}

	sqlDb, _ := db.DB()
	defer func() {
		if err := sqlDb.Close(); err != nil {
			log.Errorf("failed to close database connection: %v", err)
		} else {
			log.Info("Database connection closed")
		}
	}()


	config.Bootstrap(&config.BootstrapConfig{
		Log:    log,
		App:    app,
		Config: viperConfig,
		DB:     db,
	})

	port := viperConfig.GetInt("app.port")
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: app,
	}

	go func() {
		log.Infof("Server starting on port %d", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("could not start server: %v", err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Info("shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	log.Info("Server exiting")

}
