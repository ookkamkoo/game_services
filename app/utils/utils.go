package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
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

func pkcs7Pad(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padtext...)
}

func pkcs7Unpad(data []byte, blockSize int) []byte {
	length := len(data)
	unpadding := int(data[length-1])
	return data[:(length - unpadding)]
}

func Encrypt(text string, key string) (string, error) {
	if len(key) != 32 {
		return "", errors.New("key must be 32 bytes long")
	}

	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	paddedText := pkcs7Pad([]byte(text), block.BlockSize())
	encrypted := make([]byte, len(paddedText))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(encrypted, paddedText)

	data := append(iv, encrypted...)
	return base64.StdEncoding.EncodeToString(data), nil
}

func Decrypt(encryptedText string, key string) (string, error) {
	if len(key) != 32 {
		return "", errors.New("key must be 32 bytes long")
	}

	decoded, err := base64.StdEncoding.DecodeString(encryptedText)
	if err != nil {
		return "", errors.New("invalid base64 string")
	}

	if len(decoded) <= aes.BlockSize {
		return "", errors.New("encrypted text is too short")
	}

	iv := decoded[:aes.BlockSize]
	encrypted := decoded[aes.BlockSize:]

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", errors.New("failed to create new cipher")
	}

	decrypted := make([]byte, len(encrypted))
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(decrypted, encrypted)

	decrypted = pkcs7Unpad(decrypted, block.BlockSize())
	return string(decrypted), nil
}
