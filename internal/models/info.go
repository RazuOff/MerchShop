package models

type HistoryResponse struct {
	Coins       int         `json:"coins"`
	Inventory   []item      `json:"inventory"`
	CoinHistory coinHistory `json:"coinHistory"`
}

type item struct {
	Type     string `json:"type"`
	Quantity int    `json:"quantity"`
}

type coinHistory struct {
	Received []transaction `json:"received"`
	Sent     []transaction `json:"sent"`
}

type transaction struct {
	FromUser string `json:"fromUser"`
	ToUser   string `json:"toUser"`
	Amount   int    `json:"amount"`
}
