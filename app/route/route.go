package route

import (
	"fmt"
	"game_services/app/controllers"
	"game_services/app/middlewares"
	"game_services/app/migration"
	"game_services/app/utils"

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

	app.Get("/migration", migrateHandler)
	// Define routes for game provider APIs
	gameProvider := app.Group("/game-provider/api")
	gameProvider.Post("/balance", controllers.GameProvider)
	// Example: gameProvider.Get("/aa", controllers.GameProviderAA)
	gamePG := app.Group("/game-pg-provider/")
	gamePG.Post("/checkBalance", controllers.CheckBalancePG)
	gamePG.Post("/settleBets", controllers.SettleBetsPG)

	// Define routes for general APIs
	api := app.Group("/api")
	api.Get("/products/:categoryId", controllers.ProductsByCategory)
	api.Get("/game-list/:categoryId/:productId", middlewares.GameSeviceMiddleware(), controllers.GameList)
	// api.Post("/launch-game", controllers.LaunchGame)
	api.Post("/launch-games/:productId", controllers.LaunchGames)
	api.Post("/settingPg", controllers.SettingGamePg100)
	api.Get("/getRefoundLost", middlewares.GameSeviceMiddleware(), controllers.GetBetWinLossSummary)
	api.Post("/getReportGame", middlewares.GameSeviceMiddleware(), controllers.GetReportGame)
	api.Get("/getReportGameProduct", middlewares.GameSeviceMiddleware(), controllers.GetReportGameProduct)
	api.Get("/getReportGameByProductName", middlewares.GameSeviceMiddleware(), controllers.GetReportGameByProductName)
	// api.Get("/user-information/:username", controllers.UserInformation)
	utils.Encrypt("", "")
}

func migrateHandler(c *fiber.Ctx) error {
	// Run migrations
	migration.RunMigration()

	// Return a success response
	return c.SendString("Migrations completed successfully")
}
