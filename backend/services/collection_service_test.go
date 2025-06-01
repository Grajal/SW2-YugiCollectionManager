package services

import (
	"testing"

	"github.com/Grajal/SW2-YugiCollectionManager/backend/models"
	"github.com/Grajal/SW2-YugiCollectionManager/backend/repository"
	"github.com/Grajal/SW2-YugiCollectionManager/backend/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_collectionService_GetUserCollection(t *testing.T) {
	db := utils.SetupTestDB(
		&models.User{}, &models.Card{}, &models.UserCard{},
		&models.MonsterCard{}, &models.SpellTrapCard{},
		&models.LinkMonsterCard{}, &models.PendulumMonsterCard{},
	)

	user := &models.User{Username: "testuser", Email: "test@example.com", Password: "securepass"}
	card := &models.Card{CardYGOID: 11111, Name: "Blue-Eyes White Dragon", Type: "Monster"}
	utils.SeedTestData(db, user, card)

	userCard := &models.UserCard{
		UserID:   user.ID,
		CardID:   card.ID,
		Quantity: 2,
	}
	utils.SeedTestData(db, userCard)

	repo := repository.NewCollectionRepositoryWithDB(db)
	service := &collectionService{repo: repo}

	t.Run("returns user collection correctly", func(t *testing.T) {
		got, err := service.GetUserCollection(user.ID)
		require.NoError(t, err)
		require.Len(t, got, 1)

		assert.Equal(t, user.ID, got[0].UserID)
		assert.Equal(t, card.ID, got[0].CardID)
		assert.Equal(t, 2, got[0].Quantity)
		assert.Equal(t, "Blue-Eyes White Dragon", got[0].Card.Name)
	})
}

func Test_collectionService_GetUserCard(t *testing.T) {
	db := utils.SetupTestDB(
		&models.User{}, &models.Card{}, &models.UserCard{},
		&models.MonsterCard{}, &models.SpellTrapCard{},
		&models.LinkMonsterCard{}, &models.PendulumMonsterCard{},
	)

	user := models.User{Username: "testuser", Email: "test@example.com", Password: "pass"}
	card := models.Card{CardYGOID: 12345, Name: "Test Card", Type: "Normal Monster"}

	utils.SeedTestData(db, &user, &card)

	userCard := models.UserCard{UserID: user.ID, CardID: card.ID, Quantity: 2}
	utils.SeedTestData(db, &userCard)

	repo := repository.NewCollectionRepositoryWithDB(db)
	service := &collectionService{repo: repo}

	tests := []struct {
		name string
		args struct {
			userID uint
			cardID uint
		}
		want    *models.UserCard
		wantErr bool
	}{
		{
			name: "returns user card successfully",
			args: struct {
				userID uint
				cardID uint
			}{
				userID: user.ID,
				cardID: card.ID,
			},
			want: &models.UserCard{
				UserID:   user.ID,
				CardID:   card.ID,
				Quantity: 2,
			},
			wantErr: false,
		},
		{
			name: "returns error for non-existent user card",
			args: struct {
				userID uint
				cardID uint
			}{
				userID: user.ID,
				cardID: 999,
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.GetUserCard(tt.args.userID, tt.args.cardID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserCard() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got.UserID != tt.want.UserID || got.CardID != tt.want.CardID || got.Quantity != tt.want.Quantity {
					t.Errorf("GetUserCard() = %+v, want %+v", got, tt.want)
				}
			}
		})
	}
}

func Test_collectionService_AddCardToCollection(t *testing.T) {
	db := utils.SetupTestDB(&models.User{}, &models.Card{}, &models.UserCard{})

	user := &models.User{Username: "testuser", Email: "test@example.com", Password: "pass"}
	card := &models.Card{Name: "Blue-Eyes White Dragon", CardYGOID: 123456}
	utils.SeedTestData(db, user, card)

	repo := repository.NewCollectionRepositoryWithDB(db)
	service := &collectionService{repo: repo}

	tests := []struct {
		name string
		args struct {
			userID   uint
			cardID   uint
			quantity int
		}
		setup   func()
		wantQty int
		wantErr bool
	}{
		{
			name: "adds new card to collection",
			args: struct {
				userID   uint
				cardID   uint
				quantity int
			}{
				userID:   user.ID,
				cardID:   card.ID,
				quantity: 2,
			},
			setup:   func() {},
			wantQty: 2,
			wantErr: false,
		},
		{
			name: "increments quantity of existing card",
			args: struct {
				userID   uint
				cardID   uint
				quantity int
			}{
				userID:   user.ID,
				cardID:   card.ID,
				quantity: 3,
			},
			setup: func() {
				existing := &models.UserCard{UserID: user.ID, CardID: card.ID, Quantity: 2}
				utils.SeedTestData(db, existing)
			},
			wantQty: 5, // 2 + 3
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			utils.CleanupDB(db, &models.UserCard{})

			tt.setup()

			err := service.AddCardToCollection(tt.args.userID, tt.args.cardID, tt.args.quantity)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddCardToCollection() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				var result models.UserCard
				err := db.Where("user_id = ? AND card_id = ?", tt.args.userID, tt.args.cardID).First(&result).Error
				require.NoError(t, err)
				assert.Equal(t, tt.wantQty, result.Quantity)
			}
		})
	}
}

func Test_collectionService_DecreaseCardQuantity(t *testing.T) {
	db := utils.SetupTestDB(&models.User{}, &models.Card{}, &models.UserCard{})

	user := models.User{Email: "test@example.com"}
	card := models.Card{Name: "Decrease Card"}
	utils.SeedTestData(db, &user, &card)

	userCard := models.UserCard{UserID: user.ID, CardID: card.ID, Quantity: 5}
	db.Create(&userCard)

	repo := repository.NewCollectionRepositoryWithDB(db)
	s := &collectionService{repo: repo}

	tests := []struct {
		name string
		args struct {
			userID           uint
			cardID           uint
			quantityToRemove int
		}
		wantErr bool
	}{
		{
			name: "decreases quantity successfully",
			args: struct {
				userID           uint
				cardID           uint
				quantityToRemove int
			}{
				userID:           user.ID,
				cardID:           card.ID,
				quantityToRemove: 3,
			},
			wantErr: false,
		},
		{
			name: "returns error when user card not found",
			args: struct {
				userID           uint
				cardID           uint
				quantityToRemove int
			}{
				userID:           user.ID,
				cardID:           999,
				quantityToRemove: 1,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := s.DecreaseCardQuantity(tt.args.userID, tt.args.cardID, tt.args.quantityToRemove)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecreaseCardQuantity() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
