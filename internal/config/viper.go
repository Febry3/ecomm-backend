package config

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func NewViper(log *logrus.Logger) *viper.Viper {
	err := godotenv.Load()
	if err != nil && !os.IsNotExist(err) {
		log.Errorf("Warning: Could not load .env file. Error: %v", err)
	}
	config := viper.New()
	config.AutomaticEnv()
	config.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	return config
}
