package services

import (
	"errors"

	"github.com/Grajal/SW2-YugiCollectionManager/backend/database"
	"github.com/Grajal/SW2-YugiCollectionManager/backend/models"
	"github.com/Grajal/SW2-YugiCollectionManager/backend/repository"
	"golang.org/x/crypto/bcrypt"
)

// AuthService defines methods for authentication and user management.
type AuthService interface {
	Login(username, password string) (*models.User, error)
	Register(username, email, password string) (*models.User, error)
	GetUserByID(id uint) (*models.User, error)
}

type authService struct {
	userRepo repository.UserRepository
}

// NewAuthService creates a new instance of authService with a given user repository.
func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authService{userRepo: userRepo}
}

// Login authenticates a user by their username and password.
// Returns the user if credentials are correct.
func (s *authService) Login(username, password string) (*models.User, error) {
	var user models.User
	if err := database.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, errors.New("user not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid password")
	}

	return &user, nil
}

// Register creates a new user after checking that the username or email does not exist.
// Password is securely hashed before storing.
func (s *authService) Register(username, email, password string) (*models.User, error) {
	var existing models.User
	if err := database.DB.Where("username = ?", username).Or("email = ?", email).First(&existing).Error; err == nil {
		return nil, errors.New("username or email already exists")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	user := models.User{
		Username: username,
		Email:    email,
		Password: string(hashed),
	}

	if err := database.DB.Create(&user).Error; err != nil {
		return nil, errors.New("failed to create user")
	}

	return &user, nil
}

// GetUserByID retrieves a user by ID, preloading collection and decks.
func (s *authService) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	if err := database.DB.Preload("Collection.Card").
		Preload("Collection").
		Preload("Decks").
		First(&user, id).Error; err != nil {
		return nil, errors.New("user not found")
	}
	return &user, nil
}

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
