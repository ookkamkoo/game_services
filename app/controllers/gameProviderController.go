package controllers

import (
	"fmt"
	"game_services/app/utils"
	"time"

	"github.com/gofiber/fiber/v2"
)

type BalanceRequest struct {
	AgentUsername  string `json:"agentUsername"`
	OperatorToken  string `json:"operatorToken"`
	SeamlessKey    string `json:"seamlessKey"`
	PlayerUsername string `json:"playerUsername"`
	CurrencyCode   string `json:"currencyCode"`
	ProductName    string `json:"productName"`
	ProductId      int    `json:"productId"`
	ProductCode    string `json:"productCode"`
	EventType      int    `json:"eventType"`
	EventName      string `json:"eventName"`
	RequestUid     string `json:"requestUid"`
	RequestTime    string `json:"requestTime"`
	Timestamp      int64  `json:"timestamp"`
}

type EventDetail struct {
	IsFeature                    bool    `json:"isFeature"`
	IsFeatureBuy                 bool    `json:"isFeatureBuy"`
	JackpotRtpContributionAmount float64 `json:"jackpotRtpContributionAmount"`
	JackpotWinAmount             float64 `json:"jackpotWinAmount"`
}

type TransactionRequest struct {
	AgentUsername  string      `json:"agentUsername"`
	Amount         float64     `json:"amount"`
	CategoryId     int         `json:"categoryId"`
	CategoryName   string      `json:"categoryName"`
	CurrencyCode   string      `json:"currencyCode"`
	EventDetail    EventDetail `json:"eventDetail"`
	EventName      string      `json:"eventName"`
	EventType      int         `json:"eventType"`
	GameCode       int         `json:"gameCode"`
	GameName       string      `json:"gameName"`
	IsRefund       bool        `json:"isRefund"`
	OperatorToken  string      `json:"operatorToken"`
	PlayerUsername string      `json:"playerUsername"`
	ProductCode    string      `json:"productCode"`
	ProductId      int         `json:"productId"`
	ProductName    string      `json:"productName"`
	RequestTime    string      `json:"requestTime"`
	RequestUid     string      `json:"requestUid"`
	RoundId        string      `json:"roundId"`
	SeamlessKey    string      `json:"seamlessKey"`
	Timestamp      float64     `json:"timestamp"`
	TxnId          string      `json:"txnId"`
}

type CreditRequest struct {
	OperatorToken  string  `json:"operatorToken"`
	SeamlessKey    string  `json:"seamlessKey"`
	AgentUsername  string  `json:"agentUsername"`
	PlayerUsername string  `json:"playerUsername"`
	CurrencyCode   string  `json:"currencyCode"`
	ProductName    string  `json:"productName"`
	ProductId      int     `json:"productId"`
	ProductCode    string  `json:"productCode"`
	CategoryId     int     `json:"categoryId"`
	CategoryName   string  `json:"categoryName"`
	GameName       string  `json:"gameName"`
	GameCode       string  `json:"gameCode"`
	TxnId          string  `json:"txnId"`
	RoundId        string  `json:"roundId"`
	EventType      int     `json:"eventType"`
	EventName      string  `json:"eventName"`
	TxnStatus      string  `json:"txnStatus"`
	TxnRemark      string  `json:"txnRemark"`
	ResultInfo     string  `json:"resultInfo"`
	Amount         float64 `json:"amount"`
	Turnover       float64 `json:"turnover"`
	IsEndRound     bool    `json:"isEndRound"`
	RequestUid     string  `json:"requestUid"`
	RequestTime    string  `json:"requestTime"`
	Timestamp      string  `json:"timestamp"`
}

type RollbackRequest struct {
	OperatorToken  string  `json:"operatorToken"`
	SeamlessKey    string  `json:"seamlessKey"`
	AgentUsername  string  `json:"agentUsername"`
	PlayerUsername string  `json:"playerUsername"`
	CurrencyCode   string  `json:"currencyCode"`
	ProductName    string  `json:"productName"`
	ProductId      int     `json:"productId"`
	ProductCode    string  `json:"productCode"`
	CategoryId     int     `json:"categoryId"`
	CategoryName   string  `json:"categoryName"`
	GameName       string  `json:"gameName"`
	GameCode       string  `json:"gameCode"`
	TxnId          string  `json:"txnId"`
	RollbackTxnId  string  `json:"rollbackTxnId"`
	RoundId        string  `json:"roundId"`
	EventType      int     `json:"eventType"`
	EventName      string  `json:"eventName"`
	TxnRemark      string  `json:"txnRemark"`
	Amount         float64 `json:"amount"`
	RequestUid     string  `json:"requestUid"`
	RequestTime    string  `json:"requestTime"`
	Timestamp      string  `json:"timestamp"`
}

