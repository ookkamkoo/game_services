package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"game_services/app/database"
	"game_services/app/models"
	"game_services/app/utils"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

type BodyLoginPG struct {
	Username     string          `json:"username"`
	GameCode     string          `json:"gameCode"`
	SessionToken string          `json:"sessionToken"`
	Language     string          `json:"language"`
	Setting      json.RawMessage `json:"setting"`
}

func CheckBalancePG(c *fiber.Ctx) error {
	startTime := time.Now()
	startFormatted := startTime.Format("2006-01-02 15:04:05.000")
	fmt.Println("===================== CheckBalance ===========================")
	fmt.Println("Start date and time:", startFormatted)

	var body models.BalanceCheckResponse
	if err := c.BodyParser(&body); err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}
	fmt.Println("find user")
	fmt.Println(body)
	// find user
	data, err := getBalanceServer(body.Username)

	if err != nil {
		fmt.Println(err)
		fmt.Println("Error retrieving balance:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve balance",
		})
	}
	// find user
	now := time.Now()
	timestamp := now.UnixNano() / int64(time.Millisecond)
	body.Balance = data.Data.Balance
	body.TimestampMillis = timestamp

	endTime := time.Now()
	endFormatted := endTime.Format("2006-01-02 15:04:05.000")
	fmt.Println("End date and time:", endFormatted)

	duration := endTime.Sub(startTime)
	fmt.Println("Duration:", duration)
	fmt.Println("================================================")

	return c.JSON(body)
}

func getBalanceServer(username string) (models.ResponseData, error) {
	urlBankend := os.Getenv("urlBankend")
	apiKeyBankend := os.Getenv("apiKeyBankend")

	url := fmt.Sprintf("%s/getBalance", urlBankend)
	// url := "http://127.0.0.1:3001/services-game/getBalance"
	fmt.Println(url)
	reqBody, err := json.Marshal(map[string]interface{}{
		"username": username,
	})
	fmt.Println(username)
	if err != nil {
		return models.ResponseData{}, fmt.Errorf("failed to marshal request body: %v", err)
	}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return models.ResponseData{}, fmt.Errorf("failed to create HTTP request: %v", err)
	}

	fmt.Println("key =", apiKeyBankend)
	// Set the required headers
	req.Header.Set("x-api-key", apiKeyBankend)
	req.Header.Set("Content-Type", "application/json")

	// Execute the HTTP request
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return models.ResponseData{}, fmt.Errorf("failed to send HTTP request: %v", err)
	}
	defer resp.Body.Close()
	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		return models.ResponseData{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Decode the response body into a JSON map
	var responseMap models.ResponseData
	if err := json.NewDecoder(resp.Body).Decode(&responseMap); err != nil {
		return models.ResponseData{}, fmt.Errorf("failed to decode response body: %v", err)
	}

	fmt.Println("resp = ")
	fmt.Println(responseMap)
	return responseMap, nil
}

