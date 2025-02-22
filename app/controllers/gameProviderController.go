package controllers

import (
	"encoding/json"
	"fmt"
	"game_services/app/database"
	"game_services/app/models"
	"game_services/app/utils"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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
	IsEndRound                   bool    `json:"isEndRound"`
	JackpotRtpContributionAmount float64 `json:"jackpotRtpContributionAmount"`
	JackpotWinAmount             float64 `json:"jackpotWinAmount"`
}

type DebitCreditRequest struct {
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
	GameCode       string  `json:"gameCode"` // Leave this as a string if "95" is a string
	TxnId          string  `json:"txnId"`
	RoundId        string  `json:"roundId"`
	EventType      int     `json:"eventType"`
	EventName      string  `json:"eventName"`
	EventDetail    string  `json:"eventDetail"` // Parsed separately if needed
	Amount         float64 `json:"amount"`
	RequestUid     string  `json:"requestUid"`
	RequestTime    string  `json:"requestTime"`
	Timestamp      int64   `json:"timestamp"` // Use int64 for Unix timestamps
	IsRefund       bool    `json:"isRefund"`
}

var sumGplay struct {
	SumBetAmount    float32 `json:"sum_bet_amount"`
	SumPayoutAmount float32 `json:"sum_payout_amount"`
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

func CheckProvider(nameBody string, operatorBody string, secretBody string) bool {
	// name := os.Getenv("agentUsername")
	operator := os.Getenv("operatorToken")
	secret := os.Getenv("seamlessKey")

	// fmt.Println("nameBody = ", nameBody)
	// fmt.Println("name = ", name)
	fmt.Println("operatorBody = ", operatorBody)
	fmt.Println("operator = ", operator)
	fmt.Println("secretBody = ", secretBody)
	fmt.Println("secret = ", secret)
	if operatorBody != operator || secretBody != secret {
		return true
	}
	return false
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
	fmt.Println(req)

	responseTime := time.Now().Format("2006-01-02 15:04:05")

	if req.OperatorToken == "" || req.SeamlessKey == "" || req.CurrencyCode == "" {
		fmt.Println("OperatorToken is empty")
		response := fiber.Map{
			"code":         9999,
			"msg":          "OperatorToken is required",
			"balance":      0,
			"responseTime": responseTime,
			"responseUid":  req.RequestUid,
		}
		fmt.Println(response)
		return c.JSON(response)
	}

	if CheckProvider(req.AgentUsername, req.OperatorToken, req.SeamlessKey) {
		response := fiber.Map{
			"code":         1004,
			"msg":          "Player has Insufficient Balance to Place Bet",
			"balance":      0,
			"responseTime": responseTime,
			"responseUid":  req.RequestUid,
		}
		fmt.Println(response)
		return c.JSON(response)
	}

	// balance := 1000 // Replace with actual balance logic
	fmt.Println("data")
	data, err := getBalanceServer(req.PlayerUsername)
	fmt.Println(data)
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

	var response fiber.Map
	if data.Data.Status == "10001" {
		response = fiber.Map{
			"code":         1000,
			"msg":          data.Message,
			"balance":      data.Data.Balance,
			"responseTime": responseTime,
			"responseUid":  req.RequestUid,
		}
	} else {
		// Prepare the response
		response = fiber.Map{
			"code":         0,
			"msg":          "Successful",
			"balance":      data.Data.Balance,
			"responseTime": responseTime,
			"responseUid":  req.RequestUid,
		}
	}

	fmt.Println("response")
	fmt.Println(response)
	// Return the JSON response
	return c.JSON(response)
}

func DebitProvider(c *fiber.Ctx) error {
	fmt.Println("=============== DebitProvider =================")

	// อ่านและพิมพ์ body สำหรับตรวจสอบ
	body := c.Body()
	fmt.Println("Raw Body:", string(body))

	// พาร์ส JSON body เป็น struct DebitCreditRequest
	var req DebitCreditRequest
	if err := c.BodyParser(&req); err != nil {
		fmt.Println("Invalid request format")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code": -1,
			"msg":  "Invalid request format",
		})
	}
	fmt.Println("Parsed Request:", req)

	responseTime := time.Now().Format("2006-01-02 15:04:05")

	if req.OperatorToken == "" || req.SeamlessKey == "" || req.CurrencyCode == "" {
		fmt.Println("OperatorToken is empty")
		response := fiber.Map{
			"code":         9999,
			"msg":          "OperatorToken is required",
			"balance":      0,
			"responseTime": responseTime,
			"responseUid":  req.RequestUid,
		}
		return c.JSON(response)
	}

	if req.CurrencyCode != os.Getenv("currencyCode") {
		fmt.Println("OperatorToken is empty")
		response := fiber.Map{
			"code":         1002,
			"msg":          "Invalid currency",
			"balance":      0,
			"responseTime": responseTime,
			"responseUid":  req.RequestUid,
		}
		return c.JSON(response)
	}

	// ตรวจสอบความซ้ำของ txnId
	var existingTxn models.GplayTransactions
	if err := database.DB.Where("txn_id = ?", req.TxnId).First(&existingTxn).Error; err == nil {
		fmt.Println("Duplicate txnId detected:", req.TxnId)
		return c.JSON(fiber.Map{
			"code":         2001,
			"msg":          "Duplicate txnId",
			"balance":      0,
			"responseTime": responseTime,
			"responseUid":  req.RequestUid,
		})
	}

	if CheckProvider(req.AgentUsername, req.OperatorToken, req.SeamlessKey) {
		response := fiber.Map{
			"code":         1004,
			"msg":          "Player has Insufficient Balance to Place Bet",
			"balance":      0,
			"responseTime": responseTime,
			"responseUid":  req.RequestUid,
		}

		fmt.Println("response")
		fmt.Println(response)
		return c.JSON(response)
	}

	// พาร์ส JSON string ของ EventDetail เป็น struct EventDetail
	var eventDetail EventDetail
	if req.EventDetail != "" {
		if err := json.Unmarshal([]byte(req.EventDetail), &eventDetail); err != nil {
			fmt.Println("เกิดข้อผิดพลาดในการพาร์ส EventDetail:", err)
			// return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			// 	"code": -1,
			// 	"msg":  "รูปแบบ EventDetail ไม่ถูกต้อง",
			// })
			eventDetail.IsFeature = false
			eventDetail.IsFeatureBuy = false
		}
	} else {
		eventDetail.IsFeature = false
		eventDetail.IsFeatureBuy = false
		fmt.Println("EventDetail เป็นค่าว่าง ข้ามการพาร์ส")
	}

	// เริ่มต้น transaction
	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// ตั้งค่า amount ให้เป็นลบสำหรับการหักเงิน
	amountSettle := -float32(req.Amount)
	fmt.Println("amountSettle =", amountSettle)

	// เรียกฟังก์ชัน settleServer เพื่อดึงข้อมูลยอดเงิน
	data, err := settleServer(amountSettle, float32(req.Amount), req.PlayerUsername, req.IsRefund)
	fmt.Println(data)
	if err != nil {
		tx.Rollback() // ยกเลิก transaction หากเกิดข้อผิดพลาด
		fmt.Println("Error retrieving balance:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve balance",
		})
	}

	// เพิ่มรายการธุรกรรมใหม่ใน GplayTransactions
	var tran models.GplayTransactions
	tran.UserID = data.Data.UserID
	tran.AgentID = data.Data.AgentID
	tran.Username = data.Data.Username
	tran.KeyDeposit = data.Data.KeyDeposit
	tran.CategoryId = req.GameCode
	tran.CategoryName = req.CategoryName
	tran.ProductId = req.ProductName
	tran.ProductCode = req.ProductCode
	tran.WalletAmountBefore = data.Data.BalanceBefore
	tran.WalletAmountAfter = data.Data.BalanceAfter
	tran.BetAmount = float32(req.Amount)
	tran.PayoutAmount = 0
	tran.RoundId = req.RoundId
	tran.TxnId = req.TxnId
	tran.Status = req.EventName
	tran.GameCode = req.GameCode
	tran.PlayInfo = req.GameName
	tran.IsEndRound = false
	tran.IsFreeSpin = eventDetail.IsFeature
	tran.BuyFeature = eventDetail.IsFeatureBuy
	tran.CreatedAt = time.Now()

	// บันทึกธุรกรรมในตาราง GplayTransactions ภายใต้ transaction
	if err := tx.Create(&tran).Error; err != nil {
		tx.Rollback() // ยกเลิก transaction หากเกิดข้อผิดพลาด
		fmt.Println("Error saving transaction:", err)
		return err
	}

	// ยืนยันการทำงานของ transaction (commit)
	if err := tx.Commit().Error; err != nil {
		fmt.Println("Error committing transaction:", err)
		return err
	}

	// ส่งข้อมูลตอบกลับ
	// response := fiber.Map{
	// 	"code":         0,
	// 	"msg":          "Debit successful",
	// 	"balance":      data.Data.BalanceAfter,
	// 	"responseTime": responseTime,
	// 	"responseUid":  uuid.New().String(),
	// }

	var response fiber.Map
	if data.Data.Status == "10001" {
		response = fiber.Map{
			"code":         1000,
			"msg":          data.Message,
			"balance":      0,
			"responseTime": responseTime,
			"responseUid":  req.RequestUid,
		}
	} else if data.Data.Status == "10002" {
		response = fiber.Map{
			"code":         1006,
			"msg":          "Player has Insufficient Balance to Place Bet",
			"balance":      0,
			"responseTime": responseTime,
			"responseUid":  req.RequestUid,
		}
	} else {
		// Prepare the response
		response = fiber.Map{
			"code":         0,
			"msg":          "Successful",
			"balance":      data.Data.BalanceAfter,
			"responseTime": responseTime,
			"responseUid":  req.RequestUid,
		}
	}

	fmt.Println("Response:", response)

	return c.JSON(response)
}

