package repository

import (
	"github.com/RazuOff/MerchShop/internal/models"
	"gorm.io/gorm"
)

type Users interface {
	GetUserByUsername(username string) (*models.User, error)
}

type Repository struct {
	Users
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{}
}