type RewardRequest struct {
	OperatorToken  string  `json:"operatorToken"`
	SeamlessKey    string  `json:"seamlessKey"`
	PlayerUsername string  `json:"playerUsername"`
	CurrencyCode   string  `json:"currencyCode"`
	ProductName    string  `json:"productName"`
	ProductId      int     `json:"productId"`
	ProductCode    string  `json:"productCode"`
	CategoryId     int     `json:"categoryId"`
	CategoryName   string  `json:"categoryName"`
	EventType      int     `json:"eventType"`
	EventName      string  `json:"eventName"`
	EventDetail    string  `json:"eventDetail"`
	TxnId          string  `json:"txnId"`
	Amount         float64 `json:"amount"`
	TxnStatus      string  `json:"txnStatus"`
	TxnRemark      string  `json:"txnRemark"`
	RequestUid     string  `json:"requestUid"`
	RequestTime    string  `json:"requestTime"`
	Timestamp      string  `json:"timestamp"`
}

func BalanceProvider(c *fiber.Ctx) error {
	fmt.Println("=============== BalanceProvider =================")
	// Parse JSON body into BalanceRequest struct
	var req BalanceRequest
	if err := c.BodyParser(&req); err != nil {
		fmt.Println("Invalid request")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code": -1,
			"msg":  "Invalid request",
		})
	}

	// balance := 1000 // Replace with actual balance logic
	fmt.Println("data")
	data, err := getBalanceServer(req.PlayerUsername)
	responseTime := time.Now().Format("2006-01-02 15:04:05")
	// fmt.Println(data)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Error retrieving balance:", err)

		response := fiber.Map{
			"code":         1006,
			"msg":          "Player has Insufficient Balance to Place Bet",
			"balance":      0,
			"responseTime": responseTime,
			"responseUid":  req.RequestUid,
		}

		fmt.Println("response")
		fmt.Println(response)
		return c.JSON(response)
		// return utils.SuccessResponse(c, response, "success")
	}

	// Prepare the response
	response := fiber.Map{
		"code":         0,
		"msg":          "Successful",
		"balance":      data.Data.Balance,
		"responseTime": responseTime,
		"responseUid":  req.RequestUid,
	}
	fmt.Println("response")
	fmt.Println(response)
	// Return the JSON response
	return c.JSON(response)
}

func DebitProvider(c *fiber.Ctx) error {
	fmt.Println("=============== DebitProvider =================")
	// Parse JSON body into DebitRequest struct

	body := c.Body()
	fmt.Println("Raw Body:", string(body))

	// var req TransactionRequest
	// if err := c.BodyParser(&req); err != nil {
	// 	fmt.Println("Invalid request format")
	// 	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
	// 		"code": -1,
	// 		"msg":  "Invalid request format",
	// 	})
	// }

	// // // Example balance retrieval (replace this with actual balance logic)
	// data, err := settleServer(0, -float32(req.Amount), req.PlayerUsername)
	// if err != nil {
	// 	fmt.Println("Error retrieving balance:", err)
	// 	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
	// 		"error": "Failed to retrieve balance",
	// 	})
	// }

	// var tran models.GplayTransactions
	// tran.UserID = data.Data.UserID
	// tran.AgentID = data.Data.AgentID
	// tran.Username = data.Data.Username
	// tran.CategoryId = strconv.Itoa(req.GameCode)
	// tran.CategoryName = req.CategoryName
	// tran.ProductId = req.ProductName
	// tran.ProductCode = req.ProductCode
	// tran.WalletAmountBefore = data.Data.BalanceBefore
	// tran.WalletAmountAfter = data.Data.BalanceAfter
	// tran.BetAmount = float32(req.Amount)
	// tran.PayoutAmount = 0
	// tran.RoundId = req.RoundId
	// tran.TxnId = req.TxnId
	// tran.Status = req.EventName
	// tran.GameCode = strconv.Itoa(req.GameCode)
	// tran.PlayInfo = req.GameName
	// tran.IsEndRound = false
	// tran.IsFreeSpin = req.EventDetail.IsFeature
	// tran.BuyFeature = req.EventDetail.IsFeatureBuy
	// tran.CreatedAt = time.Now()

	// if err := database.DB.Create(&tran).Error; err != nil {
	// 	fmt.Println("pg100")
	// 	fmt.Println(err)
	// 	return err
	// }

	// currentBalance, err := getBalanceServer(req.PlayerUsername)
	// responseTime := time.Now().Format("2006-01-02 15:04:05")
	// if err != nil {
	// 	response := fiber.Map{
	// 		"code":         1006,
	// 		"msg":          "Insufficient balance",
	// 		"balance":      0,
	// 		"responseTime": responseTime,
	// 		"responseUid":  uuid.New().String(),
	// 	}
	// 	fmt.Println(response)
	// 	return c.JSON(response)
	// }

	// if currentBalance.Data.Balance < float32(req.Amount) {
	// 	// Log insufficient balance
	// 	fmt.Println("Insufficient balance for debit request:", req.PlayerUsername)

	// 	// Prepare and return the insufficient balance response
	// 	response := fiber.Map{
	// 		"code":         1006,
	// 		"msg":          "Insufficient balance",
	// 		"balance":      currentBalance.Data.Balance,
	// 		"responseTime": responseTime,
	// 		"responseUid":  uuid.New().String(),
	// 	}
	// 	fmt.Println(response)
	// 	return c.JSON(response)
	// }
	// fmt.Println("update")
	// // Deduct the requested amount from balance
	// updatedBalance := currentBalance.Data.Balance - float32(req.Amount)

	// // Log successful debit transaction
	// fmt.Printf("Debit successful for %s, amount: %.2f, new balance: %.2f\n", req.PlayerUsername, req.Amount, updatedBalance)

	// // Prepare the success response
	// response := fiber.Map{
	// 	"code":         0,
	// 	"msg":          "Debit successful",
	// 	"balance":      updatedBalance,
	// 	"responseTime": responseTime,
	// 	"responseUid":  uuid.New().String(),
	// }

	// fmt.Println(response)

	// // Return the success response with the updated balance
	return c.JSON("response")
}

