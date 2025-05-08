package services

import (
	"testing"

	"github.com/Grajal/SW2-YugiCollectionManager/backend/database"
	"github.com/Grajal/SW2-YugiCollectionManager/backend/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database")
	}
	database.DB = db
	database.DB.AutoMigrate(&models.User{})
}

func seddTestUsers() {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	database.DB.Create(&models.User{
		Username: "testuser",
		Email:    "testuser@example.com",
		Password: string(hashedPassword),
	})
}

func TestAuthenticateUser(t *testing.T) {
	setupTestDB()
	seddTestUsers()

	type args struct {
		username string
		password string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Test valid user",
			args: args{
				username: "testuser",
				password: "password",
			},
			wantErr: false,
		},
		{
			name: "Test invalid user",
			args: args{
				username: "invaliduser",
				password: "password",
			},
			wantErr: true,
		},
		{
			name: "Test invalid password",
			args: args{
				username: "testuser",
				password: "wrongpassword",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := AuthenticateUser(tt.args.username, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("AuthenticateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
