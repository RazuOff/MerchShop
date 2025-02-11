package repository

import (
	"github.com/RazuOff/MerchShop/internal/models"
	"gorm.io/gorm"
)

type Users interface {
	GetUserByUsername(username string) (*models.User, error)
	SetUser(username string, password string) (*models.User, error)
	GetUserByID(userID int) (*models.User, error)
	UpdateUser(user *models.User) error
	GetUsernameByID(userID int) (string, error)
}

type History interface {
	UpdateHistory(history *models.TransactionsHistory) error
	GetUserHistory(user *models.User) ([]models.TransactionsHistory, error)
}

type Merch interface {
	GetMerchByName(itemName string) (*models.Merch, error)
	AddMerchToUser(merch *models.Merch, user *models.User) error
	GetUsersMerch(user *models.User) ([]models.Merch, error)
	GetUserMerchAmount(user *models.User, merch *models.Merch) (int, error)
}

type Repository struct {
	Users
	History
	Merch
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{}
}
