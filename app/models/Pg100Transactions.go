package models

import "time"

type Pg100Transactions struct {
	ID                 uint      `json:"id" gorm:"primarykey;type:int2"`
	UserID             string    `json:"user_id" gorm:"type:varchar(50);unique_index;not null"`
	Username           string    `json:"username" gorm:"type:varchar(50);unique_index;not null"`
	AgentID            string    `json:"agent_id" gorm:"type:varchar(50);unique_index;not null"`
	ProductId          string    `json:"product_id" gorm:"type:varchar(50);unique_index;not null"`
	WalletAmountBefore float32   `json:"wallet_amount_before" sql:"type:decimal(10,2);not null"`
	WalletAmountAfter  float32   `json:"wallet_amount_after" sql:"type:decimal(10,2);not null"`
	BetAmount          float32   `json:"bet_amount" sql:"type:decimal(10,2);not null"`
	PayoutAmount       float32   `json:"payouta_mount" sql:"type:decimal(10,2);not null"`
	RoundId            string    `json:"round_id" gorm:"type:varchar(50);not null"`
	TxnId              string    `json:"txn_id" gorm:"type:varchar(50);not null"`
	Status             string    `json:"status" gorm:"type:varchar(15);not null"`
	GameCode           string    `json:"game_code" gorm:"type:varchar(50);not null"`
	GameId             string    `json:"game_id" gorm:"type:varchar(50);not null"`
	PlayInfo           string    `json:"play_info" gorm:"type:varchar(50);not null"`
	IsEndRound         bool      `json:"is_end_round" gorm:"comment: 0,f is friend | 1 is agent;not null"`
	CreatedAt          time.Time `json:"created_at"`
}
