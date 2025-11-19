package config

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type logrusWriter struct {
	Logger *logrus.Logger
}

func (l *logrusWriter) Printf(message string, args ...interface{}) {
	l.Logger.Debugf(message, args...)
}

func NewGorm(config *viper.Viper, log *logrus.Logger) (*gorm.DB, error) {
	host := config.GetString("database.host")
	port := config.GetInt("database.port")
	user := config.GetString("database.user")
	password := config.GetString("database.password")
	dbName := config.GetString("database.dbname")
	sslMode := config.GetString("database.sslmode")
	maxOpenConns := config.GetInt("database.max_open_conns")
	maxIdleConns := config.GetInt("database.max_idle_conns")
	connMaxLifeTime := config.GetDuration("database.conn_max_lifetime")
	connMaxIdleTime := config.GetDuration("database.conn_max_idle_time")
	timezone := config.GetString("database.time_zone")

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s&TimeZone=%s",
		user, password, host, port, dbName, sslMode, timezone)

	log.Info(dsn)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.New(&logrusWriter{Logger: log}, logger.Config{
			SlowThreshold:             time.Second * 5,
			Colorful:                  true,
			IgnoreRecordNotFoundError: true,
			ParameterizedQueries:      true,
			LogLevel:                  logger.Info,
		}),
	})

	if err != nil {
		return nil, err
	}

	connection, err := db.DB()
	if err != nil {
		return nil, err
	}

	if err := connection.Ping(); err != nil {
		return nil, err
	}

	connection.SetMaxOpenConns(maxOpenConns)
	connection.SetMaxIdleConns(maxIdleConns)
	connection.SetConnMaxIdleTime(connMaxIdleTime)
	connection.SetConnMaxLifetime(connMaxLifeTime)
	return db, nil
}
