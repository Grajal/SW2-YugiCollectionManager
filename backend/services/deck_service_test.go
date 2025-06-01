package services

import (
	"fmt"
	"testing"

	"github.com/Grajal/SW2-YugiCollectionManager/backend/database"
	"github.com/Grajal/SW2-YugiCollectionManager/backend/models"
	"github.com/Grajal/SW2-YugiCollectionManager/backend/repository"
	"github.com/Grajal/SW2-YugiCollectionManager/backend/utils"
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
