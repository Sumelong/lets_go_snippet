package storing

import (
	"os"
	"time"
)

type Config struct {
	Host     string
	Database string
	Port     string
	Driver   string
	User     string
	Password string

	ctxTimeout time.Duration
}

func NewConfigPostgres() *Config {
	return &Config{
		Host:     os.Getenv("POSTGRES_HOST"),
		Database: os.Getenv("POSTGRES_DATABASE"),
		Port:     os.Getenv("POSTGRES_PORT"),
		Driver:   os.Getenv("POSTGRES_DRIVER"),
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
	}
}

func NewConfigSqlite() *Config {
	return &Config{
		Host:     os.Getenv("SQLITE_HOST"),
		Database: os.Getenv("SQLITE_DATABASE"),
		Port:     os.Getenv("SQLITE_PORT"),
		Driver:   os.Getenv("SQLITE_DRIVER"),
		User:     os.Getenv("SQLITE_USER"),
		Password: os.Getenv("SQLITE_PASSWORD"),
	}
}
