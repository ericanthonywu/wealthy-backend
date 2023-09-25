package databases

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DBConnection struct {
	db *gorm.DB
}

func NewDBConnection() (*gorm.DB, error) {
	var sqlDB *sql.DB

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")
	SSLMode := os.Getenv("SSL_MODE")
	dbTimeZone := os.Getenv("DB_TimeZone")

	dsn := fmt.Sprintf(`host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s`,
		dbHost, dbUser, dbPassword, dbName, dbPort, SSLMode, dbTimeZone,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		PrepareStmt:            true,
		SkipDefaultTransaction: true,
	})

	if err != nil {
		logrus.Error(err.Error())
		panic(err.Error())
		return db, err
	}

	sqlDB, err = db.DB()
	if err != nil {
		logrus.Warn(err.Error())
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(250)
	sqlDB.SetConnMaxLifetime(time.Hour)

	if err != nil {
		logrus.Error(err.Error())
		panic(err.Error())
		return db, err
	}

	logrus.Info("database connected")

	if err := db.Use(otelgorm.NewPlugin()); err != nil {
		panic(err.Error())
		return db, err
	}

	return db, nil
}
