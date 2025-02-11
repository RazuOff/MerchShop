package repository

import (
	"github.com/RazuOff/MerchShop/internal/models"
	"gorm.io/gorm"
)

type TransactionsHistoryPostgre struct {
	DB *gorm.DB
}

func newTransactionsHistoryPostgre(db *gorm.DB) *TransactionsHistoryPostgre {
	return &TransactionsHistoryPostgre{DB: db}
}

func (repo *TransactionsHistoryPostgre) UpdateHistory(history *models.TransactionsHistory) error {

}
func (repo *TransactionsHistoryPostgre) GetUserHistory(user *models.User) ([]models.TransactionsHistory, error) {

}
