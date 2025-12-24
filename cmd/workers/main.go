package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/febry3/gamingin/internal/config"
	"github.com/febry3/gamingin/internal/worker/tasks"
	"github.com/hibiken/asynq"
)

func main() {
	log := config.NewLogrus()
	viperConfig := config.NewViper(log)

	asynqConfig := config.NewAsynqConfig(viperConfig)
	srv := config.NewAsynqServer(asynqConfig, log)
	mux := asynq.NewServeMux()

	mux.HandleFunc(tasks.TypeEmailDelivery, tasks.HandleEmailDeliveryTask)
	mux.HandleFunc(tasks.TypeWelcomeEmail, tasks.HandleWelcomeEmailTask)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		log.Info("Shutting down worker...")
		srv.Shutdown()
	}()

	log.Info("Starting Asynq worker...")
	if err := srv.Run(mux); err != nil {
		log.Fatalf("Could not start Asynq server: %v", err)
	}
}
