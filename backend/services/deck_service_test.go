package services

import (
	"fmt"
	"testing"

	"github.com/Grajal/SW2-YugiCollectionManager/backend/database"
	"github.com/Grajal/SW2-YugiCollectionManager/backend/models"
	"github.com/Grajal/SW2-YugiCollectionManager/backend/repository"
	"github.com/Grajal/SW2-YugiCollectionManager/backend/utils"
	"github.com/stretchr/testify/assert"
)

func TestCreateDeck(t *testing.T) {
	db := utils.SetupTestDB(&models.User{}, &models.Deck{})
	database.DB = db

	user := &models.User{Username: "testuser", Email: "test@example.com", Password: "testpassword"}
	utils.SeedTestData(db, user)

	// Instanciamos el repositorio real con la DB de test
	deckRepo := repository.NewDeckRepository()
	deckService := NewDeckService(deckRepo, nil, nil) // Solo necesitas el repo aquí

	type args struct {
		userID      uint
		name        string
		description string
	}
	tests := []struct {
		name    string
		args    args
		want    *models.Deck
		wantErr bool
	}{
		{
			name: "Create valid deck",
			args: args{
				userID:      user.ID,
				name:        "Test Deck",
				description: "Just testing",
			},
			want: &models.Deck{
				UserID:      user.ID,
				Name:        "Test Deck",
				Description: "Just testing",
			},
			wantErr: false,
		},
		{
			name: "Duplicate deck name",
			args: args{
				userID:      user.ID,
				name:        "Test Deck",
				description: "Another one",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Exceeds maximum decks per user",
			args: args{
				userID:      user.ID,
				name:        "Deck 11",
				description: "Should fail",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Exceeds maximum decks per user" {
				for i := 1; i <= 10; i++ {
					deck := &models.Deck{
						UserID:      user.ID,
						Name:        fmt.Sprintf("Deck %d", i),
						Description: "Test deck",
					}
					utils.SeedTestData(db, deck)
				}
			}

			got, err := deckService.CreateDeck(tt.args.userID, tt.args.name, tt.args.description)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateDeck() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				if got != nil {
					t.Errorf("Expected nil deck on error, got: %v", got)
				}
				return
			}

			if got == nil {
				t.Errorf("Expected deck, got nil")
				return
			}

			if got.UserID != tt.want.UserID || got.Name != tt.want.Name || got.Description != tt.want.Description {
				t.Errorf("CreateDeck() = %+v, want %+v", got, tt.want)
			}

			var dbDeck models.Deck
			if err := db.First(&dbDeck, got.ID).Error; err != nil {
				t.Errorf("Deck not found in DB: %v", err)
			}
		})
	}
}

func Test_deckService_GetCardsByDeck(t *testing.T) {
	db := utils.SetupTestDB(&models.User{}, &models.Deck{}, &models.Card{}, &models.LinkMonsterCard{}, &models.MonsterCard{}, &models.PendulumMonsterCard{}, &models.SpellTrapCard{}, &models.DeckCard{})

	user := models.User{Username: "TestUser", Email: "test@example.com", Password: "securepass"}
	utils.SeedTestData(db, &user)

	card := models.Card{Name: "Dark Magician", CardYGOID: 46986414}
	utils.SeedTestData(db, &card)

	deck := models.Deck{Name: "Test Deck", UserID: user.ID}
	utils.SeedTestData(db, &deck)

	deckCard := models.DeckCard{
		DeckID:   deck.ID,
		CardID:   card.ID,
		Quantity: 3,
		Zone:     "main",
	}
	utils.SeedTestData(db, &deckCard)

	repo := repository.NewDeckRepositoryWithDB(db)

	service := &deckService{
		repo: repo,
	}

	t.Run("returns cards for given deck", func(t *testing.T) {
		got, err := service.GetCardsByDeck(user.ID, deck.ID)

		assert.NoError(t, err)
		assert.Len(t, got, 1)

		gotCard := got[0]
		assert.Equal(t, deck.ID, gotCard.DeckID)
		assert.Equal(t, card.ID, gotCard.CardID)
		assert.Equal(t, 3, gotCard.Quantity)
		assert.Equal(t, "main", gotCard.Zone)

		assert.Equal(t, "Dark Magician", gotCard.Card.Name)
	})
}

func Test_deckService_DeleteDeck(t *testing.T) {
	db := utils.SetupTestDB(&models.Deck{}, &models.User{})

	// Creamos usuario y deck
	user := models.User{Username: "testuser", Email: "test@example.com", Password: "hashedpass"}
	utils.SeedTestData(db, &user)

	deck := models.Deck{Name: "Test Deck", UserID: user.ID}
	utils.SeedTestData(db, &deck)

	// Creamos el repo pasando la db de test
	deckRepo := repository.NewDeckRepositoryWithDB(db)

	// Creamos el servicio
	service := &deckService{
		repo:            deckRepo,
		cardService:     nil,
		deckCardService: nil,
	}

	tests := []struct {
		name string
		args struct {
			deckID uint
			userID uint
		}
		wantErr bool
	}{
		{
			name: "successfully deletes deck",
			args: struct {
				deckID uint
				userID uint
			}{
				deckID: deck.ID,
				userID: user.ID,
			},
			wantErr: false,
		},
		{
			name: "returns error if deck does not exist",
			args: struct {
				deckID uint
				userID uint
			}{
				deckID: 999, // ID inexistente
				userID: user.ID,
			},
			wantErr: true,
		},
		{
			name: "returns error if user is not the owner",
			args: struct {
				deckID uint
				userID uint
			}{
				deckID: deck.ID,
				userID: 999, // ID de usuario incorrecto
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.DeleteDeck(tt.args.deckID, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteDeck() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				var deletedDeck models.Deck
				result := db.First(&deletedDeck, tt.args.deckID)
				if result.Error == nil {
					t.Errorf("Deck was not deleted, found deck with ID %v", deletedDeck.ID)
				}
			}
		})
	}
}

func Test_deckService_GetDecksByUserID(t *testing.T) {
	db := utils.SetupTestDB(&models.User{}, &models.Deck{}, &models.DeckCard{})
	user := models.User{Username: "tester", Email: "tester@example.com", Password: "hashed"}
	utils.SeedTestData(db, &user)

	deck1 := models.Deck{Name: "Test Deck 1", UserID: user.ID}
	deck2 := models.Deck{Name: "Test Deck 2", UserID: user.ID}
	utils.SeedTestData(db, &deck1, &deck2)

	repo := repository.NewDeckRepositoryWithDB(db)
	service := &deckService{
		repo:            repo,
		cardService:     nil, // no se usa en este test
		deckCardService: nil, // no se usa en este test
	}

	tests := []struct {
		name    string
		userID  uint
		want    []models.Deck
		wantErr bool
	}{
		{
			name:    "returns decks for user",
			userID:  user.ID,
			want:    []models.Deck{deck1, deck2},
			wantErr: false,
		},
		{
			name:    "returns empty slice for unknown user",
			userID:  999,
			want:    []models.Deck{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.GetDecksByUserID(tt.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDecksByUserID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(got) != len(tt.want) {
				t.Errorf("Expected %d decks, got %d", len(tt.want), len(got))
				return
			}

			for i := range got {
				if got[i].Name != tt.want[i].Name || got[i].UserID != tt.want[i].UserID {
					t.Errorf("Deck mismatch: got %+v, want %+v", got[i], tt.want[i])
				}
			}
		})
	}
}
