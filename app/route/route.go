package route

import (
	"github.com/gofiber/fiber/v2"

	"game_services/app/controllers"
)

func SetRoute(app *fiber.App) {
	gameProvider := app.Group("/game-provider")
	gameProvider.Post("/", controllers.GameProvider)
	gameProvider.Post("/aa", controllers.GameProviderAA)
}
