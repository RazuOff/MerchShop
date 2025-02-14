package service

import (
	"testing"

	"github.com/RazuOff/MerchShop/internal/models"
	mock_repository "github.com/RazuOff/MerchShop/internal/repository/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestServiceAuth_RegistrateOrLogin(t *testing.T) {
	type mockBehavior func(repo *mock_repository.MockUsers, user *models.User)

	testTable := []struct {
		name          string
		user          models.User
		mockBehavior  mockBehavior
		expectedUser  *models.User
		expectedError *models.ServiceError
	}{
		{
			name: "OkLogin",
			user: models.User{ID: 1, Login: "test", Password: "test", Coins: 100},
			mockBehavior: func(repo *mock_repository.MockUsers, user *models.User) {
				repo.EXPECT().GetUserByUsername(user.Login).Return(user, nil)
			},
			expectedUser:  &models.User{ID: 1, Login: "test", Password: "test", Coins: 100},
			expectedError: nil,
		},
		{
			name: "WrongPassword",
			user: models.User{ID: 1, Login: "test", Password: "test", Coins: 100},

			mockBehavior: func(repo *mock_repository.MockUsers, user *models.User) {
				repo.EXPECT().GetUserByUsername(user.Login).Return(&models.User{Login: user.Login, Password: "truepassword"}, nil)
			},
			expectedUser:  nil,
			expectedError: &models.ServiceError{TextError: "incorrect password", Code: 401},
		},
		{
			name: "Registration",
			user: models.User{ID: 1, Login: "test", Password: "test", Coins: 100},

			mockBehavior: func(repo *mock_repository.MockUsers, user *models.User) {
				repo.EXPECT().GetUserByUsername(user.Login).Return(nil, nil)
				repo.EXPECT().SetUser(user.Login, user.Login).Return(user, nil)
			},
			expectedUser:  &models.User{ID: 1, Login: "test", Password: "test", Coins: 100},
			expectedError: nil,
		},
	}

	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mock_repository.NewMockUsers(ctrl)

			test.mockBehavior(mockRepo, &test.user)

			authService := AuthService{repository: mockRepo}
			user, err := authService.RegistrateOrLogin(test.user.Login, test.user.Password)

			assert.Equal(t, test.expectedUser, user)
			if test.expectedError != nil {
				assert.NotNil(t, err)
				assert.Equal(t, test.expectedError.Code, err.Code)
			} else {
				assert.Equal(t, test.expectedError, err)
			}

		})
	}
}
