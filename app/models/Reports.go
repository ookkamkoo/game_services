package models

import "time"

type BetWinLossSummary struct {
	UserID     uint    `json:"user_id"`
	BetWinloss float32 `json:"bet_winloss"`
}

type Reports struct {
	ID                 uint      `json:"id" gorm:"primarykey;type:int2"`
	UserID             uint      `json:"user_id" gorm:"type:int2;index;not null"`
	AgentID            uint      `json:"agent_id" gorm:"type:int2;index;not null"`
	Username           string    `json:"username" gorm:"type:varchar(50);index;not null"`
	RoundId            string    `json:"round_id" gorm:"type:varchar(50);not null"`
	ProductId          string    `json:"product_id" gorm:"type:varchar(50);index;not null"`
	ProductName        string    `json:"product_name" gorm:"type:varchar(50);index;not null"`
	GameId             string    `json:"game_id" gorm:"type:varchar(50);index;not null"`
	GameName           string    `json:"game_name" gorm:"type:varchar(50);index;not null"`
	WalletAmountBefore float32   `json:"wallet_amount_before" gorm:"type:decimal(10,2);not null"`
	WalletAmountAfter  float32   `json:"wallet_amount_after" gorm:"type:decimal(10,2);not null"`
	BetAmount          float32   `json:"bet_amount" gorm:"type:decimal(10,2);not null"`
	BetResult          float32   `json:"bet_result" gorm:"type:decimal(10,2);not null"`
	BetWinloss         float32   `json:"bet_winloss" gorm:"type:decimal(10,2);not null"`
	Status             string    `json:"status" gorm:"type:varchar(15);not null"`
	IP                 string    `json:"ip" gorm:"type:varchar(15);not null"`
	Description        string    `json:"description" gorm:"type:varchar(200)"`
	CreatedAt          time.Time `json:"created_at"`
}
