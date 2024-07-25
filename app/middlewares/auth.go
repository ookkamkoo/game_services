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
		apikey := c.Get("X-Api-Key")
		fmt.Println(apikey)
		OsApiKey := os.Getenv("API_KEY_BACKEND")
		if apikey == OsApiKey {
			return c.Next()
		} else {
			return utils.ErrorResponse(c, http.StatusUnauthorized, "You don't have permission to access this resource", "You don't have permission to access this resource")
		}
	}
}
