package repository

import (
	"github.com/RazuOff/MerchShop/internal/models"
	"gorm.io/gorm"
)

type UsersPostgre struct {
	DB *gorm.DB
}

func newUsersPostgre(db *gorm.DB) *UsersPostgre {
	return &UsersPostgre{DB: db}
}

func (postgre *UsersPostgre) GetUserByUsername(username string) models.User {

}
