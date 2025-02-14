package service

import (
	"fmt"
	"testing"

	"github.com/RazuOff/MerchShop/internal/models"
	"github.com/RazuOff/MerchShop/internal/repository"
	mock_repository "github.com/RazuOff/MerchShop/internal/repository/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestGenerateInfo(t *testing.T) {
	type mockBehavior func(merchRepo *mock_repository.MockMerch, histRepo *mock_repository.MockHistory, userRepo *mock_repository.MockUsers, user *models.User)

	testTable := []struct {
		name          string
		user          *models.User
		mockBehavior  mockBehavior
		expectedInfo  *models.HistoryResponse
		expectedError *models.ServiceError
	}{
		{
			name: "Success",
			user: &models.User{ID: 1, Coins: 100},
			mockBehavior: func(merchRepo *mock_repository.MockMerch, histRepo *mock_repository.MockHistory, userRepo *mock_repository.MockUsers, user *models.User) {
				userRepo.EXPECT().GetUserByID(user.ID).Return(user, nil)
				merchRepo.EXPECT().GetUsersMerch(user).Return([]models.Merch{{Name: "item1"}}, nil)
				merchRepo.EXPECT().GetUserMerchAmount(user, &models.Merch{Name: "item1"}).Return(2, nil)
				histRepo.EXPECT().GetUserHistory(user).Return([]models.TransactionsHistory{
					{SenderID: 1, ReceiverID: 2, Coins: 50},
				}, nil)
				userRepo.EXPECT().GetUsernameByID(2).Return("user2", nil)
			},
			expectedInfo: &models.HistoryResponse{
				Coins: 100,
				CoinHistory: models.CoinHistory{
					Received: []models.ReceivedTransaction{},
					Sent:     []models.SentTransaction{{ToUser: "user2", Amount: 50}},
				},
				Inventory: []models.Item{{Type: "item1", Quantity: 2}},
			},
			expectedError: nil,
		},
		{
			name: "User not found",
			user: &models.User{ID: 1},
			mockBehavior: func(merchRepo *mock_repository.MockMerch, histRepo *mock_repository.MockHistory, userRepo *mock_repository.MockUsers, user *models.User) {
				userRepo.EXPECT().GetUserByID(user.ID).Return(nil, fmt.Errorf("user not found"))
			},
			expectedInfo:  nil,
			expectedError: &models.ServiceError{TextError: "user not found", Code: 500},
		},
	}

	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			histRepo := mock_repository.NewMockHistory(ctrl)
			userRepo := mock_repository.NewMockUsers(ctrl)
			merchRepo := mock_repository.NewMockMerch(ctrl)

			mockService := NewService(&repository.Repository{Users: userRepo, History: histRepo, Merch: merchRepo})
			test.mockBehavior(merchRepo, histRepo, userRepo, test.user)
			info, err := mockService.GenerateInfo(test.user.ID)

			require.Equal(t, test.expectedError, err)
			require.Equal(t, test.expectedInfo, info)

		})
	}
}
