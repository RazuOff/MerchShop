package repository

import (
	"github.com/RazuOff/MerchShop/internal/models"
	"gorm.io/gorm"
)

//go:generate mockgen -source=repository.go -destination=mocks/mock.go

type Users interface {
	GetUserByUsername(username string) (*models.User, error)
	SetUser(username string, password string) (*models.User, error)
	GetUserByID(userID int) (*models.User, error)
	UpdateUsers(users ...*models.User) error
	GetUsernameByID(userID int) (string, error)
}

type History interface {
	AddHistory(history *models.TransactionsHistory) error
	GetUserHistory(user *models.User) ([]models.TransactionsHistory, error)
}

type Merch interface {
	GetMerchByName(itemName string) (*models.Merch, error)
	BuyMerch(merch *models.Merch, user *models.User) error
	GetUsersMerch(user *models.User) ([]models.Merch, error)
	GetUserMerchAmount(user *models.User, merch *models.Merch) (int, error)
}

type Repository struct {
	Users
	History
	Merch
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		Users:   NewUsersPostgre(db),
		History: NewTransactionsHistoryPostgre(db),
		Merch:   NewMerchPostgre(db),
	}
}
