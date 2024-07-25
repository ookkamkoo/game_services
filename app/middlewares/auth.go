package middlewares

import (
	"fmt"
	"game_services/app/utils"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	// "fmt"
)

func GameSeviceMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		apikey := c.Get("x-api-key")
		OsApiKey := os.Getenv("API_KEY_BACKEND")
		fmt.Println("apikey = ", apikey)
		fmt.Println("OsApiKey = ", OsApiKey)

		if apikey == OsApiKey {
			return c.Next()
		} else {
			return utils.ErrorResponse(c, http.StatusUnauthorized, "You don't have permission to access this resource", "You don't have permission to access this resource")
		}
	}
}
