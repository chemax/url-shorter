package util

import (
	"errors"
	"math/rand"
	"strings"
)

// AlreadyHaveThisURLError если урл уже существует
type AlreadyHaveThisURLError struct {
}

// Error для реализации интерфейса ошибки
func (au *AlreadyHaveThisURLError) Error() string {
	return "already have this url in db"
}

// DeleteTask задача на удаление
type DeleteTask struct {
	Codes  []string
	UserID string
}

// URLStructUser URL Struct
type URLStructUser struct {
	Shortcode string `json:"short_url"`
	URL       string `json:"original_url"`
}

// URLStructForBatch URL Struct For Batch
type URLStructForBatch struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

// URLStructForBatchResponse URL Struct For Batch Response
type URLStructForBatchResponse struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

// UserIDStringType User ID String Type
type UserIDStringType string

const (
	TokenCookieName      = "token"
	UserID               = UserIDStringType("userID")
	ServerAddressEnv     = "SERVER_ADDRESS"
	BaseURLEnv           = "BASE_URL"
	SavePath             = "FILE_STORAGE_PATH"
	DBConnectString      = "DATABASE_DSN"
	CodeLength           = 8
	CodeGenerateAttempts = 20
)

// ErrMissingContent контент помечен как удаленный
var ErrMissingContent = errors.New("content deleted")

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

// CheckHeaderIsValidType проверка на валидность хидера
func CheckHeaderIsValidType(header string) bool {
	return strings.Contains(header, "application/json") || strings.Contains(header, "application/x-gzip")
}

// CheckHeader проверка на валидность хидера
func CheckHeader(header string) bool {
	return strings.Contains(header, "text/plain") || strings.Contains(header, "application/x-gzip")
}

// RandStringRunes генерация псевдослучайной строки заданной длинны
func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
