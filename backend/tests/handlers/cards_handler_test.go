package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/Grajal/SW2-YugiCollectionManager/backend/database"
	"github.com/Grajal/SW2-YugiCollectionManager/backend/models"
	"github.com/Grajal/SW2-YugiCollectionManager/backend/routes"
	"github.com/stretchr/testify/assert"
)

func SetupTestDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// AutoMigrate your models
	err = db.AutoMigrate(
		&models.User{},
		&models.Card{},
		&models.SpellTrapCard{},
		&models.MonsterCard{},
		&models.LinkMonsterCard{},
		&models.PendulumMonsterCard{},
		&models.UserCard{},
		&models.Deck{},
	)
	return db, err
}

func TestGetOrFetchCard(t *testing.T) {
	database.DBConnect()
	if err := database.DB.AutoMigrate(models.User{}, models.Card{}, models.SpellTrapCard{}, models.MonsterCard{}, models.LinkMonsterCard{}, models.PendulumMonsterCard{}, models.UserCard{}, models.Deck{}); err != nil {
		panic("Failed to migrate database: " + err.Error())
	}

	router := routes.SetupRouter()
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/cards/Dark%20Magician", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var card models.Card
	err := json.Unmarshal(w.Body.Bytes(), &card)
	assert.NoError(t, err)

	assert.Equal(t, "Dark Magician", card.Name)
	assert.NotEmpty(t, card.ImageURL) // Optional: Validate expected fields
	assert.NotZero(t, card.CardYGOID) // Ensure YGOID is returned
}

func TestGetOrFetchCard_NotFound(t *testing.T) {
	database.DBConnect()
	if err := database.DB.AutoMigrate(models.User{}, models.Card{}, models.SpellTrapCard{}, models.MonsterCard{}, models.LinkMonsterCard{}, models.PendulumMonsterCard{}, models.UserCard{}, models.Deck{}); err != nil {
		panic("Failed to migrate database: " + err.Error())
	}

	router := routes.SetupRouter()
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/cards/NonExistentCardName", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response["error"], "not found")
}
