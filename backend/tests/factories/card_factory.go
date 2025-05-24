package factories

import (
	"fmt"
	"math/rand"

	"github.com/Grajal/SW2-YugiCollectionManager/backend/models"
	"github.com/Grajal/SW2-YugiCollectionManager/backend/tests/testutils"
)

// CardFactory generates a Card and saves it to the test database
func CardFactory(cardName string) models.Card {
	card := models.Card{
		CardYGOID: rand.Intn(1000000), // Random unique YGO ID
		Name:      cardName,
		Desc:      "A powerful card used for testing purposes.",
		FrameType: "normal",
		Type:      "Monster",
		ImageURL:  fmt.Sprintf("https://images.ygoprodeck.com/images/cards/%d.jpg", rand.Intn(1000000)),
	}
	testutils.TestDB.Create(&card)
	testutils.TestDB.First(&card, card.ID)
	return card
}
