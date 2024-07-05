package models

import "time"

type Reports struct {
	ID                 uint      `json:"id" gorm:"primarykey;type:int2"`
	UserID             string    `json:"user_id" gorm:"type:varchar(50);unique_index;not null"`
	Username           string    `json:"username" gorm:"type:varchar(50);unique_index;not null"`
	AgentID            string    `json:"agent_id" gorm:"type:varchar(50);unique_index;not null"`
	RoundId            string    `json:"round_id" gorm:"type:varchar(50);not null"`
	ProductId          string    `json:"product_id" gorm:"type:varchar(50);unique_index;not null"`
	ProductName        string    `json:"product_name" gorm:"type:varchar(50);unique_index;not null"`
	GameId             string    `json:"game_id" gorm:"type:varchar(50);unique_index;not null"`
	GameName           string    `json:"game_name" gorm:"type:varchar(50);unique_index;not null"`
	WalletAmountBefore float32   `json:"wallet_amount_before" sql:"type:decimal(10,2);not null"`
	WalletAmountAfter  float32   `json:"wallet_amount_after" sql:"type:decimal(10,2);not null"`
	BetAmount          float32   `json:"bet_amount" sql:"type:decimal(10,2);not null"`
	BetResult          float32   `json:"bet_result" sql:"type:decimal(10,2);not null"`
	BetWinloss         float32   `json:"bet_winloss" sql:"type:decimal(10,2);not null"`
	Status             string    `json:"status" gorm:"type:varchar(15);not null"`
	IP                 string    `json:"ip" gorm:"type:varchar(15);not null"`
	Description        string    `json:"description" gorm:"type:varchar(200);not null"`
	CreatedAt          time.Time `json:"created_at"`
}
