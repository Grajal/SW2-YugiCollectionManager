// Package database provides PostgreSQL database connection and management functionality
package database

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB is the global database connection instance
var DB *gorm.DB

// DBConnect initializes the PostgreSQL database connection
// Required environment variables:
// - DB_HOST: database host
// - DB_USER: database user
// - DB_PASSWORD: user password
// - DB_NAME: database name
// - DB_PORT: database port
// - DB_SSLMODE: SSL mode
func DBConnect() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	pw := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")
	sslmode := os.Getenv("DB_SSLMODE")

	dsn := "host=" + host + " user=" + user + " password=" + pw + " dbname=" + dbName + " port=" + port + " sslmode=" + sslmode

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Connected to database")
}