func CreditProvider(c *fiber.Ctx) error {
	// Parse JSON body into CreditRequest struct
	var req CreditRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code": -1,
			"msg":  "Invalid request format",
		})
	}

	// Example balance retrieval and credit processing (replace with actual logic)
	currentBalance := 1000.0 // Replace with actual balance retrieval logic
	updatedBalance := currentBalance + req.Amount

	// Log successful credit transaction
	fmt.Printf("Credit successful for %s, amount: %.2f, new balance: %.2f\n", req.PlayerUsername, req.Amount, updatedBalance)

	// Prepare the response with the updated balance
	response := fiber.Map{
		"code":         0,
		"msg":          "Credit successful",
		"balance":      updatedBalance,
		"responseTime": time.Now().Format("2006-01-02 15:04:05"),
		"responseUid":  req.RequestUid,
	}

	// Return the success response with the updated balance
	return utils.SuccessResponse(c, response, "success")
}

func RollbackProvider(c *fiber.Ctx) error {
	// Parse JSON body into RollbackRequest struct
	var req RollbackRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code": -1,
			"msg":  "Invalid request format",
		})
	}

	// Example balance retrieval and rollback processing (replace with actual logic)
	currentBalance := 1000.0                      // Replace with actual balance retrieval logic
	updatedBalance := currentBalance + req.Amount // Rollback typically reverses the original amount

	// Log successful rollback transaction
	fmt.Printf("Rollback successful for %s, amount: %.2f, new balance: %.2f\n", req.PlayerUsername, req.Amount, updatedBalance)

	// Prepare the response with the updated balance
	response := fiber.Map{
		"code":         0,
		"msg":          "Rollback successful",
		"balance":      updatedBalance,
		"responseTime": time.Now().Format("2006-01-02 15:04:05"),
		"responseUid":  req.RequestUid,
	}

	// Return the success response with the updated balance
	return utils.SuccessResponse(c, response, "success")
}

func RewardProvider(c *fiber.Ctx) error {
	// Parse JSON body into RewardRequest struct
	var req RewardRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code": -1,
			"msg":  "Invalid request format",
		})
	}

	// Example balance retrieval and reward processing (replace with actual logic)
	currentBalance := 1000.0                      // Replace with actual balance retrieval logic
	updatedBalance := currentBalance + req.Amount // Add reward amount to current balance

	// Log successful reward transaction
	fmt.Printf("Reward successful for %s, amount: %.2f, new balance: %.2f\n", req.PlayerUsername, req.Amount, updatedBalance)

	// Prepare the response with the updated balance
	response := fiber.Map{
		"code":         0,
		"msg":          "Reward successful",
		"balance":      updatedBalance,
		"responseTime": time.Now().Format("2006-01-02 15:04:05"),
		"responseUid":  req.RequestUid,
	}

	// Return the success response with the updated balance
	return utils.SuccessResponse(c, response, "success")
}
