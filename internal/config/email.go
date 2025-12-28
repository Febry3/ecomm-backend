package config

import (
	"github.com/spf13/viper"
	"gopkg.in/mail.v2"
)

func NewEmail(viper *viper.Viper) *mail.Dialer {
	mailHost := viper.GetString("mail.host")
	mailPort := viper.GetInt("mail.port")
	mailUsername := viper.GetString("mail.username")
	mailPassword := viper.GetString("mail.password")

	return mail.NewDialer(mailHost, mailPort, mailUsername, mailPassword)
}
