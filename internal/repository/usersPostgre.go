package repository

import "gorm.io/gorm"

type UsersPostgre struct {
	DB *gorm.DB
}

func NewUsersPostgre(db *gorm.DB) *UsersPostgre {
	return &UsersPostgre{DB: db}
}
