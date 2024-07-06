package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"game_services/app/models"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

const privateURLPG100 = "https://agent-api.pgf-asw0uz.com"
const apiKey = "OWJxTzlTNzdCRzpWWXVjZ200emhjcGFiTnZ3YzlTNWR3YWhXWk1HMmNpOQ=="

const apiKeyBankend = "BKw7jpQd8SOv7LuqPFq6MgQ4A1TflW4Ls"
const urlBankend = "https://backend.scbbbb.com/game-services"

type BodyLoginPG struct {
	Username     string `json:"username"`
	GameCode     string `json:"gameCode"`
	SessionToken string `json:"sessionToken"`
	Language     string `json:"language"`
}

func CheckBalancePG(c *fiber.Ctx) error {
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
	data, err := getBalanceServerPG(body.Username)

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

	return c.JSON(body)
}

func getBalanceServerPG(username string) (models.ResponseData, error) {
	url := fmt.Sprintf("%s/getBalance", urlBankend)
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
	fmt.Println(responseMap)
	return responseMap, nil
}

func SettleBetsPG(c *fiber.Ctx) error {
	var body models.SettleCheckResponse
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request payload",
			"error":   err.Error(),
		})
	}

	// find user
	data, err := settleServer(body)
	if err != nil {
		fmt.Println("Error retrieving balance:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve balance",
		})
	}
	fmt.Println("data = ", data)
	// err = database.DB.Transaction(func(tx *gorm.DB) error {
	// 	var pg100 models.Pg100Transactions
	// 	pg100.UserID = data.Data.UserID
	// 	pg100.Username = data.Data.Username
	// 	pg100.AgentID = data.Data.AgentID
	// 	pg100.ProductId = body.ProductId
	// 	pg100.WalletAmountBefore = data.Data.BalanceBefore
	// 	pg100.WalletAmountAfter = data.Data.BalanceAfter
	// 	pg100.BetAmount = body.Transactions[0].BetAmount
	// 	pg100.PayoutAmount = body.Transactions[0].PayoutAmount
	// 	pg100.RoundId = body.Transactions[0].RoundID
	// 	pg100.TxnId = body.Transactions[0].TxnID
	// 	pg100.Status = body.Transactions[0].Status
	// 	pg100.GameCode = body.Transactions[0].GameCode
	// 	pg100.GameId = body.Transactions[0].GameCode
	// 	pg100.PlayInfo = body.Transactions[0].PlayInfo
	// 	pg100.IsEndRound = body.Transactions[0].IsEndRound
	// 	pg100.CreatedAt = time.Now()

	// 	var winLoss = body.Transactions[0].PayoutAmount - body.Transactions[0].BetAmount
	// 	var status = ""
	// 	if winLoss > 0 {
	// 		status = "WIN"
	// 	} else if winLoss == 0 {
	// 		status = "EQ"
	// 	} else {
	// 		status = "LOSS"
	// 	}

	// 	var report models.Reports
	// 	report.UserID = data.Data.UserID
	// 	report.Username = data.Data.Username
	// 	report.AgentID = data.Data.AgentID
	// 	report.RoundId = body.Transactions[0].RoundID
	// 	report.ProductId = body.ProductId
	// 	report.ProductName = body.ProductId
	// 	report.GameId = body.Transactions[0].GameCode
	// 	report.GameName = body.Transactions[0].GameCode
	// 	report.WalletAmountBefore = data.Data.BalanceBefore
	// 	report.WalletAmountAfter = data.Data.BalanceAfter
	// 	report.BetAmount = body.Transactions[0].BetAmount
	// 	report.BetResult = body.Transactions[0].PayoutAmount
	// 	report.BetWinloss = winLoss
	// 	report.Status = status
	// 	report.IP = utils.GetIP()
	// 	report.Description = ""
	// 	report.CreatedAt = time.Now()

	// 	if err := tx.Create(&pg100).Error; err != nil {
	// 		fmt.Println("pg100")
	// 		fmt.Println(err)
	// 		return err
	// 	}

	// 	if err := tx.Create(&report).Error; err != nil {
	// 		fmt.Println("report")
	// 		fmt.Println(err)
	// 		return err
	// 	}

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

	return c.JSON(resq)
}

func settleServer(data models.SettleCheckResponse) (models.ResponseDataSettle, error) {
	url := fmt.Sprintf("%s/settleGame", urlBankend)
	fmt.Println("amount = ", data.Transactions[0].PayoutAmount-data.Transactions[0].BetAmount)
	reqBody, err := json.Marshal(map[string]interface{}{
		"username":  data.Username,
		"betsettle": data.Transactions[0].PayoutAmount - data.Transactions[0].BetAmount,
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
	url := fmt.Sprintf("%s/seamless/api/v2/login", privateURLPG100)

	fmt.Println(data.Username)
	// Marshal the data to JSON
	now := time.Now()
	sec := now.Unix()
	secStr := strconv.FormatInt(sec, 10)
	fmt.Println(secStr)
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
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Decode the response body into a JSON map
	var responseMap map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&responseMap); err != nil {
		return nil, fmt.Errorf("failed to decode response body: %v", err)
	}

	return responseMap, nil
}
