package controllers

import (
	"fmt"
	"game_services/app/database"
	"game_services/app/models"
	"game_services/app/utils"
	"net/http"
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

	fmt.Println(yesterdayStartFormatted)
	fmt.Println(yesterdayEndFormatted)

	if err := database.DB.Model(&models.Reports{}).
		Select("user_id, CAST(SUM(bet_winloss) AS FLOAT) as bet_winloss").
		Where("created_at >= ? AND created_at <= ?", yesterdayStartFormatted, yesterdayEndFormatted).Group("user_id").Having("SUM(bet_winloss) < 0").Find(&results).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to fetch bet win/loss summary",
		})
	}

	return utils.SuccessResponse(c, results, "Bet win/loss summary retrieved successfully.")
}

func GetReportGame(c *fiber.Ctx) error {
	var body models.ReportGameRequest
	if err := c.BodyParser(&body); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Failed to parse request body.", err.Error())
	}
	fmt.Println(body)
	// Create the query using the provided start and end date/time
	query := database.DB.Model(&models.Reports{}).
		Where("created_at BETWEEN ? AND ?", body.DateTimeStart, body.DateTimeEnd)

	// Apply additional filters based on request parameters
	if body.DateSelect != "" && body.DateSelect != "all" {
		query = query.Where("status = ?", body.DateSelect)
	}

	if body.Amount != 0 {
		query = query.Where("bet_winloss = ?", body.Amount)
	}

	// Apply pagination if necessary
	if body.Page > 0 && body.PageSize > 0 {
		offset := (body.Page - 1) * body.PageSize
		query = query.Offset(offset).Limit(body.PageSize)
	}

	// Clone the query for counting
	countQuery := query

	// Execute the query to find the matching report games
	var reportGames []models.Reports
	if err := query.Find(&reportGames).Error; err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Failed to get report game data.", err.Error())
	}

	// Count the total records
	var count int64
	if err := countQuery.Count(&count).Error; err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Failed to count report game data.", err.Error())
	}

	// Prepare the response
	response := fiber.Map{
		"data":         reportGames,
		"recordsTotal": count,
	}
	// Return the response
	return utils.SuccessResponse(c, response, "Get report game successfully.")
}
