package route

import (
	"github.com/gofiber/fiber/v2"

	"game_services/app/controllers"
)

func SetRoute(app *fiber.App) {
	gameProvider := app.Group("/game-provider")
	gameProvider.Post("/", controllers.GameProvider)
	gameProvider.Get("/aa", controllers.GameProviderAA)

	api := app.Group("/api")
	api.Post("/launch", controllers.LaunchProvider)
	api.Post("/products/:id", controllers.ProductsByCategory)
}
