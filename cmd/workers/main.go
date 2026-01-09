package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/febry3/gamingin/internal/config"
	"github.com/febry3/gamingin/internal/infra/payment"
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

	db, err := config.NewGorm(viperConfig, log)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	groupBuySessionRepo := pg.NewGroupBuySessionRepositoryPg(db)
	groupBuyTierRepo := pg.NewGroupBuyTierRepositoryPg(db)
	productRepo := pg.NewProductRepositoryPg(db)
	productVariantRepo := pg.NewProductVariantRepositoryPg(db)
	buyerGroupSessionRepo := pg.NewBuyerGroupBuySessionRepositoryPg(db)
	addressRepo := pg.NewAddressRepositoryPg(db)
	buyerGroupMemberRepo := pg.NewBuyerGroupMemberRepositoryPg(db)
	txManager := pg.NewTxManager(db)

	orderRepo := pg.NewOrderRepositoryPg(db)
	paymentRepo := pg.NewPaymentRepositoryPg(db)
	shippingRepo := pg.NewOrderShippingDetailRepositoryPg(db)
	stockRepo := pg.NewProductVariantStockRepositoryPg(db)
	userWalletRepo := pg.NewUserWalletRepositoryPg(db)

	asynqConfig := config.NewAsynqConfig(viperConfig)
	asynqClient := config.NewAsynqClient(asynqConfig, log)
	defer asynqClient.Close()

	midtransCoreClient := config.NewMidtransCoreApiClient(viperConfig)
	serverKey := viperConfig.GetString("midtrans.server_key")
	paymentGateway := payment.NewMidtransGateway(*midtransCoreClient, serverKey, log)

	userWalletUsecase := usecase.NewUserWalletUsecase(userWalletRepo)
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

	orderUsecase := usecase.NewOrderUsecase(
		orderRepo,
		paymentRepo,
		shippingRepo,
		addressRepo,
		productVariantRepo,
		stockRepo,
		buyerGroupSessionRepo,
		paymentGateway,
		txManager,
		asynqClient,
		log,
	)

	groupBuyHandler := worker.NewGroupBuySessionHandler(groupBuyUsecase, asynqClient, email, log)
	orderHandler := worker.NewOrderHandler(orderUsecase, groupBuyUsecase, userWalletUsecase, log)

	srv := config.NewAsynqServer(asynqConfig, log)
	mux := asynq.NewServeMux()

	mux.HandleFunc(tasks.TypeEmailDelivery, tasks.HandleEmailDeliveryTask)
	mux.HandleFunc(tasks.TypeWelcomeEmail, tasks.HandleWelcomeEmailTask)

	mux.HandleFunc(tasks.TypeGroupBuySessionEnd, groupBuyHandler.HandleSessionEnd)
	mux.HandleFunc(tasks.TypeGroupBuySessionEndMail, groupBuyHandler.HandleSessionEndMail)
	mux.HandleFunc(tasks.TypeBuyerGroupBuySessionEnd, groupBuyHandler.HandleBuyerSessionEnd)

	mux.HandleFunc(tasks.TypeOrderExpiration, orderHandler.HandleOrderExpiration)

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
