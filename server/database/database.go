package database

import (
	"C2-D2/server/models"
	"fmt"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// Initialize DB connection with GORM
func Initialize(DB_USER string, DB_PASSWORD string, DB_NAME string, DB_PORT string) {
	stringconn := fmt.Sprintf("user=%s password=%s dbname=%s port=%s sslmode=disable", DB_USER, DB_PASSWORD, DB_NAME, DB_PORT)
	var err error
	DB, err = gorm.Open(postgres.New(postgres.Config{
		DSN: stringconn,
	}), &gorm.Config{})
	if err != nil {
		logrus.Error(err)
		panic("Cannot connect to PostgreSQL!")
	}
	logrus.Info("Connected to PostgreSQL...")
}

// Ensure tables from the model library are created on the DB
func Migrate() {
	err := DB.AutoMigrate(&models.Agent{}, &models.Task{})
	if err != nil {
		logrus.Error("Database migration failed!")
	}
	logrus.Info("Database migration completed...")
}
