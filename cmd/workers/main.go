package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/febry3/gamingin/internal/config"
	"github.com/febry3/gamingin/internal/repository/pg"
	"github.com/febry3/gamingin/internal/usecase"
	"github.com/febry3/gamingin/internal/worker"
	"github.com/febry3/gamingin/internal/worker/tasks"
	"github.com/hibiken/asynq"
)

func main() {
	log := config.NewLogrus()
	viperConfig := config.NewViper(log)
	email := config.NewEmail(viperConfig)

	// Initialize database for worker
	db, err := config.NewGorm(viperConfig, log)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// Initialize repositories
	groupBuySessionRepo := pg.NewGroupBuySessionRepositoryPg(db)
	groupBuyTierRepo := pg.NewGroupBuyTierRepositoryPg(db)
	productRepo := pg.NewProductRepositoryPg(db)
	productVariantRepo := pg.NewProductVariantRepositoryPg(db)
	buyerGroupSessionRepo := pg.NewBuyerGroupBuySessionRepositoryPg(db)
	addressRepo := pg.NewAddressRepositoryPg(db)
	buyerGroupMemberRepo := pg.NewBuyerGroupMemberRepositoryPg(db)
	txManager := pg.NewTxManager(db)

	// Initialize usecase (asynqClient is nil since worker doesn't enqueue tasks from usecase)
	groupBuyUsecase := usecase.NewGroupBuyUsecase(
		addressRepo,
		groupBuySessionRepo,
		groupBuyTierRepo,
		productRepo,
		productVariantRepo,
		buyerGroupSessionRepo,
		buyerGroupMemberRepo,
		txManager,
		log,
		nil,
	)

	// Initialize Asynq
	asynqConfig := config.NewAsynqConfig(viperConfig)
	asynqClient := config.NewAsynqClient(asynqConfig, log)
	defer asynqClient.Close()

	// Initialize handler with usecase and asynq client for task chaining
	groupBuyHandler := worker.NewGroupBuySessionHandler(groupBuyUsecase, asynqClient, email, log)

	srv := config.NewAsynqServer(asynqConfig, log)
	mux := asynq.NewServeMux()

	// Register task handlers
	mux.HandleFunc(tasks.TypeEmailDelivery, tasks.HandleEmailDeliveryTask)
	mux.HandleFunc(tasks.TypeWelcomeEmail, tasks.HandleWelcomeEmailTask)

	// Group Buy related
	mux.HandleFunc(tasks.TypeGroupBuySessionEnd, groupBuyHandler.HandleSessionEnd)
	mux.HandleFunc(tasks.TypeGroupBuySessionEndMail, groupBuyHandler.HandleSessionEndMail)

	// Handle graceful shutdown
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
