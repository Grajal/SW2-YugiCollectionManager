package services

import (
	"fmt"

	"github.com/Grajal/SW2-YugiCollectionManager/backend/client"
	"github.com/Grajal/SW2-YugiCollectionManager/backend/database"
	"github.com/Grajal/SW2-YugiCollectionManager/backend/models"
	"github.com/Grajal/SW2-YugiCollectionManager/backend/utils"
)

func GetOrFetchCardByIDOrName(id int, name string) (*models.Card, error) {
	var card models.Card
	var err error

	query := database.DB.Preload("MonsterCard").
		Preload("PendulumMonsterCard").
		Preload("LinkMonsterCard").
		Preload("SpellTrapCard")

	if id > 0 {
		err = query.First(&card, "card_ygo_id = ?", id).Error
	} else {
		err = query.First(&card, "name ILIKE ?", name).Error
	}

	if err == nil {
		return &card, nil
	}

	apiCard, err := client.FetchCardByIDOrName(id, name)
	if err != nil {
		return nil, fmt.Errorf("card not found in external API: %w", err)
	}

	s3URL, err := utils.UploadImage(apiCard.ID, apiCard.ImageURL)
	if err != nil {
		return nil, fmt.Errorf("failed to upload image to S3: %w", err)
	}

	newCard := BuildCardFromAPICard(apiCard, s3URL)

	if err := database.DB.Create(&newCard).Error; err != nil {
		return nil, fmt.Errorf("failed to store new card: %w", err)

	}

	return &newCard, nil
}

func BuildCardFromAPICard(apiCard *client.APICard, imageURL string) models.Card {
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
		card.SpellTrapCard = &models.SpellTrapCard{
			Type: apiCard.Type,
		}
	case apiCard.FrameType == "link":
		card.LinkMonsterCard = &models.LinkMonsterCard{
			LinkValue:   apiCard.LinkValue,
			LinkMarkers: apiCard.LinkMarkers,
		}
		card.MonsterCard = &models.MonsterCard{
			Atk:       apiCard.Atk,
			Def:       0,
			Level:     0,
			Attribute: apiCard.Attribute,
			Race:      apiCard.Race,
		}
	case apiCard.FrameType == "pendulum":
		card.PendulumMonsterCard = &models.PendulumMonsterCard{
			Scale: apiCard.Scale,
		}
		card.MonsterCard = &models.MonsterCard{
			Atk:       apiCard.Atk,
			Def:       apiCard.Def,
			Level:     apiCard.Level,
			Attribute: apiCard.Attribute,
			Race:      apiCard.Race,
		}
	default:

		card.MonsterCard = &models.MonsterCard{
			Atk:       apiCard.Atk,
			Def:       apiCard.Def,
			Level:     apiCard.Level,
			Attribute: apiCard.Attribute,
			Race:      apiCard.Race,
		}

	}

	return card
}
