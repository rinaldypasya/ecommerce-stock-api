package db

import (
	"ecommerce-stock-api/config"
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// ConnectPostgres initializes the GORM PostgreSQL connection
func ConnectPostgres() *gorm.DB {
	dsn := config.AppConfig.DBDSN
	if dsn == "" {
		// fallback default if no environment variable set
		dsn = "host=localhost user=postgres password=postgres dbname=ecommerce port=5432 sslmode=disable TimeZone=Asia/Jakarta"
	}

	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	}

	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Ping DB to validate connection
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get generic DB interface: %v", err)
	}

	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(time.Hour)

	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("Database unreachable: %v", err)
	}

	fmt.Println("✅ Connected to PostgreSQL")
	return db
}
