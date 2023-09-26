package commons

import (
	"gorm.io/gorm"
	"sync"
)

type Config struct {
	APPPort        string   `env:"APP_PORT,required"`
	APPSecret      string   `env:"APP_SECRET"`
	APPOTPValidity string   `env:"APP_OTP_VALIDITY"`
	APPDomain      string   `env:"APP_DOMAIN"`
	APPMode        string   `env:"APP_MODE"`
	DBHost         string   `env:"DB_HOST"`
	DBPort         string   `env:"DB_PORT"`
	DBUsername     string   `env:"DB_USERNAME"`
	DBPassword     string   `env:"DB_PASSWORD"`
	DBName         string   `env:"DB_NAME"`
	DBLogLevel     string   `env:"DB_LOG_LEVEL"`
	DBTimeZone     string   `env:"DB_TimeZone"`
	CORSWhitelist  []string `env:"CORS_WHITELIST"`
	SSL            string   `env:"SSL_MODE"`
}

var (
	Gorm      *gorm.DB
	Once      sync.Once
	EnvConfig Config
)
