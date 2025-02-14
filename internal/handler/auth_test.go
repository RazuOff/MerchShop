package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/RazuOff/MerchShop/internal/config"
	"github.com/RazuOff/MerchShop/internal/models"
	"github.com/RazuOff/MerchShop/internal/service"
	mock_service "github.com/RazuOff/MerchShop/internal/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestHandler_auth(t *testing.T) {
	type mockBehavior func(s *mock_service.MockAuth, credentials Credentials)

	testTable := []struct {
		name                 string
		inputBody            string
		inputCredentials     Credentials
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Ok",
			inputBody: `{"username": "username", "password": "qwerty"}`,
			inputCredentials: Credentials{
				Username: "username",
				Password: "qwerty",
			},
			mockBehavior: func(s *mock_service.MockAuth, credentials Credentials) {
				s.EXPECT().RegistrateOrLogin(credentials.Username, credentials.Password).Return(&models.User{Login: credentials.Username, Password: credentials.Password}, nil)
				s.EXPECT().GenerateToken(credentials.Username, gomock.Any(), gomock.Any()).Return("token123", nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"token":"token123"}`,
		},
		{
			name:      "Invalid credentials",
			inputBody: `{"username": "username", "password": "wrongpassword"}`,
			inputCredentials: Credentials{
				Username: "username",
				Password: "wrongpassword",
			},
			mockBehavior: func(s *mock_service.MockAuth, credentials Credentials) {
				s.EXPECT().RegistrateOrLogin(credentials.Username, credentials.Password).Return(nil, &models.ServiceError{TextError: "wrong pass", Code: 401})
			},
			expectedStatusCode:   401,
			expectedResponseBody: `{"errors": "wrong pass"}`,
		},
		{
			name:      "Internal server error",
			inputBody: `{"username": "username", "password": "qwerty"}`,
			inputCredentials: Credentials{
				Username: "username",
				Password: "qwerty",
			},
			mockBehavior: func(s *mock_service.MockAuth, credentials Credentials) {
				s.EXPECT().RegistrateOrLogin(credentials.Username, credentials.Password).Return(nil, &models.ServiceError{TextError: "internal server error", Code: 500})
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"errors": "internal server error"}`,
		},
		{
			name:      "Empty request body",
			inputBody: ``,
			inputCredentials: Credentials{
				Username: "",
				Password: "",
			},
			mockBehavior: func(s *mock_service.MockAuth, credentials Credentials) {
			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"errors": "Incorrect request body"}`,
		},
	}

	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockAuth := mock_service.NewMockAuth(ctrl)
			test.mockBehavior(mockAuth, test.inputCredentials)

			services := &service.Service{Auth: mockAuth}

			h := &Handler{service: services, config: &config.Config{JwtKey: []byte("key")}}

			r := gin.Default()
			r.POST("/auth", h.Auth)

			w := httptest.NewRecorder()

			req := httptest.NewRequest(http.MethodPost, "/auth", bytes.NewBufferString(test.inputBody))
			req.Header.Set("Content-Type", "application/json")

			r.ServeHTTP(w, req)

			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.JSONEq(t, test.expectedResponseBody, w.Body.String())
		})
	}
}
