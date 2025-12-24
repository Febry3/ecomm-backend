package config

import (
	"context"

	"github.com/hibiken/asynq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type AsynqConfig struct {
	RedisAddr     string
	RedisPassword string
	RedisDB       int
	Concurrency   int
}

func NewAsynqConfig(config *viper.Viper) *AsynqConfig {
	return &AsynqConfig{
		RedisAddr:     config.GetString("redis.addr"),
		RedisPassword: config.GetString("redis.password"),
		RedisDB:       config.GetInt("redis.db"),
		Concurrency:   config.GetInt("asynq.concurrency"),
	}
}

func (c *AsynqConfig) GetRedisClientOpt() asynq.RedisClientOpt {
	return asynq.RedisClientOpt{
		Addr:     c.RedisAddr,
		Password: c.RedisPassword,
		DB:       c.RedisDB,
	}
}

func NewAsynqClient(config *AsynqConfig, log *logrus.Logger) *asynq.Client {
	client := asynq.NewClient(config.GetRedisClientOpt())
	log.Info("Asynq client initialized")
	return client
}

func NewAsynqServer(config *AsynqConfig, log *logrus.Logger) *asynq.Server {
	srv := asynq.NewServer(
		config.GetRedisClientOpt(),
		asynq.Config{
			Concurrency: config.Concurrency,
			Queues: map[string]int{
				"critical": 6, // High priority tasks
				"default":  3, // Normal priority tasks
				"low":      1, // Low priority tasks
			},
			ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
				log.Errorf("Error processing task %s: %v", task.Type(), err)
			}),
			Logger: log,
		},
	)
	log.Info("Asynq server initialized")
	return srv
}

func NewAsynqInspector(config *AsynqConfig) *asynq.Inspector {
	return asynq.NewInspector(config.GetRedisClientOpt())
}
