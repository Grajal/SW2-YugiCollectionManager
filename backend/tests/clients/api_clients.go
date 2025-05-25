package api_clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Grajal/SW2-YugiCollectionManager/backend/models"
	"github.com/Grajal/SW2-YugiCollectionManager/backend/routes"
	"github.com/Grajal/SW2-YugiCollectionManager/backend/tests/factories"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

// SetupTestRouter sets up the Gin router for testing purposes
func SetupTestRouter() *gin.Engine {
	router := routes.SetupRouter()
	return router
}

// GenerateTestJWT generates a JWT token for testing
func GenerateTestJWT(userID uint) (string, error) {
	// Define the token claims
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	secret := "JWT_SECRET"
	fmt.Println("Signing JWT with SECRET_KEY:", secret) // DEBUG
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// AssertResponse checks the response status and body
func AssertResponse(t *testing.T, response *httptest.ResponseRecorder, expectedStatus int, expectedBody string) {
	assert.Equal(t, expectedStatus, response.Code)
	if expectedBody != "" {
		assert.JSONEq(t, expectedBody, response.Body.String())
	}
}

// TestClient represents a reusable client for API testing
type TestClient struct {
	Router *gin.Engine
	Token  string
	User   models.User
}

// PerformRequest is a helper function to perform HTTP requests in tests
func (client *TestClient) PerformRequest(method, path string, body interface{}, headers map[string]string) *httptest.ResponseRecorder {
	var reqBody *bytes.Reader
	if body != nil {
		jsonBody, _ := json.Marshal(body)
		reqBody = bytes.NewReader(jsonBody)
	} else {
		reqBody = bytes.NewReader([]byte{})
	}

	req, _ := http.NewRequest(method, path, reqBody)

	// Set default headers
	req.Header.Set("Content-Type", "application/json")

	// Add the token to the Authorization header if available
	if client.Token != "" {
		req.Header.Set("Authorization", "Bearer "+client.Token)
	}

	// Set additional headers if provided
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	w := httptest.NewRecorder()
	client.Router.ServeHTTP(w, req)
	return w
}

// NewTestClientWithAuth initializes a new TestClient with authentication and creates a user in DB
func NewTestClient(authenticate bool) *TestClient {
	gin.SetMode(gin.ReleaseMode) // disable debug logs for testing
	router := SetupTestRouter()
	client := &TestClient{Router: router}

	if authenticate {
		user := factories.UserFactory()
		token, err := GenerateTestJWT(user.ID)
		if err != nil {
			log.Fatalf("Failed to generate JWT token: %v", err)
		}
		client.Token = token
		client.User = user
	}
	return client
}
