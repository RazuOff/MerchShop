package database

import (
	"github.com/RazuOff/MerchShop/internal/models"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func NewTestDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	db.Exec("TRUNCATE TABLE users, merch, user_merch RESTART IDENTITY CASCADE")
	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(&models.User{}, &models.Merch{}, &models.UserMerch{}, &models.TransactionsHistory{})

	return db, err
}
