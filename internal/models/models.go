package models

type Merch struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
	User  []User `gorm:"many2many:user_merch;constraint:OnDelete:CASCADE;"`
}

type TransactionsHistory struct {
	ID        int `json:"id"`
	SenderID  int `json:"sender_id"`
	ReciverID int `json:"reciver_id"`
	Coins     int `json:"coins"`
}

type User struct {
	ID       int     `json:"id"`
	Login    string  `json:"login"`
	Password string  `json:"password"`
	Coins    int     `gorm:"default:1000;check:coins >= 0" json:"coins"`
	Merch    []Merch `gorm:"many2many:user_merch;constraint:OnDelete:CASCADE;"`
}
