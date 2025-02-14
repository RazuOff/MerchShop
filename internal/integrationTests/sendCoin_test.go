package integrationtests

import (
	"bytes"
	"encoding/json"
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

func Test_sendCoins(t *testing.T) {
	db, err := database.NewTestDB()

	require.NoError(t, err)

	repo := repository.NewRepository(db)
	service := service.NewService(repo)
	h := handler.NewHandler(service, &config.Config{JwtKey: []byte("test")})

	db.Save(&models.User{Login: "OldUser", Password: "somepass"})

	router := gin.Default()
	router.POST("/auth", h.Auth)
	authRoutes := router.Group("/", h.AuthMiddleware())
	authRoutes.POST("/sendCoin", h.SendCoins)

	userJSON := `{"username": "testuser", "password": "password123"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/auth", bytes.NewBufferString(userJSON))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	var authResp handler.AuthResponce
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &authResp))

	testCases := []struct {
		name                  string
		requestBody           string
		authToken             string
		expectedCode          int
		expectedSenderCoins   int
		expectedRecieverCoins int
		expectedHistory       models.TransactionsHistory
	}{
		{
			name:                  "Success",
			requestBody:           `{"toUser":"OldUser", "amount": 100}`,
			authToken:             authResp.Token,
			expectedCode:          http.StatusOK,
			expectedSenderCoins:   900,
			expectedRecieverCoins: 1100,
			expectedHistory:       models.TransactionsHistory{ID: 1, SenderID: 2, ReceiverID: 1, Coins: 100},
		},
		{
			name:                  "User not found",
			requestBody:           `{"toUser":"UnknownUser", "amount": 100}`,
			authToken:             authResp.Token,
			expectedCode:          http.StatusBadRequest,
			expectedSenderCoins:   1000,
			expectedRecieverCoins: 1000,
		},
		{
			name:                  "Send to self",
			requestBody:           `{"toUser":"testuser", "amount": 100}`,
			authToken:             authResp.Token,
			expectedCode:          http.StatusBadRequest,
			expectedSenderCoins:   1000,
			expectedRecieverCoins: 1000,
		},
		{
			name:                  "Not enough coins",
			requestBody:           `{"toUser":"OldUser", "amount": 10000}`,
			authToken:             authResp.Token,
			expectedCode:          http.StatusBadRequest,
			expectedSenderCoins:   1000,
			expectedRecieverCoins: 1000,
		},
		{
			name:                  "No token",
			requestBody:           `{"toUser":"OldUser", "amount": 100}`,
			authToken:             "",
			expectedCode:          http.StatusUnauthorized,
			expectedSenderCoins:   1000,
			expectedRecieverCoins: 1000,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			db.Exec("UPDATE users SET coins = ? WHERE id = ?", 1000, 1)
			db.Exec("UPDATE users SET coins = ? WHERE id = ?", 1000, 2)
			db.Exec("DELETE FROM transactions_histories; DELETE FROM sqlite_sequence WHERE name='transactions_histories'")

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/sendCoin", bytes.NewBufferString(test.requestBody))
			req.Header.Set("Authorization", test.authToken)
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)

			require.Equal(t, test.expectedCode, w.Code)

			var senderUser models.User
			db.First(&senderUser, 2)
			var recieverUser models.User
			db.First(&recieverUser, 1)

			var history models.TransactionsHistory
			db.First(&history, 1)

			require.Equal(t, test.expectedSenderCoins, senderUser.Coins)
			require.Equal(t, test.expectedRecieverCoins, recieverUser.Coins)
			require.Equal(t, test.expectedHistory, history)
		})
	}
}
