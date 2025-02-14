package models

type Merch struct {
	ID    int    `json:"id"`
	Name  string `gorm:"unique;not null" json:"name"`
	Price int    `json:"price"`
}

type TransactionsHistory struct {
	ID         int `json:"id"`
	SenderID   int `json:"sender_id"`
	ReceiverID int `json:"reciver_id"`
	Coins      int `json:"coins"`
}

type User struct {
	ID       int    `json:"id"`
	Login    string `gorm:"unique;not null" json:"login"`
	Password string `gorm:"not null" json:"password"`
	Coins    int    `gorm:"default:1000;check:coins >= 0;not null" json:"coins"`
}

type UserMerch struct {
	UserID  int `gorm:"primaryKey"`
	MerchID int `gorm:"primaryKey"`
	Amount  int `gorm:"default:1;check:amount > 0"`
}
