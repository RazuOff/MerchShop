package repository

import (
	"fmt"

	"github.com/RazuOff/MerchShop/internal/models"
	"gorm.io/gorm"
)

type TransactionsHistoryPostgre struct {
	DB *gorm.DB
}

func NewTransactionsHistoryPostgre(db *gorm.DB) *TransactionsHistoryPostgre {
	return &TransactionsHistoryPostgre{DB: db}
}

func (repo *TransactionsHistoryPostgre) AddHistory(history *models.TransactionsHistory) error {

	if history == nil {
		return fmt.Errorf("transaction history is nil")
	}

	if err := repo.DB.Create(history).Error; err != nil {
		return fmt.Errorf("failed to add transaction history: %w", err)
	}
	return nil
}

func (repo *TransactionsHistoryPostgre) GetUserHistory(user *models.User) ([]models.TransactionsHistory, error) {
	var history []models.TransactionsHistory

	if err := repo.DB.Where("sender_id = ? OR receiver_id = ?", user.ID, user.ID).
		Find(&history).Error; err != nil {
		return nil, fmt.Errorf("failed to get user transaction history: %v", err)
	}

	return history, nil
}
