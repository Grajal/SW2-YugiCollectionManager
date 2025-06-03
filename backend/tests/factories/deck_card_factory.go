package factories

import (
	"github.com/Grajal/SW2-YugiCollectionManager/backend/models"
	"github.com/Grajal/SW2-YugiCollectionManager/backend/tests/testutils"
)

func DeckCardFactory(deckID uint, cardID uint, quantity int, zone string) *models.DeckCard {
	db := testutils.TestDB

	deckCard := &models.DeckCard{
		DeckID:   deckID,
		CardID:   cardID,
		Quantity: quantity,
		Zone:     zone,
	}

	if err := db.Create(deckCard).Error; err != nil {
		panic("failed to create deck card: " + err.Error())
	}

	return deckCard
}
