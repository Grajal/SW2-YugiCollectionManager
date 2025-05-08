package services

import (
	"github.com/Grajal/SW2-YugiCollectionManager/backend/database"
	"github.com/Grajal/SW2-YugiCollectionManager/backend/models"
)

func FindOrCreateUser(clerkID, username string) (models.User, error) {
	var user models.User
	err := database.DB.Where("clerk_id = ?", clerkID).First(&user).Error
	if err == nil {
		return user, nil
	}

	user = models.User{
		ClerkID:  clerkID,
		Username: username,
	}

	if err := database.DB.Create(&user).Error; err != nil {
		return user, err
	}

	return user, nil
}
