package service

import (
	"github.com/RazuOff/MerchShop/internal/config"
	"github.com/RazuOff/MerchShop/internal/models"
	"github.com/RazuOff/MerchShop/internal/repository"
)

type Auth interface {
	RegistrateOrLogin(username string, password string) (*models.User, *models.ServiceError)
	GenerateToken(username string, userID int, config config.Config) (string, *models.ServiceError)
}

type Info interface {
	GenerateInfo(userID int) (*models.HistoryResponse, *models.ServiceError)
}

type Coin interface {
	SendCoins(fromUserID int, toUserLogin string, amount int) *models.ServiceError
	BuyItem(itemName string, userID int) *models.ServiceError
}

type Service struct {
	Auth
	Info
	Coin
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Auth: NewAuthService(repo),
		Info: NewInfoService(repo),
		Coin: NewCoinService(repo),
	}
}
