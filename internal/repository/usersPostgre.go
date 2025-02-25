package repository

import (
	"fmt"
	"log"

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
	var user models.User

	if err := repo.DB.Where("login = ?", username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user: %v", err)
	}

	return &user, nil
}

func (repo *UsersPostgre) SetUser(username string, password string) (*models.User, error) {
	user := &models.User{
		Login:    username,
		Password: password,
	}

	if err := repo.DB.Create(user).Error; err != nil {
		return nil, fmt.Errorf("failed to create user: %v", err)
	}

	return user, nil
}

func (repo *UsersPostgre) GetUserByID(userID int) (*models.User, error) {
	var user models.User

	if err := repo.DB.First(&user, userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user: %v", err)
	}

	return &user, nil
}

func (repo *UsersPostgre) UpdateUsers(users ...*models.User) error {
	if len(users) == 0 {
		return fmt.Errorf("users")
	}

	tx := repo.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	for _, user := range users {
		result := tx.Model(&models.User{}).Where("id = ?", user.ID).Updates(user)
		if result.Error != nil {
			tx.Rollback()
			return fmt.Errorf("failed to update user: %v", result.Error)
		}

		if result.RowsAffected == 0 {
			log.Println(user.Coins)
			tx.Rollback()
			return fmt.Errorf("no user found with ID %d", user.ID)
		}
	}
	tx.Commit()

	return nil
}

func (repo *UsersPostgre) GetUsernameByID(userID int) (string, error) {
	var user models.User

	if err := repo.DB.First(&user, userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", fmt.Errorf("user with ID %d not found", userID)
		}
		return "", fmt.Errorf("failed to get user: %v", err)
	}

	return user.Login, nil
}
