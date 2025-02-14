package handler

import (
	"bytes"
	"net/http/httptest"
	"testing"

	"github.com/RazuOff/MerchShop/internal/models"
	"github.com/RazuOff/MerchShop/internal/service"
	mock_service "github.com/RazuOff/MerchShop/internal/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestHandler_SendCoins(t *testing.T) {
	type mockBehavior func(s *mock_service.MockCoin, userID int, toUserLogin string, amount int)

	testTable := []struct {
		name                 string
		userID               int
		toUserLogin          string
		amount               int
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
		inputBody            string
	}{
		{
			name:        "Success",
			userID:      1,
			toUserLogin: "user2",
			amount:      100,
			inputBody:   `{"toUser": "user2", "amount": 100}`,
			mockBehavior: func(s *mock_service.MockCoin, userID int, toUserLogin string, amount int) {
				s.EXPECT().SendCoins(userID, toUserLogin, amount).Return(nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: "{}",
		},
		{
			name:      "Missing user ID",
			userID:    0,
			inputBody: `{"toUser": "user2", "amount": 100}`,
			mockBehavior: func(s *mock_service.MockCoin, userID int, toUserLogin string, amount int) {
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"errors": "context does not contains user ID"}`,
		},
		{
			name:      "Invalid request body",
			userID:    1,
			inputBody: `{"toUser": "user2", "amount": "not a number"}`,
			mockBehavior: func(s *mock_service.MockCoin, userID int, toUserLogin string, amount int) {
			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"errors": "incorrect request body"}`,
		},
		{
			name:        "Service error",
			userID:      1,
			toUserLogin: "user2",
			amount:      100,
			inputBody:   `{"toUser": "user2", "amount": 100}`,
			mockBehavior: func(s *mock_service.MockCoin, userID int, toUserLogin string, amount int) {
				s.EXPECT().SendCoins(userID, toUserLogin, amount).Return(&models.ServiceError{TextError: "insufficient funds", Code: 400})
			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"errors": "insufficient funds"}`,
		},
	}

	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockCoinService := mock_service.NewMockCoin(ctrl)
			test.mockBehavior(mockCoinService, test.userID, test.toUserLogin, test.amount)

			services := &service.Service{Coin: mockCoinService}
			h := &Handler{service: services}

			r := gin.Default()
			r.POST("/send", testMiddlware(test.userID), h.SendCoins)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/send", bytes.NewBufferString(test.inputBody))
			req.Header.Set("Content-Type", "application/json")

			r.ServeHTTP(w, req)

			assert.Equal(t, test.expectedStatusCode, w.Code)

			if w.Code == 200 {
				assert.Empty(t, w.Body)
			} else {
				assert.JSONEq(t, test.expectedResponseBody, w.Body.String())
			}
		})
	}
}
