package database

import (
	"fmt"

	"github.com/RazuOff/MerchShop/internal/config"
	"github.com/RazuOff/MerchShop/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgre(config *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(config.PostgreString))
	if err != nil {
		return nil, fmt.Errorf("database connect error: %w", err)
	}

	if err := db.AutoMigrate(&models.User{}, &models.Merch{}, &models.TransactionsHistory{}); err != nil {
		return nil, fmt.Errorf("database connect error: %w", err)
	}
	db.Create(getTestMerchTable())

	return db, nil
}
