package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	// "github.com/joho/godotenv"

	"game_services/app/database"
	"game_services/app/route"
)

func main() {
	app := fiber.New()

	app.Use(cors.New())

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
		AllowMethods: "GET, POST, PUT, DELETE, HEAD, PATCH",
		MaxAge:       3600,
	}))

	// settime zone
	loc, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		fmt.Println("Error loading timezone:", err)
		return
	}
	time.Local = loc

	// Connect  Database
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatal("Error loading env")
	// 	log.Fatal(err)
	// }

	if err := database.PG_Connect(); err != nil {
		log.Fatal(err)
	}

	route.SetRoute(app)

	app.Listen(":3003")
}
