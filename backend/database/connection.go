package database

import (
	"log"
	"os"

	"github.com/Grajal/SW2-YugiCollectionManager/backend/models"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

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

// CheckIfCardExists checks if a card with a given ID exists in the database
func CheckIfCardExists(cardID uint) (bool, error) {
	var card models.Card
	result := DB.First(&card, "id = ?", cardID)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return false, result.Error
	}
	return result.RowsAffected > 0, nil
}
