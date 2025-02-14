package service

import (
	"errors"
	"fmt"
	"testing"

	"github.com/RazuOff/MerchShop/internal/models"
	"github.com/RazuOff/MerchShop/internal/repository"
	mock_repository "github.com/RazuOff/MerchShop/internal/repository/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

type userMatcher struct {
	ID       int
	Login    string
	Password string
	Coins    int
}

func (m userMatcher) Matches(x interface{}) bool {
	user, ok := x.(*models.User)
	if !ok || user == nil {
		return false
	}
	return user.ID == m.ID && user.Coins == m.Coins &&
		user.Login == m.Login && user.Password == m.Password
}

func (m userMatcher) String() string {
	return fmt.Sprintf("User with ID=%d, Coins=%d", m.ID, m.Coins)
}

func TestCoinService_SendCoin(t *testing.T) {
	type mockBehavior func(historyRep *mock_repository.MockHistory, userRep *mock_repository.MockUsers, fromUser *models.User, toUser *models.User, amount int)

	testTable := []struct {
		name          string
		fromUser      models.User
		toUser        models.User
		amount        int
		mockBehavior  mockBehavior
		expectedError *models.ServiceError
	}{
		{
			name:     "Error getting fromUser",
			fromUser: models.User{ID: 1, Login: "from", Password: "from", Coins: 100},
			toUser:   models.User{ID: 2, Login: "to", Password: "to", Coins: 100},
			amount:   50,
			mockBehavior: func(historyRep *mock_repository.MockHistory, userRep *mock_repository.MockUsers, fromUser *models.User, toUser *models.User, amount int) {
				userRep.EXPECT().GetUserByID(fromUser.ID).Return(nil, fmt.Errorf("error fetching user"))
			},
			expectedError: &models.ServiceError{TextError: "error fetching user", Code: 500},
		},
		{
			name:     "Error getting toUser",
			fromUser: models.User{ID: 1, Login: "from", Password: "from", Coins: 100},
			toUser:   models.User{ID: 2, Login: "to", Password: "to", Coins: 100},
			amount:   50,
			mockBehavior: func(historyRep *mock_repository.MockHistory, userRep *mock_repository.MockUsers, fromUser *models.User, toUser *models.User, amount int) {
				userRep.EXPECT().GetUserByID(fromUser.ID).Return(fromUser, nil)
				userRep.EXPECT().GetUserByUsername(toUser.Login).Return(nil, fmt.Errorf("error fetching user"))
			},
			expectedError: &models.ServiceError{TextError: "error fetching user", Code: 500},
		},
		{
			name:     "ToUser does not exist",
			fromUser: models.User{ID: 1, Login: "from", Password: "from", Coins: 100},
			toUser:   models.User{ID: 2, Login: "to", Password: "to", Coins: 100},
			amount:   50,
			mockBehavior: func(historyRep *mock_repository.MockHistory, userRep *mock_repository.MockUsers, fromUser *models.User, toUser *models.User, amount int) {
				userRep.EXPECT().GetUserByID(fromUser.ID).Return(fromUser, nil)
				userRep.EXPECT().GetUserByUsername(toUser.Login).Return(nil, nil)
			},
			expectedError: &models.ServiceError{TextError: "ToUserLogin does not exists", Code: 400},
		},
		{
			name:     "Send coins to self",
			fromUser: models.User{ID: 1, Login: "from", Password: "from", Coins: 100},
			toUser:   models.User{ID: 1, Login: "from", Password: "from", Coins: 100},
			amount:   50,
			mockBehavior: func(historyRep *mock_repository.MockHistory, userRep *mock_repository.MockUsers, fromUser *models.User, toUser *models.User, amount int) {
				userRep.EXPECT().GetUserByID(fromUser.ID).Return(fromUser, nil)
				userRep.EXPECT().GetUserByUsername(toUser.Login).Return(toUser, nil)
			},
			expectedError: &models.ServiceError{TextError: "You cannot send coins to yourself", Code: 400},
		},
		{
			name:     "Not enough coins",
			fromUser: models.User{ID: 1, Login: "from", Password: "from", Coins: 100},
			toUser:   models.User{ID: 2, Login: "to", Password: "to", Coins: 100},
			amount:   200,
			mockBehavior: func(historyRep *mock_repository.MockHistory, userRep *mock_repository.MockUsers, fromUser *models.User, toUser *models.User, amount int) {
				userRep.EXPECT().GetUserByID(fromUser.ID).Return(fromUser, nil)
				userRep.EXPECT().GetUserByUsername(toUser.Login).Return(toUser, nil)
			},
			expectedError: &models.ServiceError{TextError: "Not enought coins", Code: 400},
		},
		{
			name:     "Error saving result in DB",
			fromUser: models.User{ID: 1, Login: "from", Password: "from", Coins: 100},
			toUser:   models.User{ID: 2, Login: "to", Password: "to", Coins: 100},
			amount:   10,
			mockBehavior: func(historyRep *mock_repository.MockHistory, userRep *mock_repository.MockUsers, fromUser *models.User, toUser *models.User, amount int) {
				userRep.EXPECT().GetUserByID(fromUser.ID).Return(fromUser, nil)
				userRep.EXPECT().GetUserByUsername(toUser.Login).Return(toUser, nil)
				userRep.EXPECT().UpdateUsers(&models.User{ID: fromUser.ID, Login: fromUser.Login, Password: fromUser.Password, Coins: fromUser.Coins - amount},
					&models.User{ID: toUser.ID, Login: toUser.Login, Password: toUser.Password, Coins: toUser.Coins + amount}).Return(errors.New("error"))
			},
			expectedError: &models.ServiceError{TextError: "error", Code: 500},
		},
		{
			name:     "Successful transaction",
			fromUser: models.User{ID: 1, Login: "from", Password: "from", Coins: 100},
			toUser:   models.User{ID: 2, Login: "to", Password: "to", Coins: 100},
			amount:   50,
			mockBehavior: func(historyRep *mock_repository.MockHistory, userRep *mock_repository.MockUsers, fromUser *models.User, toUser *models.User, amount int) {
				userRep.EXPECT().GetUserByID(fromUser.ID).Return(fromUser, nil)
				userRep.EXPECT().GetUserByUsername(toUser.Login).Return(toUser, nil)
				userRep.EXPECT().UpdateUsers(
					userMatcher{ID: fromUser.ID, Login: fromUser.Login, Password: fromUser.Password, Coins: fromUser.Coins - amount},
					userMatcher{ID: toUser.ID, Login: toUser.Login, Password: toUser.Password, Coins: toUser.Coins + amount},
				).Return(nil)
				historyRep.EXPECT().AddHistory(&models.TransactionsHistory{SenderID: fromUser.ID, ReceiverID: toUser.ID, Coins: amount}).Return(nil)
			},
			expectedError: nil,
		},
		{
			name:     "Zero coins",
			fromUser: models.User{ID: 1, Login: "from", Password: "from", Coins: 100},
			toUser:   models.User{ID: 2, Login: "to", Password: "to", Coins: 100},
			amount:   0,
			mockBehavior: func(historyRep *mock_repository.MockHistory, userRep *mock_repository.MockUsers, fromUser *models.User, toUser *models.User, amount int) {
			},
			expectedError: &models.ServiceError{TextError: "amount needs to be greater then 0", Code: 400},
		},
	}

	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUsersRepo := mock_repository.NewMockUsers(ctrl)
			mockHistoryRepo := mock_repository.NewMockHistory(ctrl)

			test.mockBehavior(mockHistoryRepo, mockUsersRepo, &test.fromUser, &test.toUser, test.amount)

			coinService := CoinService{repository: &repository.Repository{Users: mockUsersRepo, History: mockHistoryRepo}}

			err := coinService.SendCoins(test.fromUser.ID, test.toUser.Login, test.amount)

			if err == nil {
				require.Equal(t, test.expectedError, err)
			} else {
				require.Equal(t, test.expectedError.Code, err.Code)
			}
		})
	}

}

