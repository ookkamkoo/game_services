package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"game_services/app/utils"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
)

const privateURLPG100 = "https://agent-api.pgf-asw0uz.com"
const apiKey = "OWJxTzlTNzdCRzpWWXVjZ200emhjcGFiTnZ3YzlTNWR3YWhXWk1HMmNpOQ=="

const apiKeyBankend = "BKw7jpQd8SOv7LuqPFq6MgQ4A1TflW4Ls"
const urlBankend = "https://backend.scbbbb.com/game-services/"

type BodyLoginPG struct {
	Username     string `json:"username"`
	GameCode     string `json:"gameCode"`
	SessionToken string `json:"sessionToken"`
	Language     string `json:"language"`
}

type BalanceCheckResponse struct {
	ID              string        `json:"id"`
	StatusCode      int           `json:"statusCode"`
	TimestampMillis int64         `json:"timestampMillis"`
	ProductId       string        `json:"productId"`
	Currency        string        `json:"currency"`
	Balance         float32       `json:"balance"`
	Username        string        `json:"username"`
	Transactions    []Transaction `json:"txns"`
	BalanceBefore   float32       `json:"balanceBefore"`
	BalanceAfter    float32       `json:"balanceAfter"`
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

func CheckBalancePG(c *fiber.Ctx) error {
	var body BalanceCheckResponse
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request payload",
			"error":   err.Error(),
		})
	}
	fmt.Println("find user")
	// find user
	data, err := GetBalancePG(body.Username)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": err,
		})
	}
	// find user
	now := time.Now()
	timestamp := now.UnixNano() / int64(time.Millisecond)
	body.Balance = data.Data.Balance
	body.TimestampMillis = timestamp

	return utils.SuccessResponse(c, body, "success")
}

func GetBalancePG(username string) (ResponseData, error) {
	url := fmt.Sprintf("%s/getBalance", urlBankend)
	reqBody, err := json.Marshal(map[string]interface{}{
		"username": username,
	})
	fmt.Println(username)
	if err != nil {
		return ResponseData{}, fmt.Errorf("failed to marshal request body: %v", err)
	}
	fmt.Println("11111111")
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return ResponseData{}, fmt.Errorf("failed to create HTTP request: %v", err)
	}
	fmt.Println("222222222")
	// Set the required headers
	req.Header.Set("x-api-key", apiKeyBankend)
	req.Header.Set("Content-Type", "application/json")

	// Execute the HTTP request
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return ResponseData{}, fmt.Errorf("failed to send HTTP request: %v", err)
	}
	defer resp.Body.Close()
	fmt.Println("3333333333333")
	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		return ResponseData{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	fmt.Println("44444444444444")
	// Decode the response body into a JSON map
	var responseMap ResponseData
	if err := json.NewDecoder(resp.Body).Decode(&responseMap); err != nil {
		return ResponseData{}, fmt.Errorf("failed to decode response body: %v", err)
	}
	fmt.Println(responseMap)
	return responseMap, nil
}

func SettleBetsPG(c *fiber.Ctx) error {
	var body BalanceCheckResponse
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request payload",
			"error":   err.Error(),
		})
	}
	// find user
	data, err := GetBalancePG(body.Username)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": err,
		})
	}

	// find user
	now := time.Now()
	timestamp := now.UnixNano() / int64(time.Millisecond)
	body.Balance = data.Data.Balance
	body.TimestampMillis = timestamp

	return utils.SuccessResponse(c, body, "success")
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

	// Marshal the data to JSON
	reqBody, err := json.Marshal(map[string]interface{}{
		"username":     data.Username,
		"gameCode":     data.GameCode,
		"sessionToken": data.SessionToken,
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
