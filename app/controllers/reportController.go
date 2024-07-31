package controllers

import (
	"game_services/app/database"
	"game_services/app/models"
	"game_services/app/utils"
	"time"

	"github.com/gofiber/fiber/v2"
)

type BetWinLossSummary struct {
	UserID     uint    `json:"user_id"`
	BetWinLoss float32 `json:"bet_winloss"`
}

func GetBetWinLossSummary(c *fiber.Ctx) error {
	var results []models.BetWinLossSummary

	now := time.Now()
	yesterday := now.AddDate(0, 0, -1)
	yesterdayStart := time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 0, 0, 0, 0, yesterday.Location())
	yesterdayEnd := time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 23, 59, 59, 0, yesterday.Location())
	layout := "2006-01-02 15:04:05"
	yesterdayStartFormatted := yesterdayStart.Format(layout)
	yesterdayEndFormatted := yesterdayEnd.Format(layout)

	if err := database.DB.Model(&models.Reports{}).
		Select("member_id, CAST(SUM(bet_winloss) AS FLOAT) as bet_winloss").
		Where("created_at >= ? AND created_at <= ?", yesterdayStartFormatted, yesterdayEndFormatted).Group("member_id").Having("SUM(bet_winloss) < 0").Find(&results).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to fetch bet win/loss summary",
		})
	}

	return utils.SuccessResponse(c, results, "Bet win/loss summary retrieved successfully.")
}
