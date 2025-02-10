package service

import "github.com/RazuOff/MerchShop/internal/repository"

type Service struct {
}

func NewService(repo *repository.Repository) *Service {
	return &Service{}
}
