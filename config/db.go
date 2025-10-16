package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/example/hrms-backend/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	globalDB *gorm.DB
	once     sync.Once
)

func Connect() (*gorm.DB, error) {
	var err error
	once.Do(func() {
		host := os.Getenv("DB_HOST")
		port := os.Getenv("DB_PORT")
		user := os.Getenv("DB_USER")
		pass := os.Getenv("DB_PASSWORD")
		name := os.Getenv("DB_NAME")
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC", host, user, pass, name, port)

		var db *gorm.DB
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			return
		}
		var sqlDB *sql.DB
		sqlDB, err = db.DB()
		if err != nil {
			return
		}
		// Connection pool
		sqlDB.SetMaxOpenConns(50)
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetConnMaxIdleTime(5 * time.Minute)
		sqlDB.SetConnMaxLifetime(30 * time.Minute)

		globalDB = db
	})
	return globalDB, err
}

func GetDB() *gorm.DB {
	if globalDB == nil {
		log.Fatal("database not initialized")
	}
	return globalDB
}

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.Employee{},
		&models.Attendance{},
		&models.Leave{},
	)
}
