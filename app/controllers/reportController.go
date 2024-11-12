package controllers

import (
	"fmt"
	"game_services/app/database"
	"game_services/app/models"
	"game_services/app/utils"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
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

	query2 := database.DB.Model(&models.Reports{}).
		Select("COALESCE(SUM(bet_amount), 0) AS totalBetAmount, COALESCE(SUM(bet_result), 0) AS totalBetResult, COALESCE(SUM(bet_winloss), 0) AS totalBetWinloss").
		Where("created_at BETWEEN ? AND ?", body.DateTimeStart, body.DateTimeEnd)

	// Apply additional filters based on request parameters
	if body.Username != "" {
		query = query.Where("username LIKE ?", "%"+body.Username+"%")
		query2 = query2.Where("username LIKE ?", "%"+body.Username+"%")
	}

	if body.DateSelect != "" && body.DateSelect != "all" {
		query = query.Where("status = ?", body.DateSelect)
		query2 = query2.Where("status = ?", body.DateSelect)
	}

	if body.Amount != 0 {
		query = query.Where("bet_winloss = ?", body.Amount)
		query2 = query2.Where("bet_winloss = ?", body.Amount)
	}

	// Apply pagination if necessary
	if body.Page > 0 && body.PageSize > 0 {
		offset := (body.Page - 1) * body.PageSize
		query = query.Offset(offset).Limit(body.PageSize)
		fmt.Println("Pagination Offset:", offset)
		fmt.Println("Pagination Limit:", body.PageSize)
	}

	// Clone the query for counting
	countQuery := database.DB.Model(&models.Reports{}).
		Where("created_at BETWEEN ? AND ?", body.DateTimeStart, body.DateTimeEnd)

	if body.Username != "" {
		countQuery = countQuery.Where("username LIKE ?", "%"+body.Username+"%")
	}

	if body.DateSelect != "" && body.DateSelect != "all" {
		countQuery = countQuery.Where("status = ?", body.DateSelect)
	}

	if body.Amount != 0 {
		countQuery = countQuery.Where("bet_winloss = ?", body.Amount)
	}

	// Execute the query to find the matching report games
	var reportGames []models.Reports
	if err := query.Order("id desc").Find(&reportGames).Error; err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Failed to get report game data.", err.Error())
	}

	// Count the total records
	var count int64
	if err := countQuery.Count(&count).Error; err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Failed to count report game data.", err.Error())
	}

	fmt.Println("Total Records:", count)

	// Calculate sums of bet_amount, bet_result, and bet_winloss
	var sums struct {
		TotalBetAmount  float64 `gorm:"column:totalBetAmount"`
		TotalBetResult  float64 `gorm:"column:totalBetResult"`
		TotalBetWinloss float64 `gorm:"column:totalBetWinloss"`
	}

	err := query2.Row().Scan(&sums.TotalBetAmount, &sums.TotalBetResult, &sums.TotalBetWinloss)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Failed to calculate sums.", err.Error())
	}

	// Prepare the response
	response := fiber.Map{
		"data":          reportGames,
		"recordsTotal":  count,
		"sumBetAmount":  sums.TotalBetAmount,
		"sumBetResult":  sums.TotalBetResult,
		"sumBetWinloss": sums.TotalBetWinloss,
	}

	// Return the response
	return utils.SuccessResponse(c, response, "Get report game successfully.")
}

func GetReportGameProduct(c *fiber.Ctx) error {

	// Struct for holding the sums grouped by product_id
	type SumResult struct {
		ProductName string  `json:"product_name"` // เปลี่ยนเป็น string
		WinLose     float64 `json:"win_lose"`
	}

	var sums []SumResult

	now := time.Now()
	yesterday := now.AddDate(0, 0, -1)
	yesterdayStart := time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 0, 0, 0, 0, yesterday.Location())
	yesterdayEnd := time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 23, 59, 59, 0, yesterday.Location())
	layout := "2006-01-02 15:04:05"
	yesterdayStartFormatted := yesterdayStart.Format(layout)
	yesterdayEndFormatted := yesterdayEnd.Format(layout)

	fmt.Println(yesterdayStartFormatted)
	fmt.Println(yesterdayEndFormatted)

	// Perform the query with GROUP BY product_id
	if err := database.DB.Model(&models.Reports{}).
		Select("product_name, SUM(bet_winloss) AS win_lose").
		Where("created_at BETWEEN ? AND ?", yesterdayStartFormatted, yesterdayEndFormatted).
		Group("product_name").
		Scan(&sums).Error; err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Failed to calculate sums.", err.Error())
	}

	// Prepare the response
	response := fiber.Map{
		"data": sums,
	}

	// Return the response
	return utils.SuccessResponse(c, response, "Get report game successfully.")
}