func CreditProvider(c *fiber.Ctx) error {
	fmt.Println("=============== CreditProvider =================")

	// อ่านและพิมพ์ body สำหรับตรวจสอบ
	body := c.Body()
	fmt.Println("Raw Body:", string(body))

	// พาร์ส JSON body เป็น struct DebitCreditRequest
	var req DebitCreditRequest
	if err := c.BodyParser(&req); err != nil {
		fmt.Println("Invalid request format")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code": -1,
			"msg":  "Invalid request format",
		})
	}
	fmt.Println("Parsed Request:", req)

	responseTime := time.Now().Format("2006-01-02 15:04:05")
	if req.OperatorToken == "" || req.SeamlessKey == "" || req.CurrencyCode == "" {
		fmt.Println("OperatorToken is empty")
		response := fiber.Map{
			"code":         9999,
			"msg":          "OperatorToken is required",
			"balance":      0,
			"responseTime": responseTime,
			"responseUid":  req.RequestUid,
		}
		return c.JSON(response)
	}

	// ตรวจสอบความซ้ำของ txnId
	var existingTxn models.GplayTransactions
	if err := database.DB.Where("txn_id = ?", req.TxnId).First(&existingTxn).Error; err == nil {
		fmt.Println("Duplicate txnId detected:", req.TxnId)
		return c.JSON(fiber.Map{
			"code":         2001,
			"msg":          "Duplicate txnId",
			"balance":      0,
			"responseTime": responseTime,
			"responseUid":  req.RequestUid,
		})
	}

	if CheckProvider(req.AgentUsername, req.OperatorToken, req.SeamlessKey) {
		response := fiber.Map{
			"code":         1004,
			"msg":          "Player has Insufficient Balance to Place Bet",
			"balance":      0,
			"responseTime": responseTime,
			"responseUid":  req.RequestUid,
		}

		fmt.Println("response")
		fmt.Println(response)
		return c.JSON(response)
	}

	// พาร์ส JSON string ของ EventDetail เป็น struct EventDetail
	var eventDetail EventDetail
	if req.EventDetail != "" {
		if err := json.Unmarshal([]byte(req.EventDetail), &eventDetail); err != nil {
			fmt.Println("เกิดข้อผิดพลาดในการพาร์ส EventDetail:", err)
			// return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			// 	"code": -1,
			// 	"msg":  "รูปแบบ EventDetail ไม่ถูกต้อง",
			// })

			// eventDetail.IsFeature = false
			// eventDetail.IsFeatureBuy = false
		}
	} else {
		// eventDetail.IsFeature = false
		// eventDetail.IsFeatureBuy = false
		fmt.Println("EventDetail เป็นค่าว่าง ข้ามการพาร์ส")
	}

	// ตั้งค่า amount ให้เป็นบวกสำหรับการเติมเงิน
	amountSettle := float32(req.Amount)
	fmt.Println("amountSettle =", amountSettle)

	// เรียกฟังก์ชัน settleServer เพื่อทำการเติมเงิน
	data, err := settleServer(amountSettle, float32(req.Amount), req.PlayerUsername, req.IsRefund)
	if err != nil {
		fmt.Println("Error retrieving balance:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve balance",
		})
	}

	// เพิ่มรายการธุรกรรมใหม่ใน GplayTransactions
	var tran models.GplayTransactions
	tran.UserID = data.Data.UserID
	tran.AgentID = data.Data.AgentID
	tran.Username = data.Data.Username
	tran.KeyDeposit = data.Data.KeyDeposit
	tran.CategoryId = req.GameCode
	tran.CategoryName = req.CategoryName
	tran.ProductId = req.ProductName
	tran.ProductCode = req.ProductCode
	tran.WalletAmountBefore = data.Data.BalanceBefore
	tran.WalletAmountAfter = data.Data.BalanceAfter
	tran.BetAmount = 0                      // เพราะเป็นการเติมเงิน ไม่ใช่การเดิมพัน
	tran.PayoutAmount = float32(req.Amount) // จำนวนเงินที่เติมเข้าระบบ
	tran.RoundId = req.RoundId
	tran.TxnId = req.TxnId
	tran.Status = req.EventName
	tran.GameCode = req.GameCode
	tran.PlayInfo = req.GameName
	tran.IsEndRound = eventDetail.IsEndRound
	tran.IsFreeSpin = eventDetail.IsFeature
	tran.BuyFeature = eventDetail.IsFeatureBuy
	tran.CreatedAt = time.Now()

	// บันทึกธุรกรรมในตาราง GplayTransactions ภายใต้ transaction
	if err := database.DB.Create(&tran).Error; err != nil {
		fmt.Println("Error saving transaction:", err)
		return err
	}

	if eventDetail.IsEndRound {

		// คำนวณยอดรวมของ Bet ใน round เดียวกันจากธุรกรรมที่เป็น credit
		parts := strings.Split(req.TxnId, "-")
		fmt.Println(parts[1])
		// var sumPayoutAmount float32
		// var sum
		if err := database.DB.Model(&models.GplayTransactions{}).
			Where("txn_id LIKE ?", "%"+parts[1]+"%").
			Select("COALESCE(SUM(bet_amount), 0) AS sum_bet_amount, COALESCE(SUM(payout_amount), 0) AS sum_payout_amount").
			Scan(&sumGplay).Error; err != nil {
			fmt.Println("Error calculating sum:", err)
			return err
		}

		// คำนวณยอดชนะ/แพ้ และสถานะ
		var winLoss = sumGplay.SumPayoutAmount - sumGplay.SumBetAmount
		var status = ""
		if winLoss > 0 {
			status = "WIN"
		} else if winLoss == 0 {
			status = "EQ"
		} else {
			status = "LOSS"
		}
		fmt.Printf("Total Bet Amount: %.2f, Total Payout Amount: %.2f\n", sumGplay.SumBetAmount, sumGplay.SumPayoutAmount)
		// เพิ่มรายการใน Reports ภายใต้ transaction
		fmt.Println(req.CategoryName)
		var report models.Reports
		report.UserID = data.Data.UserID
		report.Username = data.Data.Username
		report.KeyDeposit = data.Data.KeyDeposit
		report.CategoryName = req.CategoryName
		report.AgentID = data.Data.AgentID
		report.RoundId = req.RoundId
		report.RoundCheck = parts[1]
		report.ProductId = req.ProductName
		report.ProductName = req.ProductName
		report.GameId = req.GameCode
		report.GameName = req.GameName
		report.WalletAmountBefore = data.Data.BalanceBefore + (sumGplay.SumBetAmount - sumGplay.SumPayoutAmount)
		report.WalletAmountAfter = data.Data.BalanceAfter
		report.BetAmount = sumGplay.SumBetAmount
		report.BetResult = sumGplay.SumPayoutAmount
		report.BetWinloss = winLoss
		report.Status = status
		report.IP = utils.GetIP()
		report.Description = ""
		report.CreatedAt = time.Now()

		// บันทึกข้อมูลรายงานลงฐานข้อมูล
		if err := database.DB.Create(&report).Error; err != nil {
			fmt.Println("Error saving report:", err)
			return err
		}
	} else {
		if err := database.DB.Model(&models.GplayTransactions{}).
			Where("round_id = ?", req.RoundId).
			Select("COALESCE(SUM(bet_amount), 0) AS sum_bet_amount, COALESCE(SUM(payout_amount), 0) AS sum_payout_amount").
			Scan(&sumGplay).Error; err != nil {
			fmt.Println("Error calculating sum:", err)
			return err
		}

		var winLoss = sumGplay.SumPayoutAmount - sumGplay.SumBetAmount
		var status = ""
		if winLoss > 0 {
			status = "WIN"
		} else if winLoss == 0 {
			status = "EQ"
		} else {
			status = "LOSS"
		}
		if req.CategoryName == "Live Casino" {
			var report models.Reports
			report.UserID = data.Data.UserID
			report.Username = data.Data.Username
			report.KeyDeposit = data.Data.KeyDeposit
			report.CategoryName = req.CategoryName
			report.AgentID = data.Data.AgentID
			report.RoundId = req.RoundId
			report.RoundCheck = req.RoundId
			report.ProductId = req.ProductName
			report.ProductName = req.ProductName
			report.GameId = req.GameCode
			report.GameName = req.GameName
			report.WalletAmountBefore = data.Data.BalanceBefore + (sumGplay.SumBetAmount - sumGplay.SumPayoutAmount)
			report.WalletAmountAfter = data.Data.BalanceAfter
			report.BetAmount = sumGplay.SumBetAmount
			report.BetResult = sumGplay.SumPayoutAmount
			report.BetWinloss = winLoss
			report.Status = status
			report.IP = utils.GetIP()
			report.Description = ""
			report.CreatedAt = time.Now()
			if err := database.DB.Create(&report).Error; err != nil {
				fmt.Println("Error saving report:", err)
				return err
			}
		} else if req.ProductName != "PGSOFT" {
			var report models.Reports
			report.UserID = data.Data.UserID
			report.Username = data.Data.Username
			report.KeyDeposit = data.Data.KeyDeposit
			report.CategoryName = req.CategoryName
			report.AgentID = data.Data.AgentID
			report.RoundId = req.RoundId
			report.RoundCheck = req.RoundId
			report.ProductId = req.ProductName
			report.ProductName = req.ProductName
			report.GameId = req.GameCode
			report.GameName = req.GameName
			report.WalletAmountBefore = data.Data.BalanceBefore + (sumGplay.SumBetAmount - sumGplay.SumPayoutAmount)
			report.WalletAmountAfter = data.Data.BalanceAfter
			report.BetAmount = sumGplay.SumBetAmount
			report.BetResult = sumGplay.SumPayoutAmount
			report.BetWinloss = winLoss
			report.Status = status
			report.IP = utils.GetIP()
			report.Description = ""
			report.CreatedAt = time.Now()
			if err := database.DB.Create(&report).Error; err != nil {
				fmt.Println("Error saving report:", err)
				return err
			}
		}
	}

	// สร้าง response เวลาปัจจุบัน

	// Log สำหรับการทำรายการเติมเงินสำเร็จ
	fmt.Printf("Credit successful for %s, amount: %.2f, new balance: %.2f\n", req.PlayerUsername, req.Amount, data.Data.BalanceAfter)

	// สร้างข้อมูล response และส่งกลับ
	response := fiber.Map{
		"code":         0,
		"msg":          "Credit successful",
		"balance":      data.Data.BalanceAfter,
		"responseTime": responseTime,
		"responseUid":  uuid.New().String(),
	}

	fmt.Println("Response:", response)

	// ส่ง response กลับในรูปแบบ JSON
	return c.JSON(response)
}

