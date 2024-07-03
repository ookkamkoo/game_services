package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"game_services/app/utils"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

const privateURLPG100 = "https://agent-api.pgf-asw0uz.com"
const apiKey = "OWJxTzlTNzdCRzpWWXVjZ200emhjcGFiTnZ3YzlTNWR3YWhXWk1HMmNpOQ=="

type BodyLoginPG struct {
	Username     string `json:"username"`
	GameCode     string `json:"gameCode"`
	SessionToken string `json:"sessionToken"`
	Language     string `json:"language"`
}

func GP100Provider(c *fiber.Ctx) error {
	// Return the roles
	return utils.SuccessResponse(c, "success", "success")
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
