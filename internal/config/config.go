package config

import "time"

type Config struct {
	Database DBConfig
	App      AppConfig
}

type DBConfig struct {
	Host            string
	Port            int
	User            string
	Password        string
	DBName          string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

type AppConfig struct {
	Port int
	Mode string
}
