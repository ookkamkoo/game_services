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

func PGGameList() (any, error) {
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

func PGLaunchGames(data BodyLoginPG) (any, error) {
	url := fmt.Sprintf("%s/seamless/api/v2/login", privateURLPG100)
	fmt.Println({
		"username":     data.Username,
		"gameCode":     data.GameCode,
		"sessionToken": data.SessionToken,
		"language":     data.Language,
	});
	reqBody, err := json.Marshal(map[string]interface{}{
		"username":     data.Username,
		"gameCode":     data.GameCode,
		"sessionToken": data.SessionToken,
		"language":     data.Language,
	})
	
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("x-api-key", apiKey)
	fmt.Println(apiKey)
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