func RollbackProvider(c *fiber.Ctx) error {
	fmt.Println("=============== RollbackProvider =================")

	// อ่านและพิมพ์ body สำหรับตรวจสอบ
	body := c.Body()
	fmt.Println("Raw Body:", string(body))

	// พาร์ส JSON body เป็น struct DebitCreditRequest
	var req DebitCreditRequest
	if err := c.BodyParser(&req); err != nil {
		fmt.Println("Invalid request format")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code": -1,
			"msg":  "Invalid request format",
		})
	}
	fmt.Println("Parsed Request:", req)

	responseTime := time.Now().Format("2006-01-02 15:04:05")
	if req.OperatorToken == "" || req.SeamlessKey == "" || req.CurrencyCode == "" {
		fmt.Println("OperatorToken is empty")
		response := fiber.Map{
			"code":         9999,
			"msg":          "OperatorToken is required",
			"balance":      0,
			"responseTime": responseTime,
			"responseUid":  req.RequestUid,
		}
		return c.JSON(response)
	}

	if CheckProvider(req.AgentUsername, req.OperatorToken, req.SeamlessKey) {
		response := fiber.Map{
			"code":         1004,
			"msg":          "Player has Insufficient Balance to Place Bet",
			"balance":      0,
			"responseTime": responseTime,
			"responseUid":  req.RequestUid,
		}

		fmt.Println("response")
		fmt.Println(response)
		return c.JSON(response)
	}

	// พาร์ส JSON string ของ EventDetail เป็น struct EventDetail
	var eventDetail EventDetail
	if req.EventDetail != "" {
		if err := json.Unmarshal([]byte(req.EventDetail), &eventDetail); err != nil {
			fmt.Println("เกิดข้อผิดพลาดในการพาร์ส EventDetail:", err)
			// return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			// 	"code": -1,
			// 	"msg":  "รูปแบบ EventDetail ไม่ถูกต้อง",
			// })
		}
	} else {
		fmt.Println("EventDetail เป็นค่าว่าง ข้ามการพาร์ส")
	}

	// ตั้งค่า amount ให้เป็นบวกสำหรับการเติมเงิน
	amountSettle := float32(req.Amount)
	fmt.Println("amountSettle =", amountSettle)

	// เรียกฟังก์ชัน settleServer เพื่อทำการเติมเงิน
	data, err := settleServer(amountSettle, float32(req.Amount), req.PlayerUsername, false)
	if err != nil {
		fmt.Println("Error retrieving balance:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve balance",
		})
	}

	// เพิ่มรายการธุรกรรมใหม่ใน GplayTransactions
	var tran models.GplayTransactions
	tran.UserID = data.Data.UserID
	tran.AgentID = data.Data.AgentID
	tran.Username = data.Data.Username
	tran.KeyDeposit = data.Data.KeyDeposit
	tran.CategoryId = req.GameCode
	tran.CategoryName = req.CategoryName
	tran.ProductId = req.ProductName
	tran.ProductCode = req.ProductCode
	tran.WalletAmountBefore = data.Data.BalanceBefore
	tran.WalletAmountAfter = data.Data.BalanceAfter
	tran.BetAmount = 0                      // เพราะเป็นการเติมเงิน ไม่ใช่การเดิมพัน
	tran.PayoutAmount = float32(req.Amount) // จำนวนเงินที่เติมเข้าระบบ
	tran.RoundId = req.RoundId
	tran.TxnId = req.TxnId
	tran.Status = req.EventName
	tran.GameCode = req.GameCode
	tran.PlayInfo = req.GameName
	tran.IsEndRound = eventDetail.IsEndRound
	tran.IsFreeSpin = eventDetail.IsFeature
	tran.BuyFeature = eventDetail.IsFeatureBuy
	tran.CreatedAt = time.Now()

	// บันทึกธุรกรรมในตาราง GplayTransactions ภายใต้ transaction
	if err := database.DB.Create(&tran).Error; err != nil {
		fmt.Println("Error saving transaction:", err)
		return err
	}

	if eventDetail.IsEndRound {

		// คำนวณยอดรวมของ Bet ใน round เดียวกันจากธุรกรรมที่เป็น credit
		parts := strings.Split(req.TxnId, "-")
		fmt.Println(parts[1])
		// var sumPayoutAmount float32
		// var sum
		if err := database.DB.Model(&models.GplayTransactions{}).
			Where("txn_id LIKE ?", "%"+parts[1]+"%").
			Select("COALESCE(SUM(bet_amount), 0) AS sum_bet_amount, COALESCE(SUM(payout_amount), 0) AS sum_payout_amount").
			Scan(&sumGplay).Error; err != nil {
			fmt.Println("Error calculating sum:", err)
			return err
		}

		// คำนวณยอดชนะ/แพ้ และสถานะ
		var winLoss = sumGplay.SumPayoutAmount - sumGplay.SumBetAmount
		var status = ""
		if winLoss > 0 {
			status = "WIN"
		} else if winLoss == 0 {
			status = "EQ"
		} else {
			status = "LOSS"
		}
		fmt.Printf("Total Bet Amount: %.2f, Total Payout Amount: %.2f\n", sumGplay.SumBetAmount, sumGplay.SumPayoutAmount)
		// เพิ่มรายการใน Reports ภายใต้ transaction
		fmt.Println(data.Data)
		var report models.Reports
		report.UserID = data.Data.UserID
		report.Username = data.Data.Username
		report.KeyDeposit = data.Data.KeyDeposit
		report.CategoryName = req.CategoryName
		report.AgentID = data.Data.AgentID
		report.RoundId = req.RoundId
		report.RoundCheck = parts[1]
		report.ProductId = req.ProductName
		report.ProductName = req.ProductName
		report.GameId = req.GameCode
		report.GameName = req.GameName
		report.WalletAmountBefore = data.Data.BalanceBefore + (sumGplay.SumBetAmount - sumGplay.SumPayoutAmount)
		report.WalletAmountAfter = data.Data.BalanceAfter
		report.BetAmount = sumGplay.SumBetAmount
		report.BetResult = sumGplay.SumPayoutAmount
		report.BetWinloss = winLoss
		report.Status = status
		report.IP = utils.GetIP()
		report.Description = ""
		report.CreatedAt = time.Now()

		// บันทึกข้อมูลรายงานลงฐานข้อมูล
		if err := database.DB.Create(&report).Error; err != nil {
			fmt.Println("Error saving report:", err)
			return err
		}
	}

	// สร้าง response เวลาปัจจุบัน

	// Log สำหรับการทำรายการเติมเงินสำเร็จ
	fmt.Printf("Credit successful for %s, amount: %.2f, new balance: %.2f\n", req.PlayerUsername, req.Amount, data.Data.BalanceAfter)

	// สร้างข้อมูล response และส่งกลับ
	response := fiber.Map{
		"code":         0,
		"msg":          "Credit successful",
		"balance":      data.Data.BalanceAfter,
		"responseTime": responseTime,
		"responseUid":  uuid.New().String(),
	}

	fmt.Println("Response:", response)

	// ส่ง response กลับในรูปแบบ JSON
	return c.JSON(response)
}

