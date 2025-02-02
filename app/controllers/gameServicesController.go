package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"game_services/app/utils"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
)

// const privateURL = "https://api-test.gpsuperapi.com/api"
// const operator_token = "d5f5232e-bd91-47c7-acc9-0856b6e8a06a:bd8d067f-2623-4afb-ba5f-32896fce47a5"
// const key = "MNzLhy68lkH418xGYFE41XkKvoiRr2FX"

type LaunchRequest struct {
	PlayerUsername string `json:"playerUsername"`
	DeviceType     string `json:"deviceType"`
	Lang           string `json:"lang"`
	ReturnURL      string `json:"returnUrl"`
	PlayerIP       string `json:"playerIp"`
	LaunchCode     string `json:"launchCode"`
	CurrencyCode   string `json:"currencyCode"`
	AuthToken      string `json:"authToken"`
}

var privateURL string
var operator_token string
var key string

func SetValueFormENV() {
	privateURL = os.Getenv("privateURL")
	operator_token = os.Getenv("operator_token")
	key = os.Getenv("key")
}

func ProductsByCategory(c *fiber.Ctx) error {
	SetValueFormENV()
	// Get the categoryId query parameter from the request
	categoryId := c.Params("categoryId")
	if categoryId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "categoryId is required",
		})
	}

	// Construct the URL with the categoryId query parameter
	url := fmt.Sprintf("%s/v1/products?categoryId=%s", privateURL, categoryId)
	fmt.Println(url)
	// Create a new GET request
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to create request",
			"error":   err.Error(),
		})
	}
	fmt.Println("operator_token : " + operator_token)
	fmt.Println("key : " + key)
	encrypted, err := utils.Encrypt(operator_token, key)
	if err != nil {
		fmt.Println("Encryption error:", err)
		return err
	}
	fmt.Println(encrypted)
	req.Header.Set("X-Authorization-Token", encrypted)
	// Send the request
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch products",
			"error":   err.Error(),
		})
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch products",
			"error":   fmt.Sprintf("unexpected status code: %d", resp.StatusCode),
		})
	}

	// Decode the response body into a JSON array
	var responseMap map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&responseMap); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to decode response",
			"error":   err.Error(),
		})
	}

	// Return the response map
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Response decoded successfully",
		"data":    responseMap,
	})
}

func GameList(c *fiber.Ctx) error {
	SetValueFormENV()
	// Get the categoryId query parameter from the request
	categoryId := c.Params("categoryId")
	productId := c.Params("productId")
	if categoryId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "categoryId is required",
		})
	}

	if productId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "productId is required",
		})
	}

	if productId == "pg100" {

		data, err := PGGameList()
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": err,
			})
		}

		// Return the response map
		return c.JSON(fiber.Map{
			"status":  "success",
			"message": "Response decoded successfully",
			"data":    data,
		})

	} else {
		url := fmt.Sprintf("%s/v1/games?categoryId=%s&productId=%s&page=1&limit=100", privateURL, categoryId, productId)

		// Create a new GET request
		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "Failed to create request",
				"error":   err.Error(),
			})
		}

		fmt.Println("operator_token : " + operator_token)
		fmt.Println("key : " + key)
		fmt.Println("url : " + url)

		encrypted, err := utils.Encrypt(operator_token, key)
		if err != nil {
			fmt.Println("Encryption error:", err)
			return err
		}
		req.Header.Set("X-Authorization-Token", encrypted)
		client := http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "Failed to fetch game list",
				"error":   err.Error(),
			})
		}
		defer resp.Body.Close()

		// Check the response status code
		if resp.StatusCode != http.StatusOK {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "Failed to fetch game list",
				"error":   fmt.Sprintf("unexpected status code: %d", resp.StatusCode),
			})
		}

		// Decode the response body into a JSON array
		var responseMap map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&responseMap); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "Failed to decode response",
				"error":   err.Error(),
			})
		}

		// Return the response map
		return c.JSON(fiber.Map{
			"status":  "success",
			"message": "Response decoded successfully",
			"data":    responseMap,
		})
	}
}

