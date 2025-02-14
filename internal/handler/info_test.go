package handler

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/RazuOff/MerchShop/internal/models"
	"github.com/RazuOff/MerchShop/internal/service"
	mock_service "github.com/RazuOff/MerchShop/internal/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestHandler_GetInfo(t *testing.T) {
	type mockBehavior func(s *mock_service.MockInfo, userID int)

	testTable := []struct {
		name                string
		user                models.User
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedResponseOk  models.HistoryResponse
		expectedResponseErr models.ErrorResponce
	}{
		{
			name: "Success",
			user: models.User{ID: 1, Login: "test", Password: "test", Coins: 1000},
			mockBehavior: func(s *mock_service.MockInfo, userID int) {
				s.EXPECT().GenerateInfo(userID).Return(&models.HistoryResponse{Coins: 1000}, nil)
			},
			expectedStatusCode: 200,
			expectedResponseOk: models.HistoryResponse{Coins: 1000},
		},
		{
			name: "User ID not found",
			user: models.User{ID: 0},
			mockBehavior: func(s *mock_service.MockInfo, userID int) {
			},
			expectedStatusCode:  500,
			expectedResponseErr: models.ErrorResponce{Errors: "context does not contains user ID"},
		},
		{
			name: "Service error",
			user: models.User{ID: 1, Login: "test", Password: "test", Coins: 1000},
			mockBehavior: func(s *mock_service.MockInfo, userID int) {
				s.EXPECT().GenerateInfo(userID).Return(nil, &models.ServiceError{TextError: "service error", Code: 400})
			},
			expectedStatusCode:  400,
			expectedResponseErr: models.ErrorResponce{Errors: "service error"},
		},
	}

	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockInfoService := mock_service.NewMockInfo(ctrl)
			test.mockBehavior(mockInfoService, test.user.ID)

			services := &service.Service{Info: mockInfoService}
			h := &Handler{service: services}

			r := gin.Default()
			r.GET("/info", testMiddlware(test.user.ID), h.GetInfo)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/info", nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, test.expectedStatusCode, w.Code)
			if w.Code == 200 {
				rsep, _ := json.Marshal(test.expectedResponseOk)
				assert.Equal(t, rsep, w.Body.Bytes())
			} else {
				assert.NotEmpty(t, test.expectedResponseErr)
			}
		})
	}
}
