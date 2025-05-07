package services

import (
	"context"
	"errors"

	"github.com/Grajal/SW2-YugiCollectionManager/backend/database"
	"github.com/Grajal/SW2-YugiCollectionManager/backend/models"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/api/idtoken"
)

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

func AuthenticateWithGoogle(idToken string, clientID string) (models.User, error) {
	var user models.User

	payload, err := idtoken.Validate(context.Background(), idToken, clientID)
	if err != nil {
		return user, errors.New("invalid Google ID token")
	}

	email := payload.Claims["email"].(string)
	username := payload.Claims["name"].(string)

	err = database.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		return user, nil
	}

	newUser := models.User{
		Email:    email,
		Username: username,
	}

	if err := database.DB.Create(&newUser).Error; err != nil {
		return user, err
	}

	return newUser, nil
}
