package database

import "github.com/RazuOff/MerchShop/internal/models"

func getTestMerchTable() *[]models.Merch {
	table := []models.Merch{
		{
			Name:  "t-shirt",
			Price: 80,
		},
		{
			Name:  "cup",
			Price: 20,
		},
		{
			Name:  "book",
			Price: 50,
		},
		{
			Name:  "pen",
			Price: 10,
		},
		{
			Name:  "powerbank",
			Price: 200,
		},
		{
			Name:  "hoody",
			Price: 300,
		},
		{
			Name:  "umbrella",
			Price: 200,
		},
		{
			Name:  "socks",
			Price: 10,
		},
		{
			Name:  "wallet",
			Price: 50,
		},
		{
			Name:  "pink-hoody",
			Price: 500,
		},
	}

	return &table
}
