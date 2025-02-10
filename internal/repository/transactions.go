package repository

import "gorm.io/gorm"

type TransactionsHistoryPostgre struct {
	DB *gorm.DB
}

func NewTransactionsHistoryPostgre(db *gorm.DB) *TransactionsHistoryPostgre {
	return &TransactionsHistoryPostgre{DB: db}
}
