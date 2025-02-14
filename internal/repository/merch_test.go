package repository

import (
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/RazuOff/MerchShop/internal/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestRepository_GetMerchByName(t *testing.T) {
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

	repo := &MerchPostgre{DB: gormDB}

	tests := []struct {
		name           string
		itemName       string
		mockSetup      func()
		expectedMerch  *models.Merch
		expectedErrStr error
	}{
		{
			name:     "Merch Found",
			itemName: "t-shirt",
			mockSetup: func() {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "merches" WHERE name = $1 ORDER BY "merches"."id" LIMIT $2`)).
					WithArgs("t-shirt", 1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "price"}).AddRow(1, "t-shirt", 20))
			},
			expectedMerch:  &models.Merch{ID: 1, Name: "t-shirt", Price: 20},
			expectedErrStr: nil,
		},
		{
			name:     "Merch Not Found",
			itemName: "NonExistent",
			mockSetup: func() {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "merches" WHERE name = $1 ORDER BY "merches"."id" LIMIT $2`)).
					WithArgs("NonExistent", 1).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			expectedMerch:  nil,
			expectedErrStr: nil,
		},
		{
			name:     "DB Error",
			itemName: "ErrorItem",
			mockSetup: func() {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "merches" WHERE name = $1 ORDER BY "merches"."id" LIMIT $2`)).
					WithArgs("ErrorItem", 1).
					WillReturnError(errors.New("db error"))
			},
			expectedMerch:  nil,
			expectedErrStr: fmt.Errorf("error fetching item by name %v", errors.New("db error")),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mockSetup()

			merch, err := repo.GetMerchByName(test.itemName)

			if test.expectedErrStr == nil {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}

			assert.Equal(t, test.expectedMerch, merch)
		})
	}

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_BuyMerch(t *testing.T) {
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

	repo := &MerchPostgre{DB: gormDB}

	user := &models.User{ID: 1, Login: "Test", Password: "Test", Coins: 100}
	merch := &models.Merch{ID: 1, Name: "testMerch", Price: 20}

	tests := []struct {
		name          string
		mockSetup     func()
		expectedError error
	}{
		{
			name: "Successful purchase - merch exists",
			mockSetup: func() {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "user_merches" WHERE user_id = $1 AND merch_id = $2 ORDER BY "user_merches"."user_id" LIMIT $3`)).
					WithArgs(user.ID, merch.ID, 1).
					WillReturnRows(sqlmock.NewRows([]string{"user_id", "merch_id", "amount"}).AddRow(user.ID, merch.ID, 1))

				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "user_merches" SET "amount"=$1 WHERE "user_id" = $2 AND "merch_id" = $3`)).
					WithArgs(2, user.ID, merch.ID).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "users" SET "login"=$1,"password"=$2,"coins"=$3 WHERE "id" = $4`)).
					WithArgs(user.Login, user.Password, user.Coins-merch.Price, user.ID).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit()
			},
			expectedError: nil,
		},
		{
			name: "Successful purchase - new merch for the user",
			mockSetup: func() {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "user_merches" WHERE user_id = $1 AND merch_id = $2 ORDER BY "user_merches"."user_id" LIMIT $3`)).
					WithArgs(user.ID, merch.ID, 1).
					WillReturnError(gorm.ErrRecordNotFound)

				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "user_merches" ("user_id","merch_id","amount") VALUES ($1,$2,$3)`)).
					WithArgs(user.ID, merch.ID, 1).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "users" SET "login"=$1,"password"=$2,"coins"=$3 WHERE "id" = $4`)).
					WithArgs(user.Login, user.Password, user.Coins-merch.Price, user.ID).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit()
			},
			expectedError: nil,
		},
		{
			name: "Error during transaction - DB error when saving user merch",
			mockSetup: func() {

				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "user_merches" WHERE user_id = $1 AND merch_id = $2 ORDER BY "user_merches"."user_id" LIMIT $3`)).
					WithArgs(user.ID, merch.ID, 1).
					WillReturnRows(sqlmock.NewRows([]string{"user_id", "merch_id", "amount"}).AddRow(user.ID, merch.ID, 1))

				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "user_merches" SET "amount"=$1 WHERE "user_id" = $2 AND "merch_id" = $3`)).
					WithArgs(2, user.ID, merch.ID).
					WillReturnError(errors.New("db error"))

				mock.ExpectRollback()
			},
			expectedError: errors.New("db error"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mockSetup()

			err := repo.BuyMerch(merch, user)

			if test.expectedError == nil {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Equal(t, test.expectedError.Error(), err.Error())
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestRepository_GetUsersMerch(t *testing.T) {
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

	repo := &MerchPostgre{DB: gormDB}

	tests := []struct {
		name           string
		user           *models.User
		mockSetup      func()
		expectedMerch  []models.Merch
		expectedErrStr error
	}{
		{
			name: "Successful Get Users Merch",
			user: &models.User{ID: 1},
			mockSetup: func() {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT "merches"."id","merches"."name","merches"."price" FROM "merches"
				 JOIN user_merches ON user_merches.merch_id = merches.id WHERE user_merches.user_id = $1`)).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "price"}).
						AddRow(1, "t-shirt", 20).
						AddRow(2, "hat", 15))
			},
			expectedMerch: []models.Merch{
				{ID: 1, Name: "t-shirt", Price: 20},
				{ID: 2, Name: "hat", Price: 15},
			},
			expectedErrStr: nil,
		},
		{
			name: "No Merch Found",
			user: &models.User{ID: 2},
			mockSetup: func() {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT "merches"."id","merches"."name","merches"."price" FROM "merches"
				 JOIN user_merches ON user_merches.merch_id = merches.id WHERE user_merches.user_id = $1`)).
					WithArgs(2).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "price"}))
			},
			expectedMerch:  []models.Merch{},
			expectedErrStr: nil,
		},
		{
			name: "Error while fetching Merch",
			user: &models.User{ID: 1},
			mockSetup: func() {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT "merches"."id","merches"."name","merches"."price" FROM "merches"
				 JOIN user_merches ON user_merches.merch_id = merches.id WHERE user_merches.user_id = $1`)).
					WithArgs(1).
					WillReturnError(fmt.Errorf("db error"))
			},
			expectedMerch:  nil,
			expectedErrStr: fmt.Errorf("failed to get user's merch: %v", fmt.Errorf("db error")),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mockSetup()

			merchList, err := repo.GetUsersMerch(test.user)

			if test.expectedErrStr == nil {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Equal(t, test.expectedErrStr.Error(), err.Error())
			}

			assert.Equal(t, test.expectedMerch, merchList)
		})
	}

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_GetUserMerchAmount(t *testing.T) {
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

	repo := &MerchPostgre{DB: gormDB}

	tests := []struct {
		name           string
		user           *models.User
		merch          *models.Merch
		mockSetup      func()
		expectedAmount int
		expectedErrStr error
	}{
		{
			name:  "User has merch",
			user:  &models.User{ID: 1},
			merch: &models.Merch{ID: 1},
			mockSetup: func() {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "user_merches" WHERE user_id = $1 AND merch_id = $2 ORDER BY "user_merches"."user_id" LIMIT $3`)).
					WithArgs(1, 1, 1).
					WillReturnRows(sqlmock.NewRows([]string{"user_id", "merch_id", "amount"}).AddRow(1, 1, 5))
			},
			expectedAmount: 5,
			expectedErrStr: nil,
		},
		{
			name:  "User has no merch",
			user:  &models.User{ID: 2},
			merch: &models.Merch{ID: 2},
			mockSetup: func() {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "user_merches" WHERE user_id = $1 AND merch_id = $2 ORDER BY "user_merches"."user_id" LIMIT $3`)).
					WithArgs(2, 2, 1).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			expectedAmount: 0,
			expectedErrStr: nil,
		},
		{
			name:  "Error while fetching merch amount",
			user:  &models.User{ID: 1},
			merch: &models.Merch{ID: 1},
			mockSetup: func() {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "user_merches" WHERE user_id = $1 AND merch_id = $2 ORDER BY "user_merches"."user_id" LIMIT $3`)).
					WithArgs(1, 1, 1).
					WillReturnError(fmt.Errorf("db error"))
			},
			expectedAmount: 0,
			expectedErrStr: fmt.Errorf("failed to get merch amount: %v", fmt.Errorf("db error")),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mockSetup()

			amount, err := repo.GetUserMerchAmount(test.user, test.merch)

			if test.expectedErrStr == nil {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Equal(t, test.expectedErrStr.Error(), err.Error())
			}

			assert.Equal(t, test.expectedAmount, amount)
		})
	}

	assert.NoError(t, mock.ExpectationsWereMet())
}
