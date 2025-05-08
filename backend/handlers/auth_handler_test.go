package handlers

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Grajal/SW2-YugiCollectionManager/backend/database"
	"github.com/Grajal/SW2-YugiCollectionManager/backend/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode((gin.TestMode))

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database")
	}

	database.DB = db
	if err := database.DB.AutoMigrate(&models.User{}); err != nil {
		panic("failed to migrate database")
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	database.DB.Create(&models.User{
		Username: "testuser",
		Email:    "testuser@example.com",
		Password: string(hashed),
	})

	r := gin.Default()
	r.POST("/api/login", Login)
	return r
}

func TestLogin(t *testing.T) {
	router := setupTestRouter()

	tests := []struct {
		name         string
		payload      string
		wantStatus   int
		wantContains string
	}{
		{
			name:         "Test valid login",
			payload:      `{"username": "testuser", "password": "password"}`,
			wantStatus:   200,
			wantContains: "Login successful",
		},
		{
			name:         "Test invalid login",
			payload:      `{"username": "testuser", "password": "wrongpassword"}`,
			wantStatus:   401,
			wantContains: "invalid password",
		},
		{
			name:         "Test missing username",
			payload:      `{"password": "password"}`,
			wantStatus:   400,
			wantContains: "Invalid input",
		},
		{
			name:         "Test missing password",
			payload:      `{"username": "testuser"}`,
			wantStatus:   400,
			wantContains: "Invalid input",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/api/login", strings.NewReader(tt.payload))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
			assert.Contains(t, w.Body.String(), tt.wantContains)
		})
	}
}
