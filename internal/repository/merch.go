package repository

import "gorm.io/gorm"

type MerchPostgre struct {
	DB *gorm.DB
}

func NewMerchPostgre(db *gorm.DB) *MerchPostgre {
	return &MerchPostgre{DB: db}
}
