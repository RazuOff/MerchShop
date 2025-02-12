package service

import (
	"github.com/RazuOff/MerchShop/internal/models"
	"github.com/RazuOff/MerchShop/internal/repository"
)

type InfoService struct {
	repository *repository.Repository
}

func NewInfoService(repo *repository.Repository) *InfoService {
	return &InfoService{repository: repo}
}

func (service *InfoService) GenerateInfo(userID int) (*models.HistoryResponse, *models.ServiceError) {
	var history models.HistoryResponse
	history.CoinHistory = models.CoinHistory{Received: []models.ReceivedTransaction{}, Sent: []models.SentTransaction{}}
	history.Inventory = []models.Item{}

	user, err := service.repository.GetUserByID(userID)
	if err != nil {
		return nil, &models.ServiceError{TextError: err.Error(), Code: 500}
	}
	if user == nil {
		return nil, &models.ServiceError{TextError: "User not found", Code: 500}
	}

	history.Coins = user.Coins
	usersMerch, err := service.repository.GetUsersMerch(user)
	if err != nil {
		return nil, &models.ServiceError{TextError: err.Error(), Code: 500}
	}

	for _, merch := range usersMerch {
		count, err := service.repository.GetUserMerchAmount(user, &merch)
		if err != nil {
			return nil, &models.ServiceError{TextError: err.Error(), Code: 500}
		}
		if count == 0 {
			continue
		}

		history.Inventory = append(history.Inventory, models.Item{Type: merch.Name, Quantity: count})
	}

	transactionHistory, err := service.repository.GetUserHistory(user)
	if err != nil {
		return nil, &models.ServiceError{TextError: err.Error(), Code: 500}
	}

	if len(transactionHistory) == 0 {
		return &history, nil
	}

	recievedTransactions := make(map[string]int)
	sentedTransactions := make(map[string]int)

	for _, transaction := range transactionHistory {
		if transaction.ReceiverID == userID {
			fromUser, err := service.repository.GetUsernameByID(transaction.SenderID)

			if err != nil {
				return nil, &models.ServiceError{TextError: err.Error(), Code: 500}
			}
			recievedTransactions[fromUser] += transaction.Coins
		}

		if transaction.SenderID == userID {
			toUser, err := service.repository.GetUsernameByID(transaction.ReceiverID)

			if err != nil {
				return nil, &models.ServiceError{TextError: err.Error(), Code: 500}
			}
			sentedTransactions[toUser] += transaction.Coins
		}
	}

	for fromUser, amount := range recievedTransactions {
		history.CoinHistory.Received = append(history.CoinHistory.Received, models.ReceivedTransaction{FromUser: fromUser, Amount: amount})
	}

	for toUser, amount := range sentedTransactions {
		history.CoinHistory.Sent = append(history.CoinHistory.Sent, models.SentTransaction{ToUser: toUser, Amount: amount})
	}

	return &history, nil
}
