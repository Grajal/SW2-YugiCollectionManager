package services

import (
	"encoding/json"
	"log"

	"github.com/Grajal/SW2-YugiCollectionManager/backend/client"
	"github.com/Grajal/SW2-YugiCollectionManager/backend/models"
)

type CardFactory interface {
	BuildCardFromAPI(apiCard *client.APICard, imageURL string) *models.Card
}

type cardFactory struct{}

func NewCardFactory() CardFactory {
	return &cardFactory{}
}

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

func (f *cardFactory) buildMonster(api *client.APICard) *models.MonsterCard {
	return &models.MonsterCard{
		Atk:       api.Atk,
		Def:       api.Def,
		Level:     api.Level,
		Attribute: api.Attribute,
		Race:      api.Race,
	}
}

func (f *cardFactory) buildSpellTrap(api *client.APICard) *models.SpellTrapCard {
	return &models.SpellTrapCard{
		Type: api.Type,
	}
}

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

func (f *cardFactory) buildPendulum(api *client.APICard) *models.PendulumMonsterCard {
	return &models.PendulumMonsterCard{
		Atk:       api.Atk,
		Def:       api.Def,
		Level:     api.Level,
		Attribute: api.Attribute,
		Scale:     api.Scale,
	}
}