func GetReportGameByProductName(c *fiber.Ctx) error {

	// Struct for holding the sums grouped by product_name
	type SumResult struct {
		ProductName string  `json:"product_name"`
		WinLose     float64 `json:"win_lose"`
	}

	var sums []SumResult

	// Retrieve and validate query parameters
	dateTimeStart := c.Query("dateTimeStart")
	dateTimeEnd := c.Query("dateTimeEnd")

	// Check if dates are provided
	if dateTimeStart == "" || dateTimeEnd == "" {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Missing required date parameters.", "dateTimeStart or dateTimeEnd is missing.")
	}

	// Optional: Parse dates to ensure they are valid (assuming format "2006-01-02 15:04:05")
	_, err := time.Parse("2006-01-02 15:04:05", dateTimeStart)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid dateTimeStart format.", err.Error())
	}
	_, err = time.Parse("2006-01-02 15:04:05", dateTimeEnd)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid dateTimeEnd format.", err.Error())
	}

	// Perform the query with GROUP BY product_name
	if err := database.DB.Model(&models.Reports{}).
		Select("product_name, SUM(bet_winloss) AS win_lose").
		Where("created_at BETWEEN ? AND ?", dateTimeStart, dateTimeEnd).
		Group("product_name").
		Scan(&sums).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.ErrorResponse(c, http.StatusNotFound, "No records found for the specified date range.", "")
		}
		return utils.ErrorResponse(c, http.StatusBadRequest, "Failed to calculate sums.", err.Error())
	}

	// Prepare the response
	response := fiber.Map{
		"data": sums,
	}
	fmt.Println(response)

	// Return the response
	return utils.SuccessResponse(c, response, "Get report game successfully.")
}

func GetReportGameByCategoryName(c *fiber.Ctx) error {

	// Struct for holding the sums grouped by product_name
	type SumResult struct {
		ProductName string  `json:"product_name"`
		WinLose     float64 `json:"win_lose"`
	}

	var sums []SumResult

	// Retrieve and validate query parameters
	dateTimeStart := c.Query("dateTimeStart")
	dateTimeEnd := c.Query("dateTimeEnd")

	// Check if dates are provided
	if dateTimeStart == "" || dateTimeEnd == "" {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Missing required date parameters.", "dateTimeStart or dateTimeEnd is missing.")
	}

	// Optional: Parse dates to ensure they are valid (assuming format "2006-01-02 15:04:05")
	_, err := time.Parse("2006-01-02 15:04:05", dateTimeStart)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid dateTimeStart format.", err.Error())
	}
	_, err = time.Parse("2006-01-02 15:04:05", dateTimeEnd)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid dateTimeEnd format.", err.Error())
	}

	// Perform the query with GROUP BY product_name
	if err := database.DB.Model(&models.Reports{}).
		Select("category_name, SUM(bet_winloss) AS win_lose").
		Where("created_at BETWEEN ? AND ?", dateTimeStart, dateTimeEnd).
		Group("category_name").
		Scan(&sums).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.ErrorResponse(c, http.StatusNotFound, "No records found for the specified date range.", "")
		}
		return utils.ErrorResponse(c, http.StatusBadRequest, "Failed to calculate sums.", err.Error())
	}

	// Prepare the response
	response := fiber.Map{
		"data": sums,
	}
	fmt.Println(response)

	// Return the response
	return utils.SuccessResponse(c, response, "Get report game successfully.")
}
