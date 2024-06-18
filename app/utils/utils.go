package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
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

func Encrypt(text string, key string) (string, error) {
	if len(key) != 32 {
		return "", errors.New("key must be 32 bytes long")
	}

	iv := make([]byte, aes.BlockSize) // Initialization vector
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, len(text))
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext, []byte(text))

	// Combine IV and encrypted text
	result := fmt.Sprintf("%s:%s", hex.EncodeToString(iv), hex.EncodeToString(ciphertext))
	return base64.StdEncoding.EncodeToString([]byte(result)), nil
}

// Decrypt decrypts the encrypted text using AES-256-CBC algorithm.
func Decrypt(encryptedText string, key string) (string, error) {
	if len(key) != 32 {
		return "", errors.New("key must be 32 bytes long")
	}

	decoded, err := base64.StdEncoding.DecodeString(encryptedText)
	if err != nil {
		return "", err
	}

	parts := string(decoded)
	ivHex := parts[:32]
	encrypted := parts[33:]

	iv, err := hex.DecodeString(ivHex)
	if err != nil {
		return "", err
	}

	if len(iv) != aes.BlockSize {
		return "", errors.New("IV must be 16 bytes long")
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	ciphertext, err := hex.DecodeString(encrypted)
	if err != nil {
		return "", err
	}

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return string(ciphertext), nil
}
