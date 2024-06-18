package controllers

import (
	"encoding/json"
	"fmt"
	"game_services/app/utils"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

const privateURL = "https://api-test.gpsuperapi.com/api"

const operator_token = "445e6ffe-6d36-4fe9-93a7-f4ca25289058:be423e20-e1ba-4296-9044-3cfc6f7424cd"
const key = "MNzLhy68lkH418xGYFE41XkKvoiRr2FX"

func LaunchProvider(c *fiber.Ctx) error {
	// Define the data for the POST request
	// launchData := map[string]string{
	// 	"playerUsername": "exampleUser",
	// 	"deviceType":     "desktop",
	// 	"lang":           "en",
	// 	"returnUrl":      "https://example.com/return",
	// 	"playerIp":       "192.168.1.1",
	// 	"gameCode":       "ABC123",
	// 	"authToken":      "your-auth-token",
	// }
	// payload, _ := json.Marshal(launchData)

	// // Make the POST request
	// resp, err := http.Post("https://example.com/api/v1/launch", "application/json", bytes.NewBuffer(payload))
	// if err != nil {
	// 	// Handle the error
	// 	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
	// 		"status":  "error",
	// 		"message": "Failed to make POST request",
	// 		"error":   err.Error(),
	// 	})
	// }
	// defer resp.Body.Close()

	// // Check the response status code
	// if resp.StatusCode != http.StatusOK {
	// 	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
	// 		"status":  "error",
	// 		"message": "Unexpected status code",
	// 		"code":    resp.StatusCode,
	// 	})
	// }

	// Return the success response
	return utils.SuccessResponse(c, "success", "success")
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
	// fmt.Println(url)

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
	// fmt.Println("--------------------------------------")
	// fmt.Println(encrypted)
	// fmt.Println("--------------------------------------")
	req.Header.Set("X-Authorization-Token", encrypted)

	// decode, err := utils.Decrypt(encrypted, key)
	// if err != nil {
	// 	fmt.Println("Encryption error:", err)
	// 	return err
	// }
	// fmt.Println(decode)

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
