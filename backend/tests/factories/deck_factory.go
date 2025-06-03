package factories

import (
	"fmt"
	"math/rand"

	"github.com/Grajal/SW2-YugiCollectionManager/backend/models"
	"github.com/Grajal/SW2-YugiCollectionManager/backend/tests/testutils"
)

func DeckFactory(user *models.User, withCards bool) *models.Deck {
	db := testutils.TestDB

	if user == nil {
		userVal := UserFactory()
		user = &userVal
	}

	deck := &models.Deck{
		UserID:      user.ID,
		Name:        fmt.Sprintf("Test Deck %d", rand.Intn(10000)),
		Description: "A test deck created by factory",
	}

	if err := db.Create(deck).Error; err != nil {
		panic("failed to create deck: " + err.Error())
	}

	if withCards {
		for i := 0; i < 5; i++ {
			card := CardFactory(fmt.Sprintf("Card %d", i))
			DeckCardFactory(deck.ID, card.ID, 1, "Main") // Default to Main deck
		}
	}

	return deck
}
