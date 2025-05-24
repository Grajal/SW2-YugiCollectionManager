package factories

import (
	"github.com/Grajal/SW2-YugiCollectionManager/backend/models"
	"github.com/Grajal/SW2-YugiCollectionManager/backend/tests/testutils"

	"golang.org/x/crypto/bcrypt"
)

// CreateUserFactory generates a user and saves it to the test database
// Email: random, Password: "password"
func UserFactory() models.User {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	user := models.User{
		Username: testutils.GenerateRandomString(10),
		Email:    testutils.GenerateRandomString(10) + "@example.com",
		Password: string(hashedPassword),
	}
	testutils.TestDB.Create(&user)
	testutils.TestDB.First(&user, user.ID)
	return user
}
