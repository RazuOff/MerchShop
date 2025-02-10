package service

import (
	"github.com/RazuOff/MerchShop/internal/config"
	"github.com/RazuOff/MerchShop/internal/models"
	"github.com/RazuOff/MerchShop/internal/repository"
)

type AuthService interface {
	RegistrateOrLogin(username string, password string) (models.User, error)
	GenerateToken(username string, userID int, config config.Config) (string, error)
}

type InfoService interface {
	GenerateInfo(userID string) (models.HistoryResponse, error)
}

type CoinService interface {
	SendCoins(fromUserID string, toUserID string, amount int) error
	BuyItem(itemName string, userID string) (bool, error)
}

type Service struct {
	AuthService
	InfoService
	CoinService
}

func NewService(repo *repository.Repository) *Service {
	return &Service{}
}