func RewardProvider(c *fiber.Ctx) error {
	fmt.Println("=============== RewardProvider =================")

	// อ่านและพิมพ์ body สำหรับตรวจสอบ
	body := c.Body()
	fmt.Println("Raw Body:", string(body))

	// พาร์ส JSON body เป็น struct DebitCreditRequest
	var req DebitCreditRequest
	if err := c.BodyParser(&req); err != nil {
		fmt.Println("Invalid request format")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code": -1,
			"msg":  "Invalid request format",
		})
	}
	fmt.Println("Parsed Request:", req)

	responseTime := time.Now().Format("2006-01-02 15:04:05")
	if req.OperatorToken == "" || req.SeamlessKey == "" || req.CurrencyCode == "" {
		fmt.Println("OperatorToken is empty")
		response := fiber.Map{
			"code":         9999,
			"msg":          "OperatorToken is required",
			"balance":      0,
			"responseTime": responseTime,
			"responseUid":  req.RequestUid,
		}
		return c.JSON(response)
	}

	// ตรวจสอบความซ้ำของ txnId
	var existingTxn models.GplayTransactions
	if err := database.DB.Where("txn_id = ?", req.TxnId).First(&existingTxn).Error; err == nil {
		fmt.Println("Duplicate txnId detected:", req.TxnId)
		return c.JSON(fiber.Map{
			"code":         2001,
			"msg":          "Duplicate txnId",
			"balance":      0,
			"responseTime": responseTime,
			"responseUid":  req.RequestUid,
		})
	}

	if CheckProvider(req.AgentUsername, req.OperatorToken, req.SeamlessKey) {
		response := fiber.Map{
			"code":         1004,
			"msg":          "Player has Insufficient Balance to Place Bet",
			"balance":      0,
			"responseTime": responseTime,
			"responseUid":  req.RequestUid,
		}

		fmt.Println("response")
		fmt.Println(response)
		return c.JSON(response)
	}

	// พาร์ส JSON string ของ EventDetail เป็น struct EventDetail
	var eventDetail EventDetail
	if req.EventDetail != "" {
		if err := json.Unmarshal([]byte(req.EventDetail), &eventDetail); err != nil {
			fmt.Println("เกิดข้อผิดพลาดในการพาร์ส EventDetail:", err)
			// return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			// 	"code": -1,
			// 	"msg":  "รูปแบบ EventDetail ไม่ถูกต้อง",
			// })
		}
	} else {
		fmt.Println("EventDetail เป็นค่าว่าง ข้ามการพาร์ส")
	}

	// ตั้งค่า amount ให้เป็นบวกสำหรับการเติมเงิน
	amountSettle := float32(req.Amount)
	fmt.Println("amountSettle =", amountSettle)

	// เรียกฟังก์ชัน settleServer เพื่อทำการเติมเงิน
	data, err := settleServer(amountSettle, float32(req.Amount), req.PlayerUsername, false)
	if err != nil {
		fmt.Println("Error retrieving balance:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve balance",
		})
	}

	// เพิ่มรายการธุรกรรมใหม่ใน GplayTransactions
	var tran models.GplayTransactions
	tran.UserID = data.Data.UserID
	tran.AgentID = data.Data.AgentID
	tran.Username = data.Data.Username
	tran.KeyDeposit = data.Data.KeyDeposit
	tran.CategoryId = req.GameCode
	tran.CategoryName = req.CategoryName
	tran.ProductId = req.ProductName
	tran.ProductCode = req.ProductCode
	tran.WalletAmountBefore = data.Data.BalanceBefore
	tran.WalletAmountAfter = data.Data.BalanceAfter
	tran.BetAmount = 0                      // เพราะเป็นการเติมเงิน ไม่ใช่การเดิมพัน
	tran.PayoutAmount = float32(req.Amount) // จำนวนเงินที่เติมเข้าระบบ
	tran.RoundId = req.RoundId
	tran.TxnId = req.TxnId
	tran.Status = req.EventName
	tran.GameCode = req.GameCode
	tran.PlayInfo = req.GameName
	tran.IsEndRound = eventDetail.IsEndRound
	tran.IsFreeSpin = eventDetail.IsFeature
	tran.BuyFeature = eventDetail.IsFeatureBuy
	tran.CreatedAt = time.Now()

	// บันทึกธุรกรรมในตาราง GplayTransactions ภายใต้ transaction
	if err := database.DB.Create(&tran).Error; err != nil {
		fmt.Println("Error saving transaction:", err)
		return err
	}

	if eventDetail.IsEndRound {

		// คำนวณยอดรวมของ Bet ใน round เดียวกันจากธุรกรรมที่เป็น credit
		parts := strings.Split(req.TxnId, "-")
		fmt.Println(parts[1])
		// var sumPayoutAmount float32
		// var sum
		if err := database.DB.Model(&models.GplayTransactions{}).
			Where("txn_id LIKE ?", "%"+parts[1]+"%").
			Select("COALESCE(SUM(bet_amount), 0) AS sum_bet_amount, COALESCE(SUM(payout_amount), 0) AS sum_payout_amount").
			Scan(&sumGplay).Error; err != nil {
			fmt.Println("Error calculating sum:", err)
			return err
		}

		// คำนวณยอดชนะ/แพ้ และสถานะ
		var winLoss = sumGplay.SumPayoutAmount - sumGplay.SumBetAmount
		var status = ""
		if winLoss > 0 {
			status = "WIN"
		} else if winLoss == 0 {
			status = "EQ"
		} else {
			status = "LOSS"
		}
		fmt.Printf("Total Bet Amount: %.2f, Total Payout Amount: %.2f\n", sumGplay.SumBetAmount, sumGplay.SumPayoutAmount)
		// เพิ่มรายการใน Reports ภายใต้ transaction
		fmt.Println(data.Data)
		var report models.Reports
		report.UserID = data.Data.UserID
		report.Username = data.Data.Username
		report.KeyDeposit = data.Data.KeyDeposit
		report.CategoryName = req.CategoryName
		report.AgentID = data.Data.AgentID
		report.RoundId = req.RoundId
		report.RoundCheck = parts[1]
		report.ProductId = req.ProductName
		report.ProductName = req.ProductName
		report.GameId = req.GameCode
		report.GameName = req.GameName
		report.WalletAmountBefore = data.Data.BalanceBefore + (sumGplay.SumBetAmount - sumGplay.SumPayoutAmount)
		report.WalletAmountAfter = data.Data.BalanceAfter
		report.BetAmount = sumGplay.SumBetAmount
		report.BetResult = sumGplay.SumPayoutAmount
		report.BetWinloss = winLoss
		report.Status = status
		report.IP = utils.GetIP()
		report.Description = ""
		report.CreatedAt = time.Now()

		// บันทึกข้อมูลรายงานลงฐานข้อมูล
		if err := database.DB.Create(&report).Error; err != nil {
			fmt.Println("Error saving report:", err)
			return err
		}
	}

	// สร้าง response เวลาปัจจุบัน

	// Log สำหรับการทำรายการเติมเงินสำเร็จ
	fmt.Printf("Credit successful for %s, amount: %.2f, new balance: %.2f\n", req.PlayerUsername, req.Amount, data.Data.BalanceAfter)

	// สร้างข้อมูล response และส่งกลับ
	response := fiber.Map{
		"code":         0,
		"msg":          "Credit successful",
		"balance":      data.Data.BalanceAfter,
		"responseTime": responseTime,
		"responseUid":  uuid.New().String(),
	}

	fmt.Println("Response:", response)

	// ส่ง response กลับในรูปแบบ JSON
	return c.JSON(response)
}
