package utils

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
)

func GetIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println("Error:", err)
		return ""
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}

	return ""
}

func ErrorResponse(c *fiber.Ctx, status int, message, errorMsg string) error {
	return c.Status(status).JSON(fiber.Map{
		"status":  "error",
		"message": message,
		"error":   errorMsg,
		"time":    time.Now().Format("2006-01-02 15:04:05"),
	})
}

func SuccessResponse(c *fiber.Ctx, data interface{}, message string) error {
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": message,
		"data":    data,
		"time":    time.Now().Format("2006-01-02 15:04:05"),
	})
}
