package repository

import (
	"fmt"

	"github.com/RazuOff/MerchShop/internal/models"
	"gorm.io/gorm"
)

type MerchPostgre struct {
	DB *gorm.DB
}

func NewMerchPostgre(db *gorm.DB) *MerchPostgre {
	return &MerchPostgre{DB: db}
}

func (repo *MerchPostgre) GetMerchByName(itemName string) (*models.Merch, error) {
	var merch models.Merch
	if err := repo.DB.Where("name = ?", itemName).First(&merch).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("error fetching item by name: %v", err)
	}
	return &merch, nil

}
func (repo *MerchPostgre) AddMerchToUser(merch *models.Merch, user *models.User) error {
	var userMerch models.UserMerch
	err := repo.DB.Where("user_id = ? AND merch_id = ?", user.ID, merch.ID).First(&userMerch).Error
	if err == nil {
		userMerch.Amount++
		return repo.DB.Save(&userMerch).Error
	} else if err == gorm.ErrRecordNotFound {
		return repo.DB.Create(&models.UserMerch{
			UserID:  user.ID,
			MerchID: merch.ID,
			Amount:  1,
		}).Error
	}

	return err
}
func (repo *MerchPostgre) GetUsersMerch(user *models.User) ([]models.Merch, error) {
	if user == nil {
		return nil, fmt.Errorf("invalid user")
	}

	var merchList []models.Merch
	err := repo.DB.
		Joins("JOIN user_merches ON user_merches.merch_id = merches.id").
		Where("user_merches.user_id = ?", user.ID).
		Find(&merchList).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get user's merch: %v", err)
	}

	return merchList, nil
}
func (repo *MerchPostgre) GetUserMerchAmount(user *models.User, merch *models.Merch) (int, error) {
	var userMerch models.UserMerch

	err := repo.DB.Where("user_id = ? AND merch_id = ?", user.ID, merch.ID).First(&userMerch).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, nil
		}
		return 0, fmt.Errorf("failed to get merch amount: %v", err)
	}

	return userMerch.Amount, nil
}
