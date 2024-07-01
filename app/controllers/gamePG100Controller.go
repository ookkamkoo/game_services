package controllers

import (
	"game_services/app/utils"

	"github.com/gofiber/fiber/v2"
)

func GP100Provider(c *fiber.Ctx) error {
	// Return the roles
	return utils.SuccessResponse(c, "success", "success")
}
