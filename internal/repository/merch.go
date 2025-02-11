package repository

import (
	"github.com/RazuOff/MerchShop/internal/models"
	"gorm.io/gorm"
)

type MerchPostgre struct {
	DB *gorm.DB
}

func newMerchPostgre(db *gorm.DB) *MerchPostgre {
	return &MerchPostgre{DB: db}
}

func (repo *MerchPostgre) GetMerchByName(itemName string) (*models.Merch, error) {

}
func (repo *MerchPostgre) AddMerchToUser(merch *models.Merch, user *models.User) error {

}
func (repo *MerchPostgre) GetUsersMerch(user *models.User) ([]models.Merch, error) {

}
func (repo *MerchPostgre) GetUserMerchAmount(user *models.User, merch *models.Merch) (int, error) {

}
