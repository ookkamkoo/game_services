package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"game_services/app/database"
	"game_services/app/models"
	"game_services/app/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
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

	fmt.Println("111111111111111111111")
	// find user
	data, err := settleServer(body)

	if err != nil {
		fmt.Println("Error retrieving balance:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve balance",
		})
	}

	err = database.DB.Transaction(func(tx *gorm.DB) error {
		var pg100 models.Pg100Transactions
		pg100.UserID = data.UserID
		pg100.Username = data.Username
		pg100.AgentID = data.AgentID
		pg100.ProductId = "pg100"
		pg100.WalletAmountBefore = data.BalanceBefore
		pg100.WalletAmountAfter = data.BalanceAfter
		pg100.BetAmount = body.Transactions.BetAmount
		pg100.PayoutAmount = body.Transactions.PayoutAmount
		pg100.RoundId = body.Transactions.RoundID
		pg100.TxnId = body.Transactions.TxnID
		pg100.Status = body.Transactions.Status
		pg100.GameCode = body.Transactions.GameCode
		pg100.GameId = body.Transactions.GameCode
		pg100.PlayInfo = body.Transactions.PlayInfo
		pg100.IsEndRound = body.Transactions.IsEndRound
		pg100.CreatedAt = time.Now()
		return nil
	})
	fmt.Println(data)
	// find user
	// now := time.Now()
	// timestamp := now.UnixNano() / int64(time.Millisecond)
	// body.Balance = data.Data.Balance
	// body.TimestampMillis = timestamp

	return utils.SuccessResponse(c, body, "success")
}

func settleServer(data models.SettleCheckResponse) (models.SettleCheckResponseFormBackend, error) {
	url := fmt.Sprintf("%s/settleGame", urlBankend)
	reqBody, err := json.Marshal(map[string]interface{}{
		"username":  data.Username,
		"betsettle": data.Transactions.PayoutAmount - data.Transactions.BetAmount,
	})
	if err != nil {
		return models.SettleCheckResponseFormBackend{}, fmt.Errorf("failed to marshal request body: %v", err)
	}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return models.SettleCheckResponseFormBackend{}, fmt.Errorf("failed to create HTTP request: %v", err)
	}
	// Set the required headers
	req.Header.Set("x-api-key", apiKeyBankend)
	req.Header.Set("Content-Type", "application/json")

	// Execute the HTTP request
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return models.SettleCheckResponseFormBackend{}, fmt.Errorf("failed to send HTTP request: %v", err)
	}
	defer resp.Body.Close()
	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		return models.SettleCheckResponseFormBackend{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	// Decode the response body into a JSON map
	var responseMap models.SettleCheckResponse
	if err := json.NewDecoder(resp.Body).Decode(&responseMap); err != nil {
		return models.SettleCheckResponseFormBackend{}, fmt.Errorf("failed to decode response body: %v", err)
	}
	fmt.Println(responseMap)
	return models.SettleCheckResponseFormBackend{}, nil
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