func SettleBetsPG(c *fiber.Ctx) error {
	startTime := time.Now()
	startFormatted := startTime.Format("2006-01-02 15:04:05.000")
	fmt.Println("==================== SettleBets ============================")
	fmt.Println("Start date and time:", startFormatted)

	var body models.SettleCheckResponse
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request payload",
			"error":   err.Error(),
		})
	}
	fmt.Println("body = ", body)
	// find user
	amountSettle := body.Transactions[0].PayoutAmount - body.Transactions[0].BetAmount
	data, err := settleServer(amountSettle, body.Username, false)
	if err != nil {
		fmt.Println("Error retrieving balance:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve balance",
		})
	}
	fmt.Println("data = ", data)
	// err = database.DB.Transaction(func(tx *gorm.DB) error {
	var pg100 models.Pg100Transactions
	pg100.UserID = data.Data.UserID
	pg100.Username = data.Data.Username
	pg100.KeyDeposit = data.Data.KeyDeposit
	pg100.AgentID = data.Data.AgentID
	pg100.ProductId = body.ProductId
	pg100.WalletAmountBefore = data.Data.BalanceBefore
	pg100.WalletAmountAfter = data.Data.BalanceAfter
	pg100.BetAmount = body.Transactions[0].BetAmount
	pg100.PayoutAmount = body.Transactions[0].PayoutAmount
	pg100.RoundId = body.Transactions[0].RoundID
	pg100.TxnId = body.Transactions[0].TxnID
	pg100.Status = body.Transactions[0].Status
	pg100.GameCode = body.Transactions[0].GameCode
	pg100.GameId = body.Transactions[0].GameCode
	pg100.PlayInfo = body.Transactions[0].PlayInfo
	pg100.IsEndRound = body.Transactions[0].IsEndRound
	pg100.CreatedAt = time.Now()

	var winLoss = body.Transactions[0].PayoutAmount - body.Transactions[0].BetAmount
	var status = ""
	if winLoss > 0 {
		status = "WIN"
	} else if winLoss == 0 {
		status = "EQ"
	} else {
		status = "LOSS"
	}

	var report models.Reports
	report.UserID = data.Data.UserID
	report.Username = data.Data.Username
	report.AgentID = data.Data.AgentID
	report.KeyDeposit = data.Data.KeyDeposit
	report.CategoryName = "Pg100"
	report.RoundId = body.Transactions[0].RoundID
	report.RoundCheck = body.Transactions[0].RoundID
	report.ProductId = body.ProductId
	report.ProductName = body.ProductId
	report.GameId = body.Transactions[0].GameCode
	report.GameName = body.Transactions[0].GameCode
	report.WalletAmountBefore = data.Data.BalanceBefore
	report.WalletAmountAfter = data.Data.BalanceAfter
	report.BetAmount = body.Transactions[0].BetAmount
	report.BetResult = body.Transactions[0].PayoutAmount
	report.BetWinloss = winLoss
	report.Status = status
	report.IP = utils.GetIP()
	report.Description = ""
	report.CreatedAt = time.Now()

	fmt.Println("report = ")
	fmt.Println(report)

	if err := database.DB.Create(&pg100).Error; err != nil {
		fmt.Println("pg100")
		fmt.Println(err)
		return err
	}

	if err := database.DB.Create(&report).Error; err != nil {
		fmt.Println("report")
		fmt.Println(err)
		return err
	}

	// 	return nil
	// })

	var resq models.SettleCheckResponse
	now := time.Now()
	resq.ID = body.ID
	resq.TimestampMillis = now.UnixNano() / int64(time.Millisecond)
	resq.ProductId = body.ProductId
	resq.Currency = body.Currency
	resq.Username = body.Username
	resq.BalanceBefore = data.Data.BalanceBefore
	resq.BalanceAfter = 0

	if err != nil {
		resq.StatusCode = 50001
		return c.JSON(resq)
	}

	fmt.Println(data.Data.Status)
	statusInt, err := strconv.Atoi(data.Data.Status)
	if err != nil {
		log.Fatalf("Failed to convert status to int: %v", err)
	}

	resq.BalanceBefore = data.Data.BalanceBefore
	resq.BalanceAfter = data.Data.BalanceAfter
	resq.StatusCode = statusInt
	fmt.Println("resq = ", resq)

	endTime := time.Now()
	endFormatted := endTime.Format("2006-01-02 15:04:05.000")
	fmt.Println("End date and time:", endFormatted)

	duration := endTime.Sub(startTime)
	fmt.Println("Duration:", duration)
	fmt.Println("================================================")

	return c.JSON(resq)
}

