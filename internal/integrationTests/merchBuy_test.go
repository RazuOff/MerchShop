package integrationtests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/RazuOff/MerchShop/internal/config"
	"github.com/RazuOff/MerchShop/internal/database"
	"github.com/RazuOff/MerchShop/internal/handler"
	"github.com/RazuOff/MerchShop/internal/models"
	"github.com/RazuOff/MerchShop/internal/repository"
	"github.com/RazuOff/MerchShop/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestBuyItemIntegration(t *testing.T) {
	db, err := database.NewTestDB()
	require.NoError(t, err)

	repo := repository.NewRepository(db)
	service := service.NewService(repo)
	h := handler.NewHandler(service, &config.Config{JwtKey: []byte("test")})

	router := gin.Default()
	router.POST("/auth", h.Auth)
	authRoutes := router.Group("/", h.AuthMiddleware())
	authRoutes.GET("/buy/:item", h.BuyItem)

	userJSON := `{"username": "testuser", "password": "password123"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/auth", bytes.NewBufferString(userJSON))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	var authResp handler.AuthResponce
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &authResp))

	testCases := []struct {
		name              string
		authToken         string
		merchName         string
		envSetup          func()
		expectedCode      int
		expectedCoins     int
		expectedUserMerch models.UserMerch
	}{
		{
			name:      "Success",
			authToken: authResp.Token,
			merchName: "test",
			envSetup: func() {
				db.Create(&models.Merch{Name: "test", Price: 100})
			},
			expectedCode:      http.StatusOK,
			expectedCoins:     900,
			expectedUserMerch: models.UserMerch{UserID: 1, MerchID: 1, Amount: 1},
		},
		{
			name:      "Invalid Token",
			authToken: "invalidToken",
			merchName: "test",
			envSetup: func() {
				db.Create(&models.Merch{Name: "test", Price: 100})
			},
			expectedCode:      http.StatusUnauthorized,
			expectedCoins:     1000,
			expectedUserMerch: models.UserMerch{},
		},
		{
			name:      "Insufficient Funds",
			authToken: authResp.Token,
			merchName: "test",
			envSetup: func() {
				db.Create(&models.Merch{Name: "test", Price: 10000})
			},
			expectedCode:      http.StatusBadRequest,
			expectedCoins:     1000,
			expectedUserMerch: models.UserMerch{},
		},
		{
			name:      "Nonexistent Item",
			authToken: authResp.Token,
			merchName: "doesnt exists",
			envSetup: func() {
				db.Create(&models.Merch{Name: "test", Price: 100})
			},
			expectedCode:      http.StatusBadRequest,
			expectedCoins:     1000,
			expectedUserMerch: models.UserMerch{},
		},
	}
	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			db.Exec("DELETE FROM merches; DELETE FROM sqlite_sequence WHERE name='merches'")
			db.Exec("DELETE FROM user_merches; DELETE FROM sqlite_sequence WHERE name='user_merches'")
			db.Exec("UPDATE users SET coins = ? WHERE id = ?", 1000, 1)

			test.envSetup()

			w = httptest.NewRecorder()
			req, _ = http.NewRequest("GET", fmt.Sprintf("/buy/%s", test.merchName), nil)
			req.Header.Set("Authorization", test.authToken)
			router.ServeHTTP(w, req)

			require.Equal(t, test.expectedCode, w.Code)

			var updatedUser models.User
			db.First(&updatedUser, 1)

			var userMerch models.UserMerch
			db.First(&userMerch, 1)

			require.Equal(t, test.expectedCoins, updatedUser.Coins)
			require.Equal(t, test.expectedUserMerch, userMerch)
		})
	}
}
