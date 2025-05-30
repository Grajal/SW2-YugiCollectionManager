package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/Grajal/SW2-YugiCollectionManager/backend/client"
	"github.com/Grajal/SW2-YugiCollectionManager/backend/database"
	"github.com/Grajal/SW2-YugiCollectionManager/backend/models"
	"github.com/Grajal/SW2-YugiCollectionManager/backend/utils"
	"gorm.io/gorm"
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

func GetCardByID(id uint) (*models.Card, error) {
	var card models.Card
	err := database.DB.Preload("MonsterCard").Preload("SpellTrapCard").Preload("LinkMonsterCard").Preload("PendulumMonsterCard").First(&card, "id = ?", id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("card with ID %d not foun in internal database", id)
	} else if err != nil {
		return nil, fmt.Errorf("failed to retrieve card from database: %w", err)
	}

	return &card, nil
}

// BuildCardFromAPICard constructs a models.Card object from a given APICard retrieved from the external API.
// It maps the general card fields and dynamically assigns the appropriate substructure based on the card type.
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
		linkMarkersJSON, err := json.Marshal(apiCard.LinkMarkers)
		if err != nil {
			log.Printf("Error marshaling link markers: %v", err)
			break
		}

		card.LinkMonsterCard = &models.LinkMonsterCard{
			LinkValue:   apiCard.LinkValue,
			LinkMarkers: string(linkMarkersJSON),
			Atk:         apiCard.Atk,
			Level:       0,
			Attribute:   apiCard.Attribute,
			Race:        apiCard.Race,
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

// GetCards retrieves a paginated list of cards from the database, including all their subtype associations.
// Parameters:
//   - limit: the maximum number of cards to retrieve.
//   - offset: the number of cards to skip before starting to return results (used for pagination).
func GetCards(limit int, offset int) ([]models.Card, error) {
	var cards []models.Card

	err := database.DB.Preload("MonsterCard").Preload("SpellTrapCard").Preload("LinkMonsterCard").Preload("PendulumMonsterCard").Limit(limit).Offset(offset).Find(&cards).Error

	if err != nil {
		return nil, fmt.Errorf("failed to retrieve cards: %w", err)
	}

	return cards, nil
}

// CountAllCards returns the total number of cards stored in the database.
func CountAllCards() (int64, error) {
	var count int64
	err := database.DB.Model(&models.Card{}).Count(&count).Error
	return count, err
}

// EnsureMinimumCards checks if there are at least 'min' cards stored in the database.
// If not, it fetches the missing number of random cards from the external YGOProDeck API
// and stores them in the database (including uploading their images to S3).
func EnsureMinimumCards(min int) error {
	var count int64
	err := database.DB.Model(&models.Card{}).Count(&count).Error
	if err != nil {
		return fmt.Errorf("failed to count cards: %w", err)
	}

	if count >= int64(min) {
		return nil
	}

	needed := int(int64(min) - count)
	apiCards, err := client.FetchRandomCards(needed)
	if err != nil {
		return fmt.Errorf("failed to fetch new cards: %w", err)
	}

	for _, apiCard := range apiCards {
		if apiCard.ID == 0 {
			fmt.Println("Invalid card received from API:", apiCard)
			continue
		}

		var existing models.Card
		err := database.DB.Where("card_ygo_id = ?", apiCard.ID).First(&existing).Error
		if err == nil {
			continue
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("error checking card: %w", err)
		}

		imageURL := ""
		if len(apiCard.CardImages) > 0 {
			imageURL = apiCard.CardImages[0].ImageURL
		}

		s3URL, err := utils.UploadCardImageToS3(apiCard.ID, imageURL)
		if err != nil {
			continue
		}

		newCard := BuildCardFromAPICard(&apiCard, s3URL)
		if err := database.DB.Create(&newCard).Error; err != nil {
			return fmt.Errorf("failed to store new card: %w", err)
		}
	}

	return nil
}

// GetCardsByFilters returns cards filtered by name, type, and frameType from the database.
// If no cards match, it attempts to fetch new ones from the external API and save them.
func GetFilteredCards(name, cardType, frameType string, limit, offset int) ([]models.Card, error) {
	db := database.DB.Model(&models.Card{}).
		Preload("MonsterCard").
		Preload("SpellTrapCard").
		Preload("LinkMonsterCard").
		Preload("PendulumMonsterCard")

	if name != "" {
		db = db.Where("name ILIKE ?", "%"+name+"%")
	}
	if cardType != "" {
		db = db.Where("type = ?", cardType)
	}
	if frameType != "" {
		db = db.Where("frame_type = ?", frameType)
	}

	var cards []models.Card
	err := db.Limit(limit).Offset(offset).Find(&cards).Error
	return cards, err
}

// CountFilteredCards returns the number of cards in the database that match the provided filters.
// It supports filtering by name (case-insensitive, partial match), card type, and frameType.
func CountFilteredCards(name, cardType, frameType string) (int64, error) {
	db := database.DB.Model(&models.Card{})

	if name != "" {
		db = db.Where("name ILIKE ?", "%"+name+"%")
	}
	if cardType != "" {
		db = db.Where("type = ?", cardType)
	}
	if frameType != "" {
		db = db.Where("frame_type = ?", frameType)
	}

	var count int64
	err := db.Count(&count).Error
	return count, err
}

// FetchAndStoreCardsByName fetches cards from the external YGOProDeck API based on a name filter.
// It stores only the cards that don't already exist in the database, uploading their images to S3.
// It returns a list of the stored (or previously existing) cards.
func FetchAndStoreCardsByName(name string) ([]models.Card, error) {
	apiCards, err := client.FetchCardsByName(name)
	if err != nil {
		return nil, err
	}
	if len(apiCards) == 0 {
		return nil, nil
	}

	var stored []models.Card
	for _, apiCard := range apiCards {
		if apiCard.ID == 0 || apiCard.Name == "" {
			continue
		}

		var existing models.Card
		err := database.DB.Where("card_ygo_id = ?", apiCard.ID).First(&existing).Error
		if err == nil {
			stored = append(stored, existing)
			continue
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			continue
		}

		s3URL, err := utils.UploadCardImageToS3(apiCard.ID, apiCard.CardImages[0].ImageURL)
		if err != nil {
			continue
		}

		card := BuildCardFromAPICard(&apiCard, s3URL)
		if err := database.DB.Create(&card).Error; err == nil {
			stored = append(stored, card)
		}
	}

	return stored, nil
}