func settleServer(amount float32, username string, is_refund bool) (models.ResponseDataSettle, error) {
	apiKeyBankend := os.Getenv("apiKeyBankend")
	urlBankend := os.Getenv("urlBankend")
	url := fmt.Sprintf("%s/settleGame", urlBankend)
	// fmt.Println("betAmount = ", betAmount)
	fmt.Println("amount = ", amount)
	reqBody, err := json.Marshal(map[string]interface{}{
		"username":  username,
		"betsettle": amount,
		"is_refund": is_refund,
	})
	if err != nil {
		return models.ResponseDataSettle{}, fmt.Errorf("failed to marshal request body: %v", err)
	}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return models.ResponseDataSettle{}, fmt.Errorf("failed to create HTTP request: %v", err)
	}
	// Set the required headers
	req.Header.Set("x-api-key", apiKeyBankend)
	req.Header.Set("Content-Type", "application/json")

	// Execute the HTTP request
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return models.ResponseDataSettle{}, fmt.Errorf("failed to send HTTP request: %v", err)
	}
	defer resp.Body.Close()
	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		return models.ResponseDataSettle{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	// Decode the response body into a JSON map
	var responseMap models.ResponseDataSettle
	if err := json.NewDecoder(resp.Body).Decode(&responseMap); err != nil {
		return models.ResponseDataSettle{}, fmt.Errorf("failed to decode response body: %v", err)
	}
	// fmt.Println(responseMap)
	return responseMap, nil
}

func PGGameList() (map[string]interface{}, error) {
	privateURLPG100 := os.Getenv("PRIVATE_URL_PG100")
	apiKey := os.Getenv("apiKey")
	url := fmt.Sprintf("%s/seamless/api/v2/games", privateURLPG100)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("x-api-key", apiKey)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)

	}

	// Decode the response body into a JSON array
	var responseMap map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&responseMap); err != nil {
		return nil, err
	}
	return responseMap, nil
}

func PGLaunchGames(data BodyLoginPG) (map[string]interface{}, error) {
	privateURLPG100 := os.Getenv("PRIVATE_URL_PG100")
	apiKey := os.Getenv("apiKey")
	url := fmt.Sprintf("%s/seamless/api/v2/login", privateURLPG100)
	fmt.Println("sssssssssss")
	// fmt.Println(data.Username)
	// Marshal the data to JSON
	now := time.Now()
	sec := now.Unix()
	secStr := strconv.FormatInt(sec, 10)
	// fmt.Println(secStr)
	reqBody, err := json.Marshal(map[string]interface{}{
		"username":     data.Username,
		"gameCode":     data.GameCode,
		"sessionToken": secStr,
		"language":     data.Language,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %v", err)
	}

	// Create a new HTTP POST request
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %v", err)
	}

	// Set the required headers
	req.Header.Set("x-api-key", apiKey)
	req.Header.Set("Content-Type", "application/json")

	// Execute the HTTP request
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send HTTP request: %v", err)
	}
	defer resp.Body.Close()

	// Check the response status code
	fmt.Println("PGLaunchGames = ", resp.StatusCode)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Decode the response body into a JSON map
	var responseMap map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&responseMap); err != nil {
		return nil, fmt.Errorf("failed to decode response body: %v", err)
	}
	fmt.Println(responseMap)
	return responseMap, nil
}

func PGSettingGame(data json.RawMessage) error {
	// url := fmt.Sprintf("%s/seamless/api/v2/setGameSetting", privateURLPG100)
	url := "https://agent-api.u17fz.com/seamless/api/v2/setGameSetting"
	apiKey := "OWJxTzlTNzdCRzpWWXVjZ200emhjcGFiTnZ3YzlTNWR3YWhXWk1HMmNpOQ=="
	fmt.Println(url)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %v", err)
	}

	// Set the required headers
	req.Header.Set("x-api-key", apiKey)
	req.Header.Set("Content-Type", "application/json")

	// Execute the HTTP request
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send HTTP request: %v", err)
	}
	defer resp.Body.Close()

	// Check the response status code
	fmt.Println("PGSettingGame = ", resp.StatusCode)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Decode the response body into a JSON map
	var responseMap map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&responseMap); err != nil {
		return fmt.Errorf("failed to decode response body: %v", err)
	}
	fmt.Println(responseMap)
	return nil
}
