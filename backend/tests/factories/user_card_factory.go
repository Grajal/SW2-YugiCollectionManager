package factories

import (
	"log"

	"github.com/Grajal/SW2-YugiCollectionManager/backend/models"
	"github.com/Grajal/SW2-YugiCollectionManager/backend/tests/testutils"
)

// UserCardFactory creates and returns a UserCard with the given user ID, card ID, and quantity.
// If a user or card is not provided, you can chain with UserFactory or CardFactory beforehand.
func UserCardFactory(userID, cardID uint, quantity int) models.UserCard {

	db := testutils.TestDB

	userCard := models.UserCard{
		UserID:   userID,
		CardID:   cardID,
		Quantity: quantity,
	}

	err := db.Create(&userCard).Error
	if err != nil {
		log.Panicf("failed to create user card: %v", err)
	}

	return userCard
}
