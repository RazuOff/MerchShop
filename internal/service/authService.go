package service

import (
	"time"

	"github.com/RazuOff/MerchShop/internal/config"
	"github.com/RazuOff/MerchShop/internal/models"
	"github.com/RazuOff/MerchShop/internal/repository"
	"github.com/golang-jwt/jwt"
)

type AuthService struct {
	repository repository.Users
}

func NewAuthService(repo *repository.Repository) *AuthService {
	return &AuthService{repository: repo.Users}
}

func (service *AuthService) GenerateToken(username string, userID int, config config.Config) (string, *models.ServiceError) {

	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &models.Claims{
		ID:       userID,
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(config.JwtKey)
	if err != nil {
		return "", &models.ServiceError{TextError: err.Error(), Code: 500}
	}
	return tokenString, nil

}

func (service *AuthService) RegistrateOrLogin(username string, password string) (*models.User, *models.ServiceError) {
	user, err := service.repository.GetUserByUsername(username)
	if err != nil {
		return nil, &models.ServiceError{TextError: err.Error(), Code: 500}
	}

	if user == nil {
		user, err := service.registrate(username, password)
		if err != nil {
			return nil, err
		}
		return user, nil
	} else {
		return user, nil
	}
}

func (service *AuthService) registrate(username string, password string) (*models.User, *models.ServiceError) {
	user, err := service.repository.SetUser(username, password)
	if err != nil {
		return nil, &models.ServiceError{TextError: err.Error(), Code: 500}
	}

	return user, nil
}