func TestCoinService_BuyItem(t *testing.T) {
	type mockBehavior func(merchRep *mock_repository.MockMerch, userRep *mock_repository.MockUsers, merch *models.Merch, user *models.User)

	testTable := []struct {
		name          string
		user          models.User
		merch         models.Merch
		mockBehavior  mockBehavior
		expectedError *models.ServiceError
	}{
		{
			name:  "Successful Purchase",
			user:  models.User{ID: 1, Login: "from", Password: "from", Coins: 100},
			merch: models.Merch{ID: 1, Name: "test_item", Price: 100},
			mockBehavior: func(mockMerchRepo *mock_repository.MockMerch, mockUserRepo *mock_repository.MockUsers, merch *models.Merch, user *models.User) {
				mockMerchRepo.EXPECT().GetMerchByName(merch.Name).Return(merch, nil)
				mockUserRepo.EXPECT().GetUserByID(user.ID).Return(user, nil)
				mockMerchRepo.EXPECT().BuyMerch(merch, user).Return(nil)

			},
			expectedError: nil,
		},
		{
			name:  "Item Not Found",
			user:  models.User{ID: 1, Login: "from", Password: "from", Coins: 100},
			merch: models.Merch{ID: 1, Name: "test_item", Price: 100},
			mockBehavior: func(mockMerchRepo *mock_repository.MockMerch, mockUserRepo *mock_repository.MockUsers, merch *models.Merch, user *models.User) {
				mockMerchRepo.EXPECT().GetMerchByName(merch.Name).Return(nil, nil)
			},
			expectedError: &models.ServiceError{TextError: "Merch is not found", Code: 400},
		},
		{
			name:  "Insufficient Funds",
			user:  models.User{ID: 1, Login: "from", Password: "from", Coins: 100},
			merch: models.Merch{ID: 1, Name: "test_item", Price: 1000},
			mockBehavior: func(mockMerchRepo *mock_repository.MockMerch, mockUserRepo *mock_repository.MockUsers, merch *models.Merch, user *models.User) {
				mockMerchRepo.EXPECT().GetMerchByName(merch.Name).Return(merch, nil)
				mockUserRepo.EXPECT().GetUserByID(user.ID).Return(user, nil)
			},
			expectedError: &models.ServiceError{TextError: "Not enought coins to buy this item ", Code: 400},
		},
	}

	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockUsersRepo := mock_repository.NewMockUsers(ctrl)
			mockMerchRepo := mock_repository.NewMockMerch(ctrl)

			mockService := NewService(&repository.Repository{Merch: mockMerchRepo, Users: mockUsersRepo})
			test.mockBehavior(mockMerchRepo, mockUsersRepo, &test.merch, &test.user)
			err := mockService.BuyItem(test.merch.Name, test.user.ID)

			if err == nil {
				require.Equal(t, test.expectedError, err)
			} else {
				require.Equal(t, test.expectedError.Code, err.Code)
			}

		})
	}
}
