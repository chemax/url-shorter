package util

import (
	"math/rand"
	"strings"
)

type AlreadyHaveThisURLError struct {
}

func (au *AlreadyHaveThisURLError) Error() string {
	return "already have this url in db"
}

type URLStructUser struct {
	Shortcode string `json:"short_url"`
	URL       string `json:"original_url"`
}
type URLStructForBatch struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}
type URLStructForBatchResponse struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

const (
	UserID               = "userID"
	ServerAddressEnv     = "SERVER_ADDRESS"
	BaseURLEnv           = "BASE_URL"
	SavePath             = "FILE_STORAGE_PATH"
	DBConnectString      = "DATABASE_DSN"
	CodeLength           = 8
	CodeGenerateAttempts = 20
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

func CheckHeaderIsValidType(header string) bool {
	return strings.Contains(header, "application/json") || strings.Contains(header, "application/x-gzip")
}
func CheckHeader(header string) bool {
	return strings.Contains(header, "text/plain") || strings.Contains(header, "application/x-gzip")
}

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
