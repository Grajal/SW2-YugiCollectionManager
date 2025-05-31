package services

import (
	"encoding/json"
	"log"

	"github.com/Grajal/SW2-YugiCollectionManager/backend/client"
	"github.com/Grajal/SW2-YugiCollectionManager/backend/models"
)

// CardFactory defines the interface for creating Card model instances from API card data.
type CardFactory interface {
	BuildCardFromAPI(apiCard *client.APICard, imageURL string) *models.Card
}

type cardFactory struct{}

func NewCardFactory() CardFactory {
	return &cardFactory{}
}

// BuildCardFromAPI creates a models.Card from the given APICard data and image URL.
// Depending on the type or frameType, it also creates and attaches the appropriate card subtype:
// SpellTrapCard, LinkMonsterCard, PendulumMonsterCard, or MonsterCard.
func (f *cardFactory) BuildCardFromAPI(apiCard *client.APICard, imageURL string) *models.Card {
	card := models.Card{
		CardYGOID: apiCard.ID,
		Name:      apiCard.Name,
		Desc:      apiCard.Desc,
		FrameType: apiCard.FrameType,
		Type:      apiCard.Type,
		ImageURL:  imageURL,
	}

	switch {
	case apiCard.Type == "Spell Card" || apiCard.Type == "Trap Card":
		card.SpellTrapCard = f.buildSpellTrap(apiCard)
	case apiCard.FrameType == "link":
		card.LinkMonsterCard = f.buildLink(apiCard)
	case apiCard.FrameType == "pendulum":
		card.PendulumMonsterCard = f.buildPendulum(apiCard)
	default:
		card.MonsterCard = f.buildMonster(apiCard)
	}

	return &card
}

// buildMonster creates a basic MonsterCard using ATK, DEF, Level, Attribute, and Race from the API.
func (f *cardFactory) buildMonster(api *client.APICard) *models.MonsterCard {
	return &models.MonsterCard{
		Atk:       api.Atk,
		Def:       api.Def,
		Level:     api.Level,
		Attribute: api.Attribute,
		Race:      api.Race,
	}
}

// buildSpellTrap creates a SpellTrapCard using the card type (Spell or Trap) from the API.
func (f *cardFactory) buildSpellTrap(api *client.APICard) *models.SpellTrapCard {
	return &models.SpellTrapCard{
		Type: api.Type,
	}
}

// buildLink creates a LinkMonsterCard with LinkValue, LinkMarkers (as JSON string), ATK, Attribute, and Race.
// Level is always set to 0 for Link monsters.
func (f *cardFactory) buildLink(api *client.APICard) *models.LinkMonsterCard {
	linkMarkersJSON, err := json.Marshal(api.LinkMarkers)
	if err != nil {
		log.Printf("Error marshaling link markers: %v", err)
	}

	return &models.LinkMonsterCard{
		LinkValue:   api.LinkValue,
		LinkMarkers: string(linkMarkersJSON),
		Atk:         api.Atk,
		Level:       0,
		Attribute:   api.Attribute,
		Race:        api.Race,
	}
}

// buildPendulum creates a PendulumMonsterCard with ATK, DEF, Level, Attribute, and Scale from the API.
func (f *cardFactory) buildPendulum(api *client.APICard) *models.PendulumMonsterCard {
	return &models.PendulumMonsterCard{
		Atk:       api.Atk,
		Def:       api.Def,
		Level:     api.Level,
		Attribute: api.Attribute,
		Scale:     api.Scale,
	}
}
