package repository

import (
	"github.com/Grajal/SW2-YugiCollectionManager/backend/database"
	"github.com/Grajal/SW2-YugiCollectionManager/backend/models"
	"gorm.io/gorm"
)

// UserRepository defines the interface for user-related database operations.
type UserRepository interface {
	FindByUsername(username string) (*models.User, error)
	FindByUsernameOrEmail(username, email string) (*models.User, error)
	Create(user *models.User) error
	FindByID(id uint) (*models.User, error)
}

type userRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new instance of userRepository using the default DB.
func NewUserRepository() UserRepository {
	return &userRepository{
		db: database.DB,
	}
}

// FindByUsername returns a user by username, or an error if not found.
func (r *userRepository) FindByUsername(username string) (*models.User, error) {
	var user models.User
	if err := database.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByUsernameOrEmail returns a user by username or email, whichever matches first.
func (r *userRepository) FindByUsernameOrEmail(username, email string) (*models.User, error) {
	var user models.User
	if err := database.DB.Where("username = ? OR email = ?", username, email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// Create inserts a new user into the database.
func (r *userRepository) Create(user *models.User) error {
	return database.DB.Create(user).Error
}

// FindByID retrieves a user by ID, including their collection and decks.
func (r *userRepository) FindByID(id uint) (*models.User, error) {
	var user models.User
	if err := database.DB.Preload("Collection.Card").Preload("Collection").Preload("Decks").First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
