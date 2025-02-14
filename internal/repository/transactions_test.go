package repository

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/RazuOff/MerchShop/internal/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestRepository_AddHistory(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unexpected error when opening mock database: %s", err)
	}
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("unexpected error when opening GORM DB: %s", err)
	}

	repo := &TransactionsHistoryPostgre{DB: gormDB}

	tests := []struct {
		name          string
		history       *models.TransactionsHistory
		mockSetup     func()
		expectedError error
	}{
		{
			name: "Successful history addition",
			history: &models.TransactionsHistory{
				ID:         1,
				SenderID:   1,
				ReceiverID: 2,
				Coins:      100,
			},
			mockSetup: func() {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "transactions_histories" ("sender_id","receiver_id","coins","id") VALUES ($1,$2,$3,$4) RETURNING "id"`)).
					WithArgs(1, 2, 100, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				mock.ExpectCommit()
			},
			expectedError: nil,
		},
		{
			name:          "Nil history object",
			history:       nil,
			mockSetup:     func() {},
			expectedError: fmt.Errorf("transaction history is nil"),
		},
		{
			name: "Database error",
			history: &models.TransactionsHistory{
				ID:         2,
				SenderID:   3,
				ReceiverID: 4,
				Coins:      50,
			},
			mockSetup: func() {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "transactions_histories" ("sender_id","receiver_id","coins","id") VALUES ($1,$2,$3,$4) RETURNING "id"`)).
					WithArgs(3, 4, 50, 2).
					WillReturnError(fmt.Errorf("db error"))
				mock.ExpectRollback()
			},
			expectedError: fmt.Errorf("failed to add transaction history: db error"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mockSetup()

			err := repo.AddHistory(test.history)

			if test.expectedError == nil {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Equal(t, test.expectedError.Error(), err.Error())
			}
		})
	}

	assert.NoError(t, mock.ExpectationsWereMet())
}
