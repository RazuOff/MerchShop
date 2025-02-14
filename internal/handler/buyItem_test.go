package handler

import (
	"log"
	"net/http/httptest"
	"testing"

	"github.com/RazuOff/MerchShop/internal/models"
	"github.com/RazuOff/MerchShop/internal/service"
	mock_service "github.com/RazuOff/MerchShop/internal/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestHandler_BuyItem(t *testing.T) {
	type mockBehavior func(s *mock_service.MockCoin, itemName string, userID int)

	testTable := []struct {
		name                 string
		itemName             string
		userID               int
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:     "Success",
			itemName: "item1",
			userID:   1,
			mockBehavior: func(s *mock_service.MockCoin, itemName string, userID int) {
				s.EXPECT().BuyItem(itemName, userID).Return(nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: ``,
		},
		{
			name:     "Item not found",
			itemName: "item1",
			userID:   1,
			mockBehavior: func(s *mock_service.MockCoin, itemName string, userID int) {

				s.EXPECT().BuyItem(itemName, userID).Return(&models.ServiceError{TextError: "item not found", Code: 400})
			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"errors": "item not found"}`,
		},
		{
			name:     "Missing user ID",
			itemName: "item1",
			userID:   0,
			mockBehavior: func(s *mock_service.MockCoin, itemName string, userID int) {
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"errors": "context does not contains user ID"}`,
		},
	}

	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockCoinService := mock_service.NewMockCoin(ctrl)
			test.mockBehavior(mockCoinService, test.itemName, test.userID)

			services := &service.Service{Coin: mockCoinService}
			h := &Handler{service: services}

			r := gin.Default()
			r.GET("/buy/:item", testMiddlware(test.userID), h.BuyItem)

			w := httptest.NewRecorder()

			req := httptest.NewRequest("GET", "/buy/"+test.itemName, nil)
			req.Header.Set("Content-Type", "application/json")

			r.ServeHTTP(w, req)

			log.Printf("ERROR %s", w.Body.String())

			assert.Equal(t, test.expectedStatusCode, w.Code)

			if w.Code == 200 {
				assert.Empty(t, w.Body.String())
			} else {
				assert.JSONEq(t, test.expectedResponseBody, w.Body.String())
			}
		})
	}
}

func testMiddlware(id int) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("id", id)
		c.Next()
	}
}
