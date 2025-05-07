package services

import (
	"errors"

	"github.com/Grajal/SW2-YugiCollectionManager/backend/database"
	"github.com/Grajal/SW2-YugiCollectionManager/backend/models"
	"golang.org/x/crypto/bcrypt"
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

// func AuthenticateWithClerk(sessionToken string) (models.User, error) {
// 	var user models.User

// 	client, err := clerk.NewClient(os.Getenv("CLERK_SECRET_KEY"))
// 	if err != nil {
// 		return user, errors.New("failed to initialize Clerk client")
// 	}

// 	session, err := client.Sessions().VerifyToken(sessionToken)
// 	if err != nil {
// 		return user, errors.New("invalid session token")
// 	}

// 	clerkUser, err := client.Users().Read(session.UserID)
// 	if err != nil {
// 		return user, errors.New("failed to get user from Clerk")
// 	}

// 	email := clerkUser.EmailAddresses[0].EmailAddress
// 	username := clerkUser.Username

// 	err = database.DB.Where("email = ?", email).First(&user).Error
// 	if err != nil {
// 		newUser := models.User{
// 			Email:    email,
// 			Username: username,
// 		}

// 		if err := database.DB.Create(&newUser).Error; err != nil {
// 			return user, err
// 		}

// 		return newUser, nil
// 	}

// 	return user, nil
// }
