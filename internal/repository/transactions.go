package repository

import "gorm.io/gorm"

type TransactionsHistoryPostgre struct {
	DB *gorm.DB
}

func newTransactionsHistoryPostgre(db *gorm.DB) *TransactionsHistoryPostgre {
	return &TransactionsHistoryPostgre{DB: db}
}
