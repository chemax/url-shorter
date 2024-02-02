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

// URLWithShort URL Struct
type URLWithShort struct {
	Shortcode string `json:"short_url"`
	URL       string `json:"original_url"`
}

// URLForBatch URL Struct For Batch
type URLForBatch struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

// URLForBatchResponse URL Struct For Batch Response
type URLForBatchResponse struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

// UserIDStringType User ID String Type
type UserIDStringType string

// const объявляю всякие константы (тест идиот? как я должен это документировать?)
const (
	// TokenCookieName Token Cookie Name
	TokenCookieName = "token"
	// UserID отдельный тип для приведения к нему из контекста
	UserID = UserIDStringType("userID")
	// ServerAddressEnv Server Address Env
	ServerAddressEnv = "SERVER_ADDRESS"
	// BaseURLEnv Base URL Env
	BaseURLEnv = "BASE_URL"
	// SavePath save path env
	SavePath = "FILE_STORAGE_PATH"
	// DBConnectString DB Connect String env
	DBConnectString = "DATABASE_DSN"
	//CodeLength длина сокращенного кода урл
	CodeLength = 8
	// CodeGenerateAttempts число попыток создать уникальный код
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
