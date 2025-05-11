// Package database provides PostgreSQL database connection and management functionality
package database

import (
	"log"
	"os"
	"time"

	"github.com/Grajal/SW2-YugiCollectionManager/backend/models"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB is the global database connection instance
var DB *gorm.DB

// DBConnect initializes the PostgreSQL database connection
// Required environment variables:
// - PGHOST: database host
// - PGUSER: database user
// - PGPASSWORD: user password
// - PGNAME: database name
// - PGPORT: database port
// - PGSSLMODE: SSL mode
func DBConnect() {
	if _, err := os.Stat(".env"); err == nil {
		if err := godotenv.Load(); err != nil {
			log.Fatalf("Error loading .env file")
		} else {
			log.Println("Loaded .env file")
		}
	} else {
		log.Println("No .env file found; assuming database is already connected")
	}

	host := os.Getenv("PGHOST")
	user := os.Getenv("PGUSER")
	pw := os.Getenv("PGPASSWORD")
	dbName := os.Getenv("PGDATABASE")
	port := os.Getenv("PGPORT")
	sslmode := os.Getenv("PGSSLMODE")

	dsn := "host=" + host + " user=" + user + " password=" + pw + " dbname=" + dbName + " port=" + port + " sslmode=" + sslmode

	var err error
	maxRetries := 10
	for i := 0; i < maxRetries; i++ {
		DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			log.Println("Connected to database")
			return
		}

		log.Printf("Failed to connect to database (attempt %d/%d): %v", i+1, maxRetries, err)
		time.Sleep(2 * time.Second)
	}

	log.Fatalf("Failed to connect to database after %d attempts: %v", maxRetries, err)
}

// CheckIfCardExists checks if a card with a given ID exists in the database
func CheckIfCardExists(cardID uint) (bool, error) {
	var card models.Card
	result := DB.First(&card, "id = ?", cardID)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return false, result.Error
	}
	return result.RowsAffected > 0, nil
}
