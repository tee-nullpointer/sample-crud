package db

import (
	"database/sql"
	"fmt"
	"sample-crud/internal/config"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	db    *gorm.DB
	sqlDB *sql.DB
)

func Init(config config.DatabaseConfig) *gorm.DB {
	zap.L().Info("Initializing database connection")
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		config.Host, config.User, config.Password, config.Name, config.Port, config.SSLMode)
	var dbErr error
	db, dbErr = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if dbErr != nil {
		panic(fmt.Sprintf("Fail to initialize database connection: %v", dbErr))
	}
	var sqlDBErr error
	sqlDB, sqlDBErr = db.DB()
	if sqlDBErr != nil {
		panic(fmt.Sprintf("Fail to initialize database connection: %v", sqlDBErr))
	}
	zap.L().Info("Database connection established")
	zap.L().Info("Setting database connection pool parameters",
		zap.Int("Max Connection", config.MaxConnection),
		zap.Int("Max Idle", config.MaxIdle),
		zap.Duration("Max Lifetime", config.MaxLifetime),
		zap.Duration("Max Idle Time", config.MaxIdleTime))
	sqlDB.SetMaxOpenConns(config.MaxConnection)
	sqlDB.SetMaxIdleConns(config.MaxIdle)
	sqlDB.SetConnMaxLifetime(config.MaxLifetime)
	sqlDB.SetConnMaxIdleTime(config.MaxIdleTime)
	return db
}

func ShutDown() {
	zap.L().Info("Shutting down database connection")
	if sqlDB != nil {
		if err := sqlDB.Close(); err != nil {
			zap.L().Error("Fail to close database connection", zap.Error(err))
			return
		}
	}
	zap.L().Info("Database connection closed")
}
