package config

import (
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/spf13/viper"
)

type MidtransConfig struct {
	serverKey string
	clientKey string
}

func NewMidtransConfig(config *viper.Viper) *MidtransConfig {
	return &MidtransConfig{
		serverKey: config.GetString("midtrans.server_key"),
		clientKey: config.GetString("midtrans.client_key"),
	}
}

func NewMidtransCoreApiClient(config *viper.Viper) *coreapi.Client {
	var client coreapi.Client
	serverKey := config.GetString("midtrans.server_key")
	client.New(serverKey, midtrans.Sandbox)
	return &client
}
