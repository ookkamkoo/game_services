package controllers

import (
	"fmt"
	"game_services/app/database"
	"game_services/app/models"
	"game_services/app/utils"
	"log"
	"math"
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

	// Load the desired time zone (e.g., "Asia/Bangkok")
	location, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to load time zone",
		})
	}

	currentTime := time.Now()
	fmt.Println("🕒 เวลาปัจจุบัน:", currentTime)

	// Get yesterday's start and end times in the specified time zone
	now := time.Now().In(location)
	yesterday := now.AddDate(0, 0, -1)
	yesterdayStart := time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 0, 0, 0, 0, location)
	yesterdayEnd := time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 23, 59, 59, 0, location)

	// Format dates for SQL query
	layout := "2006-01-02 15:04:05"
	yesterdayStartFormatted := yesterdayStart.Format(layout)
	yesterdayEndFormatted := yesterdayEnd.Format(layout)

	// Debug output
	fmt.Println("Date")
	fmt.Println("Start:", yesterdayStartFormatted)
	fmt.Println("End:", yesterdayEndFormatted)

	// Query database for summary
	if err := database.DB.Model(&models.Reports{}).
		Select("user_id, CAST(SUM(bet_winloss) AS FLOAT) as bet_winloss").
		Where("created_at >= ? AND created_at <= ?", yesterdayStartFormatted, yesterdayEndFormatted).
		Group("user_id").
		Having("SUM(bet_winloss) < 0").
		Find(&results).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to fetch bet win/loss summary",
		})
	}

	// Return successful response
	return utils.SuccessResponse(c, results, "Bet win/loss summary retrieved successfully.")
}

