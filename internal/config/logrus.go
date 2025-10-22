package config

import "github.com/sirupsen/logrus"

func NewLogrus() *logrus.Logger {
	log := logrus.New()
	log.SetLevel(logrus.DebugLevel)
	log.SetFormatter(&logrus.JSONFormatter{})
	return log
}
