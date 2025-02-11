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
		history.Inventory = append(history.Inventory, models.Item{Type: merch.Name, Quantity: count})
	}

	transactionHistory, err := service.repository.GetUserHistory(user)

	//repos recievd coins

	var recievedTransactions []models.ReceivedTransaction
	var sentedTransactions []models.SentTransaction

	for _, transaction := range transactionHistory {
		if transaction.ReceiverID == userID {
			fromUser, err := service.repository.GetUsernameByID(transaction.SenderID)

			if err != nil {
				return nil, &models.ServiceError{TextError: err.Error(), Code: 500}
			}

			recievedTransactions = append(recievedTransactions,
				models.ReceivedTransaction{FromUser: fromUser, Amount: transaction.Coins})
		}

		if transaction.SenderID == userID {
			toUser, err := service.repository.GetUsernameByID(transaction.ReceiverID)

			if err != nil {
				return nil, &models.ServiceError{TextError: err.Error(), Code: 500}
			}

			sentedTransactions = append(sentedTransactions,
				models.SentTransaction{ToUser: toUser, Amount: transaction.Coins})
		}
	}

	history.CoinHistory.Received = recievedTransactions
	history.CoinHistory.Sent = sentedTransactions

	return &history, nil
}