func GetWinLostAlliance(c *fiber.Ctx) error {
	var results []models.BetWinLossSummaryAlliance
	var list []models.Reports

	// รับค่าจาก Query Parameter
	dateStart := c.Query("dateStart")
	timeStart := c.Query("timeStart")
	dateEnd := c.Query("dateEnd")
	timeEnd := c.Query("timeEnd")
	agent_id := c.Query("agent_id")

	page := c.QueryInt("page", 1)
	pageSize := c.QueryInt("pageSize", 10)

	offset := (page - 1) * pageSize

	// ตรวจสอบ agent_id
	if agent_id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "agent_id is required",
		})
	}

	// ตรวจสอบและจัดการฟอร์แมตของวันที่
	layout := "2006-01-02 15:04:05"
	startDateTime := fmt.Sprintf("%s %s", dateStart, timeStart)
	endDateTime := fmt.Sprintf("%s %s", dateEnd, timeEnd)

	startTime, err := time.Parse(layout, startDateTime)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid dateStart or timeStart format. Expected format: YYYY-MM-DD HH:MM:SS",
		})
	}

	endTime, err := time.Parse(layout, endDateTime)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid dateEnd or timeEnd format. Expected format: YYYY-MM-DD HH:MM:SS",
		})
	}

	if !startTime.Before(endTime) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Start time must be earlier than end time.",
		})
	}

	// Debug output
	fmt.Println("Start:", startDateTime)
	fmt.Println("End:", endDateTime)

	// Query database for summary grouped by agent_id
	if err := database.DB.Model(&models.Reports{}).
		Select("agent_id, CAST(SUM(bet_winloss) AS FLOAT) as bet_winloss").
		Where("agent_id = ? AND created_at >= ? AND created_at <= ?", agent_id, startDateTime, endDateTime).
		Group("agent_id").
		Having("SUM(bet_winloss) < 0").
		Find(&results).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to fetch bet win/loss summary: " + err.Error(),
		})
	}

	// Query all records with pagination
	if err := database.DB.Where("agent_id = ? AND created_at >= ? AND created_at <= ?", agent_id, startDateTime, endDateTime).
		Order("id desc").Offset(offset).Limit(pageSize).
		Find(&list).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to fetch bet win/loss records: " + err.Error(),
		})
	}

	// Query total count
	var count int64
	if err := database.DB.Model(&models.Reports{}).
		Where("agent_id = ? AND created_at >= ? AND created_at <= ?", agent_id, startDateTime, endDateTime).
		Count(&count).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to fetch total record count: " + err.Error(),
		})
	}

	// ตรวจสอบว่ามีผลลัพธ์หรือไม่
	// if len(results) == 0 && len(list) == 0 {
	// 	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
	// 		"error":   false,
	// 		"message": "No bet win/loss data found for the specified criteria",
	// 	})
	// }

	// Prepare response data
	data := map[string]interface{}{
		"recordsTotal": count,
		"totalSummary": results,
		"dataList":     list,
		"pagination": map[string]interface{}{
			"page":      page,
			"pageSize":  pageSize,
			"total":     count,
			"totalPage": (count + int64(pageSize) - 1) / int64(pageSize), // คำนวณจำนวนหน้าทั้งหมด
		},
	}

	// Return successful response
	return utils.SuccessResponse(c, data, "Bet win/loss summary retrieved successfully.")
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

	if body.Game != "" {
		query = query.Where("game_name LIKE ?", "%"+body.Game+"%")
		query2 = query2.Where("game_name LIKE ?", "%"+body.Game+"%")
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
		CategoryName string  `json:"category_name"` // เปลี่ยนเป็น string
		WinLose      float64 `json:"win_lose"`
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
		Select("category_name, SUM(bet_winloss) AS win_lose").
		Where("created_at BETWEEN ? AND ?", yesterdayStartFormatted, yesterdayEndFormatted).
		Group("category_name").
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

func GetReportGameByCategorySum(c *fiber.Ctx) error {
	type SumResult struct {
		ReportDate   string  `json:"report_date"` // เพิ่มฟิลด์สำหรับวันที่
		CategoryName string  `json:"category_name"`
		BetAmount    float64 `json:"bet_amount"`
		BetResult    float64 `json:"bet_result"`
		BetWinLoss   float64 `json:"bet_winloss"`
	}

	var sums []SumResult

	// รับค่าพารามิเตอร์วันที่จาก query
	dateTimeStart := c.Query("dateTimeStart")
	dateTimeEnd := c.Query("dateTimeEnd")

	if dateTimeStart == "" || dateTimeEnd == "" {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Missing required date parameters.", "dateTimeStart or dateTimeEnd is missing.")
	}

	// ตั้งค่าโซนเวลาเป็น Asia/Bangkok
	loc, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		log.Println("❌ Error loading timezone:", err)
		return utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to load timezone.", err.Error())
	}

	// แปลง string → time.Time โดยใช้ `time.ParseInLocation`
	startTime, err := time.ParseInLocation("2006-01-02 15:04:05", dateTimeStart, loc)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid dateTimeStart format.", err.Error())
	}
	endTime, err := time.ParseInLocation("2006-01-02 15:04:05", dateTimeEnd, loc)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid dateTimeEnd format.", err.Error())
	}

	// Debug log
	fmt.Println("✅ Start Time:", startTime)
	fmt.Println("✅ End Time:", endTime)

	// Query database โดยใช้ GROUP BY วันที่ และ category_name
	if err := database.DB.Model(&models.Reports{}).
		Select("DATE(created_at) AS report_date, COALESCE(category_name, 'Unknown') AS category_name, SUM(bet_amount) AS bet_amount, SUM(bet_result) AS bet_result, SUM(bet_winloss) AS bet_winloss").
		Where("created_at BETWEEN ? AND ?", startTime, endTime).
		Group("report_date, category_name").
		Order("report_date ASC, category_name ASC").
		Scan(&sums).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.ErrorResponse(c, http.StatusNotFound, "No records found for the specified date range.", "")
		}
		return utils.ErrorResponse(c, http.StatusBadRequest, "Failed to calculate sums.", err.Error())
	}

	// โครงสร้างข้อมูลสำหรับเก็บรายงานแยกตามวัน
	reportData := make(map[string]map[string]fiber.Map)

	// หมวดหมู่เริ่มต้นที่ต้องการให้แสดง
	defaultCategories := []string{
		"Sportsbook", "Live Casino", "Slot Game", "Fishing Hunter", "Game Card",
		"Lotto", "E-Sport", "Poker Game", "Keno", "Crypto Trading", "Pg100",
	}

	// สร้างโครงสร้างข้อมูลให้แต่ละวันมีหมวดหมู่เกมครบถ้วน
	for _, record := range sums {
		reportDate := record.ReportDate

		// ตรวจสอบว่า key วันที่มีอยู่ใน map หรือไม่
		if _, exists := reportData[reportDate]; !exists {
			reportData[reportDate] = make(map[string]fiber.Map)

			// เพิ่มหมวดหมู่เริ่มต้นให้กับวันที่นั้น
			for _, category := range defaultCategories {
				reportData[reportDate][category] = fiber.Map{
					"name":        category,
					"bet_amount":  0.0,
					"bet_result":  0.0,
					"bet_winloss": 0.0,
				}
			}
		}

		// อัปเดตค่าผลรวมที่ดึงมา
		reportData[reportDate][record.CategoryName] = fiber.Map{
			"name":        record.CategoryName,
			"bet_amount":  record.BetAmount,
			"bet_result":  record.BetResult,
			"bet_winloss": math.Round((record.BetResult-record.BetAmount)*100) / 100,
		}
	}

	// Debug log
	fmt.Println("📊 Report Data:", reportData)

	// เตรียม response
	response := fiber.Map{
		"data": reportData,
	}

	// ส่ง response
	return utils.SuccessResponse(c, response, "Get report game successfully.")
}

