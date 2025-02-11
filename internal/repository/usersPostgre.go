package repository

import (
	"github.com/RazuOff/MerchShop/internal/models"
	"gorm.io/gorm"
)

type UsersPostgre struct {
	DB *gorm.DB
}

func NewUsersPostgre(db *gorm.DB) *UsersPostgre {
	return &UsersPostgre{DB: db}
}

func (repo *UsersPostgre) GetUserByUsername(username string) (*models.User, error) {

}
func (repo *UsersPostgre) SetUser(username string, password string) (*models.User, error) {

}
func (repo *UsersPostgre) GetUserByID(userID int) (*models.User, error) {

}
func (repo *UsersPostgre) UpdateUser(user *models.User) error {

}
func (repo *UsersPostgre) GetUsernameByID(userID int) (string, error) {

}
