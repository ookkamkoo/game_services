package route

import (
	"fmt"
	"game_services/app/controllers"

	"github.com/gofiber/fiber/v2"
)

// logMiddleware logs incoming requests
func logMiddleware(c *fiber.Ctx) error {
	fmt.Printf("Request received: %s %s\n", c.Method(), c.Path())
	return c.Next()
}

// SetRoute sets up all routes for the application
func SetRoute(app *fiber.App) {
	// Apply the logMiddleware to log all incoming requests
	app.Use(logMiddleware)

	// Define routes for game provider APIs
	gameProvider := app.Group("/game-provider/api")
	gameProvider.Post("/balance", controllers.GameProvider)
	// Example: gameProvider.Get("/aa", controllers.GameProviderAA)
	gamePG := app.Group("/game-pg-provider/")
	gamePG.Get("/checkBalance", controllers.CheckBalance)

	// Define routes for general APIs
	api := app.Group("/api")
	// api.Get("/products/:categoryId", controllers.ProductsByCategory)
	api.Get("/game-list/:categoryId/:productId", controllers.GameList)
	// api.Post("/launch-game", controllers.LaunchGame)
	api.Post("/launch-games/:productId", controllers.LaunchGames)
	// api.Get("/user-information/:username", controllers.UserInformation)
}