func GetReportGameByCategorySumByKey(c *fiber.Ctx) error {

	// Struct for holding the sums grouped by product_name
	type SumResult struct {
		CategoryName string  `json:"category_name"`
		BetAmount    float64 `json:"bet_amount"`
		BetResult    float64 `json:"bet_result"`
		BetWinLoss   float64 `json:"bet_winloss"`
	}

	var sums []SumResult

	key := c.Query("key")
	// Perform the query with GROUP BY product_name
	if err := database.DB.Model(&models.Reports{}).
		Select("category_name, SUM(bet_amount) AS bet_amount, SUM(bet_result) AS bet_result, SUM(bet_winloss) AS bet_winloss").
		Where("key_deposit = ?", key).
		Group("category_name").
		Scan(&sums).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.ErrorResponse(c, http.StatusNotFound, "No records found for the specified date range.", "")
		}
		return utils.ErrorResponse(c, http.StatusBadRequest, "Failed to calculate sums.", err.Error())
	}
	// fmt.Println(sums)
	reportData := map[string]fiber.Map{
		"Sportsbook":      {"name": "Sportsbook", "bet_amount": 0, "bet_result": 0, "bet_winloss": 0},
		"Live Casino":     {"name": "Live Casino", "bet_amount": 0, "bet_result": 0, "bet_winloss": 0},
		"Slot Game":       {"name": "Slot Game", "bet_amount": 0, "bet_result": 0, "bet_winloss": 0},
		"Fishing Hunter":  {"name": "Fishing Hunter", "bet_amount": 0, "bet_result": 0, "bet_winloss": 0},
		"Game Card":       {"name": "Game Card", "bet_amount": 0, "bet_result": 0, "bet_winloss": 0},
		"Lotto":           {"name": "Lotto", "bet_amount": 0, "bet_result": 0, "bet_winloss": 0},
		"E-Sport":         {"name": "E-Sport", "bet_amount": 0, "bet_result": 0, "bet_winloss": 0},
		"Poker Game":      {"name": "Poker Game", "bet_amount": 0, "bet_result": 0, "bet_winloss": 0},
		"Keno":            {"name": "Keno", "bet_amount": 0, "bet_result": 0, "bet_winloss": 0},
		"Crypto Tradding": {"name": "Crypto Trading", "bet_amount": 0, "bet_result": 0, "bet_winloss": 0},
		"Pg100":           {"name": "Pg100", "bet_amount": 0, "bet_result": 0, "bet_winloss": 0},
	}

	for _, v := range sums {
		if data, exists := reportData[v.CategoryName]; exists {
			data["bet_amount"] = v.BetAmount
			data["bet_result"] = v.BetResult
			data["bet_winloss"] = math.Round((v.BetResult-v.BetAmount)*100) / 100
			reportData[v.CategoryName] = data
		}
	}

	// Prepare the response
	response := fiber.Map{
		"data": reportData,
	}
	fmt.Println(response)

	// Return the response
	return utils.SuccessResponse(c, response, "Get report game successfully.")
}

func GetReportGameByCategoryName(c *fiber.Ctx) error {

	// Struct สำหรับเก็บข้อมูลการคำนวณ
	type SumResult struct {
		CategoryName string  `json:"category_name"`
		WinLose      float64 `json:"win_lose"`
	}

	// Struct สำหรับรับค่า request
	type RequestBody struct {
		Username  string `json:"username"`
		KeyDepost string `json:"key_depost"`
	}

	// พาร์ส JSON body
	var body RequestBody
	if err := c.BodyParser(&body); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body format.", err.Error())
	}

	var sums []SumResult

	// ทำการ query ข้อมูลพร้อม GROUP BY category_name
	if err := database.DB.Model(&models.Reports{}).
		Select("category_name, SUM(bet_amount) AS win_lose").
		Where("key_deposit = ? AND username = ?", body.KeyDepost, body.Username).
		Group("category_name").
		Scan(&sums).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.ErrorResponse(c, http.StatusNotFound, "No records found for the specified parameters.", "")
		}
		return utils.ErrorResponse(c, http.StatusBadRequest, "Failed to calculate sums.", err.Error())
	}

	// เตรียมข้อมูลประเภทเกมพร้อมค่าเริ่มต้น
	reportData := map[string]fiber.Map{
		"Sportsbook":      {"name": "Sportsbook", "win_lose": 0},
		"Live Casino":     {"name": "Live Casino", "win_lose": 0},
		"Slot Game":       {"name": "Slot Game", "win_lose": 0},
		"Fishing Hunter":  {"name": "Fishing Hunter", "win_lose": 0},
		"Game Card":       {"name": "Game Card", "win_lose": 0},
		"Lotto":           {"name": "Lotto", "win_lose": 0},
		"E-Sport":         {"name": "E-Sport", "win_lose": 0},
		"Poker Game":      {"name": "Poker Game", "win_lose": 0},
		"Keno":            {"name": "Keno", "win_lose": 0},
		"Crypto Tradding": {"name": "Crypto Trading", "win_lose": 0},
		"Pg100":           {"name": "Pg100", "win_lose": 0},
	}

	// อัปเดตค่า win_lose สำหรับประเภทเกมที่มีข้อมูลใน sums
	for _, v := range sums {
		if data, exists := reportData[v.CategoryName]; exists {
			data["win_lose"] = v.WinLose
			reportData[v.CategoryName] = data
		}
	}

	// เตรียมข้อมูลสำหรับการตอบกลับ
	response := fiber.Map{
		"data": reportData,
	}

	// ส่งการตอบกลับ
	return utils.SuccessResponse(c, response, "Get report game successfully.")
}
