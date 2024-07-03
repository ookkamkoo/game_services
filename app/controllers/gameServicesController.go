package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"game_services/app/utils"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

const privateURL = "https://api-test.gpsuperapi.com/api"
const operator_token = "445e6ffe-6d36-4fe9-93a7-f4ca25289058:be423e20-e1ba-4296-9044-3cfc6f7424cd"
const key = "MNzLhy68lkH418xGYFE41XkKvoiRr2FX"

type LaunchRequest struct {
	PlayerUsername string `json:"playerUsername"`
	DeviceType     string `json:"deviceType"`
	Lang           string `json:"lang"`
	ReturnURL      string `json:"returnUrl"`
	PlayerIP       string `json:"playerIp"`
	LaunchCode     string `json:"launchCode"`
	AuthToken      string `json:"authToken"`
}

func ProductsByCategory(c *fiber.Ctx) error {
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

func LaunchGame(c *fiber.Ctx) error {
	// Parse JSON request body into LaunchRequest struct
	var launchReq LaunchRequest
	if err := c.BodyParser(&launchReq); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request payload",
			"error":   err.Error(),
		})
	}

	// Validate required fields
	if launchReq.PlayerUsername == "" || launchReq.LaunchCode == "" {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "playerUsername, launchCode  required",
		})
	}

	encrypted, err := utils.Encrypt(operator_token, key)
	if err != nil {
		fmt.Println("Encryption error:", err)
		return err
	}

	// Construct the URL for the external API request
	url := fmt.Sprintf("%s/v1/launch", privateURL)
	reqBody, err := json.Marshal(map[string]interface{}{
		"playerUsername": launchReq.PlayerUsername,
		"deviceType":     launchReq.DeviceType,
		"lang":           launchReq.Lang,
		"returnUrl":      launchReq.ReturnURL,
		"playerIp":       launchReq.PlayerIP,
		"launchCode":     launchReq.LaunchCode,
		// "authToken":      encrypted,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to prepare request body",
			"error":   err.Error(),
		})
	}

	// Create a new POST request
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to create request",
			"error":   err.Error(),
		})
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Authorization-Token", encrypted)

	// Perform the HTTP request
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to make external API call",
			"error":   err.Error(),
		})
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": fmt.Sprintf("Unexpected status code: %d", resp.StatusCode),
		})
	}

	// Decode the response body into a JSON map
	var responseMap map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&responseMap); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to decode response",
			"error":   err.Error(),
		})
	}

	// Return the response
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Response decoded successfully",
		"data":    responseMap,
	})
}

func LaunchGames(c *fiber.Ctx) error {
	productId := c.Params("productId")
	if productId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "productId is required",
		})
	}

	var body BodyLoginPG
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request payload",
			"error":   err.Error(),
		})
	}

	if productId == "pg100" {
		data, err := PGLaunchGames(body)
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
	}
	return utils.SuccessResponse(c, "success", "success")
}
