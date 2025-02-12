package service

import (
	"github.com/RazuOff/MerchShop/internal/models"
	"github.com/RazuOff/MerchShop/internal/repository"
)

type CoinService struct {
	repository *repository.Repository
}

func NewCoinService(repo *repository.Repository) *CoinService {
	return &CoinService{repository: repo}
}

func (service *CoinService) SendCoins(fromUserID int, toUserLogin string, amount int) *models.ServiceError {

	fromUser, err := service.repository.GetUserByID(fromUserID)
	if err != nil {
		return &models.ServiceError{TextError: err.Error(), Code: 500}
	}

	toUser, err := service.repository.GetUserByUsername(toUserLogin)
	if err != nil {
		return &models.ServiceError{TextError: err.Error(), Code: 500}
	}

	if toUser == nil {
		return &models.ServiceError{TextError: "ToUserLogin does not exists", Code: 400}
	}

	if fromUser.ID == toUser.ID {
		return &models.ServiceError{TextError: "You cannot send coins to yourself", Code: 400}
	}

	if fromUser.Coins-amount < 0 {
		return &models.ServiceError{TextError: "Not enought coins", Code: 400}
	}

	fromUser.Coins -= amount
	toUser.Coins += amount

	if err = service.repository.UpdateUser(fromUser); err != nil {
		return &models.ServiceError{TextError: err.Error(), Code: 500}
	}

	if err = service.repository.UpdateUser(toUser); err != nil {
		return &models.ServiceError{TextError: err.Error(), Code: 500}
	}

	var history models.TransactionsHistory

	history.SenderID = fromUser.ID
	history.ReceiverID = toUser.ID
	history.Coins = amount

	if err = service.repository.AddHistory(&history); err != nil {
		return &models.ServiceError{TextError: err.Error(), Code: 500}
	}

	return nil
}

func (service *CoinService) BuyItem(itemName string, userID int) *models.ServiceError {
	merch, err := service.repository.GetMerchByName(itemName)
	if err != nil {
		return &models.ServiceError{TextError: err.Error(), Code: 500}
	}

	if merch == nil {
		return &models.ServiceError{TextError: "Merch is not found", Code: 400}
	}

	user, err := service.repository.GetUserByID(userID)
	if err != nil {
		return &models.ServiceError{TextError: err.Error(), Code: 500}
	}

	if user.Coins-merch.Price < 0 {
		return &models.ServiceError{TextError: "Not enought coins to buy this item ", Code: 400}
	}

	user.Coins -= merch.Price

	if err := service.repository.AddMerchToUser(merch, user); err != nil {
		return &models.ServiceError{TextError: err.Error(), Code: 500}
	}

	if err := service.repository.UpdateUser(user); err != nil {
		return &models.ServiceError{TextError: err.Error(), Code: 500}
	}

	return nil
}
