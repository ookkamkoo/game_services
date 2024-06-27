package route

import (
	"fmt"

	"github.com/gofiber/fiber/v2"

	"game_services/app/controllers"
)

func SetRoute(app *fiber.App) {

	fmt.Println("sssss")
	gameProvider := app.Group("/game-provider/api")
	gameProvider.Post("/balance", controllers.GameProvider)
	// gameProvider.Get("/aa", controllers.GameProviderAA)

	api := app.Group("/api")
	api.Get("/products/:categoryId", controllers.ProductsByCategory)
	api.Get("/game-list/:categoryId/:productId", controllers.GameList)
	api.Post("/launch-game", controllers.LaunchGame)
	api.Get("/user-information/:username", controllers.UserInformation)

}