func UserInformation(c *fiber.Ctx) error {
	SetValueFormENV()
	// Get the categoryId query parameter from the request
	username := c.Params("username")
	if username == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "username is required",
		})
	}

	url := fmt.Sprintf("%s/v1/user-information?q=%s", privateURL, username)

	// Create a new GET request
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to create request",
			"error":   err.Error(),
		})
	}

	encrypted, err := utils.Encrypt(operator_token, key)
	if err != nil {
		fmt.Println("Encryption error:", err)
		return err
	}
	req.Header.Set("X-Authorization-Token", encrypted)
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch username",
			"error":   err.Error(),
		})
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch username",
			"error":   fmt.Sprintf("unexpected status code: %d", resp.StatusCode),
		})
	}

	// Decode the response body into a JSON array
	var responseMap map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&responseMap); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to decode response",
			"error":   err.Error(),
		})
	}

	// Return the response map
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Response decoded successfully",
		"data":    responseMap,
	})
}

func launchGameGplay(launchReq LaunchRequest) (string, error) {
	// SetValueFormENV()
	// Parse JSON request body into LaunchRequest struct
	fmt.Println("aaaaaaaaaaaaa")
	// Validate required fields
	fmt.Println(launchReq)
	if launchReq.PlayerUsername == "" || launchReq.LaunchCode == "" {
		return "", errors.New("playerUsername and launchCode are required")
	}

	encrypted, err := utils.Encrypt(operator_token, key)
	if err != nil {
		fmt.Println("Encryption error:", err)
		return "", fmt.Errorf("failed to encrypt token: %w", err)
	}
	url := fmt.Sprintf("%s/v1/launch", privateURL)
	reqBody, err := json.Marshal(map[string]interface{}{
		"playerUsername": launchReq.PlayerUsername,
		"deviceType":     launchReq.DeviceType,
		"lang":           launchReq.Lang,
		"returnUrl":      launchReq.ReturnURL,
		"playerIp":       launchReq.PlayerIP,
		"launchCode":     launchReq.LaunchCode,
		"currencyCode":   launchReq.CurrencyCode,
	})

	fmt.Println(reqBody)

	if err != nil {
		return "", fmt.Errorf("failed to prepare request body: %w", err)
	}

	// Create a new POST request
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Authorization-Token", encrypted)

	// Perform the HTTP request
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to make external API call: %w", err)
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Decode the response body into a JSON map
	var responseMap struct {
		Code int                    `json:"code"`
		Data map[string]interface{} `json:"data"`
		Msg  string                 `json:"msg"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&responseMap); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	// Check if the response contains an error code
	if responseMap.Code != 0 {
		return "", fmt.Errorf("error from external API: %s", responseMap.Msg)
	}

	// Retrieve the game URL from the response data
	gameUrl, ok := responseMap.Data["gameUrl"].(string)
	if !ok {
		return "", errors.New("missing gameUrl in response")
	}

	// Return the game URL on success
	return gameUrl, nil
}

func LaunchGames(c *fiber.Ctx) error {
	SetValueFormENV()
	productId := c.Params("productId")
	if productId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "productId is required",
		})
	}
	fmt.Println("sssssssssssssssssss")
	fmt.Println(productId)
	if productId == "pg100" {
		var body BodyLoginPG
		if err := c.BodyParser(&body); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "Invalid request payload",
				"error":   err.Error(),
			})
		}

		data, err := PGLaunchGames(body)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "Error Launch Game",
				"error":   err.Error(),
			})
		}

		err = PGSettingGame(body.Setting)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "Error Setting Game",
				"error":   err.Error(),
			})
		}

		// Return the response map
		return c.JSON(fiber.Map{
			"status":  "success",
			"message": "Response decoded successfully",
			"data":    data,
		})
	} else {
		var launchReq LaunchRequest
		if err := c.BodyParser(&launchReq); err != nil {
			fmt.Println(err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "Invalid request payload",
				"error":   err.Error(),
			})
		}
		fmt.Println("inqqqqqqqqqqqqqqqqqq")
		gameUrl, err := launchGameGplay(launchReq)
		if err != nil {
			fmt.Println(err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": err.Error(),
			})
		}

		// return c.JSON(fiber.Map{
		// 	"status":  "success",
		// 	"gameUrl": gameUrl,
		// })
		response := fiber.Map{
			"status": "success",
			"data": fiber.Map{
				"url": gameUrl,
			},
		}
		return utils.SuccessResponse(c, response, "success")
	}
}

func SettingGamePg100(c *fiber.Ctx) error {
	SetValueFormENV()
	var body json.RawMessage
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request payload",
			"error":   err.Error(),
		})
	}

	fmt.Println("gameConfigJSON:", string(body))
	err := PGSettingGame(body)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Error Setting Game",
			"error":   err.Error(),
		})
	}

	return utils.SuccessResponse(c, "success", "success")
}
