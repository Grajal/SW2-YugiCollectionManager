package services

import (
	"errors"

	"github.com/Grajal/SW2-YugiCollectionManager/backend/database"
	"github.com/Grajal/SW2-YugiCollectionManager/backend/models"
	"golang.org/x/crypto/bcrypt"
)

// AuthenticateUser checks if the provided username and password are valid
// and returns the corresponding user if they are.
func AuthenticateUser(username, password string) (models.User, error) {
	var user models.User
	if err := database.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return user, errors.New("user not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return user, errors.New("invalid password")
	}

	return user, nil
}
