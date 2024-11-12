package models

import "time"

type BalanceCheckResponse struct {
	ID              string  `json:"id"`
	StatusCode      int     `json:"statusCode"`
	TimestampMillis int64   `json:"timestampMillis"`
	ProductId       string  `json:"productId"`
	Currency        string  `json:"currency"`
	Balance         float32 `json:"balance"`
	Username        string  `json:"username"`
}

type SettleCheckResponse struct {
	ID              string        `json:"id"`
	StatusCode      int           `json:"statusCode"`
	TimestampMillis int64         `json:"timestampMillis"`
	ProductId       string        `json:"productId"`
	Currency        string        `json:"currency"`
	BalanceBefore   float32       `json:"balanceBefore"`
	BalanceAfter    float32       `json:"balanceAfter"`
	Username        string        `json:"username"`
	Transactions    []Transaction `json:"txns"`
}

type SettleCheckResponseFormBackend struct {
	Status        string  `json:"status"`
	UserID        uint    `json:"user_id"`
	AgentID       uint    `json:"agent_id"`
	KeyDeposit    string  `json:"key_deposit"`
	Username      string  `json:"username"`
	UsernameAgent string  `json:"username_agent"`
	BalanceBefore float32 `json:"amount_before"`
	BalanceAfter  float32 `json:"amount_after"`
}

type Transaction struct {
	ID            string  `json:"id"`
	Status        string  `json:"status"`
	RoundID       string  `json:"roundId"`
	BetAmount     float32 `json:"betAmount"`
	PayoutAmount  float32 `json:"payoutAmount"`
	GameCode      string  `json:"gameCode"`
	PlayInfo      string  `json:"playInfo"`
	TxnID         string  `json:"txnId"`
	IsFreeSpin    bool    `json:"isFreeSpin"`
	BuyFeature    bool    `json:"buyFeature"`
	BonusFreeSpin bool    `json:"bonusFreeSpin"`
	IsEndRound    bool    `json:"isEndRound"`
}

type ResponseData struct {
	Data struct {
		Balance  float32 `json:"balance"`
		Username string  `json:"username"`
	} `json:"data"`
	Message string `json:"message"`
	Status  string `json:"status"`
	Time    string `json:"time"`
}

type ResponseDataSettle struct {
	Data    SettleCheckResponseFormBackend `json:"data"`
	Message string                         `json:"message"`
	Status  string                         `json:"status"`
	Time    string                         `json:"time"`
}

type Pg100Transactions struct {
	ID                 uint      `json:"id" gorm:"primarykey;type:int2"`
	UserID             uint      `json:"user_id" gorm:"type:int2;index;not null"`
	AgentID            uint      `json:"agent_id" gorm:"type:int2;index;not null"`
	Username           string    `json:"username" gorm:"type:varchar(50);index;not null"`
	KeyDeposit         string    `json:"key_deposit" gorm:"type:varchar(50);index;not null"`
	ProductId          string    `json:"product_id" gorm:"type:varchar(50);index;not null"`
	WalletAmountBefore float32   `json:"wallet_amount_before" gorm:"type:decimal(10,2);not null"`
	WalletAmountAfter  float32   `json:"wallet_amount_after" gorm:"type:decimal(10,2);not null"`
	BetAmount          float32   `json:"bet_amount" gorm:"type:decimal(10,2);not null"`
	PayoutAmount       float32   `json:"payouta_mount" gorm:"type:decimal(10,2);not null"`
	RoundId            string    `json:"round_id" gorm:"type:varchar(50);not null"`
	TxnId              string    `json:"txn_id" gorm:"type:varchar(50);not null"`
	Status             string    `json:"status" gorm:"type:varchar(15);not null"`
	GameCode           string    `json:"game_code" gorm:"type:varchar(50);not null"`
	GameId             string    `json:"game_id" gorm:"type:varchar(50);not null"`
	PlayInfo           string    `json:"play_info" gorm:"type:varchar(50);not null"`
	IsFreeSpin         bool      `json:"is_free_spin" gorm:"comment: 0,f is friend | 1 is agent;not null"`
	BuyFeature         bool      `json:"buy_feature" gorm:"comment: 0,f is friend | 1 is agent;not null"`
	BonusFreeSpin      bool      `json:"bonus_free_spin" gorm:"comment: 0,f is friend | 1 is agent;not null"`
	IsEndRound         bool      `json:"is_end_round" gorm:"comment: 0,f is friend | 1 is agent;not null"`
	CreatedAt          time.Time `json:"created_at"`
}
