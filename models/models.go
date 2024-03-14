// Package models содержит общие модели используемые в разных пакетах
package models

import (
	"errors"
)

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

// Stats структура статистики сервиса
type Stats struct {
	URLs  int64 `json:"urls"`
	Users int64 `json:"users"`
}

// UserIDStringType User ID String Type
type UserIDStringType string

// const объявляю всякие константы (тест идиот? как я должен это документировать?)
const (
	// TokenCookieName Token Cookie Name
	TokenCookieName = "token"
	// UserID отдельный тип для приведения к нему из контекста
	UserID = UserIDStringType("userID")
	// AccessToken название токена для спелчекра
	AccessToken = UserIDStringType("access-token")
	// ServerAddressEnv Server Address Env
	ServerAddressEnv = "SERVER_ADDRESS"
	// BaseURLEnv Base URL Env
	BaseURLEnv = "BASE_URL"
	// SavePath save path env
	SavePath = "FILE_STORAGE_PATH"
	// DBConnectString DB Connect String env
	DBConnectString = "DATABASE_DSN"
	//TrustedSubnet trusted subnet setup in CIDR format
	TrustedSubnet = "TRUSTED_SUBNET"
	//CodeLength длина сокращенного кода урл
	CodeLength = 8
	// CodeGenerateAttempts число попыток создать уникальный код
	CodeGenerateAttempts = 20
	// HTTPSEnabled https flag env
	HTTPSEnabled = "ENABLE_HTTPS"
	// CONFIG path to config file (json)
	CONFIG = "CONFIG"
	//RealIP название хидера в котором мы ожидаем реальный айпи клиента от, например, прокси.
	RealIP = "X-Real-IP"
)

// AlreadyHaveThisURLError если урл уже существует
type AlreadyHaveThisURLError struct {
}

// ErrMissingContent контент помечен как удаленный
var ErrMissingContent = errors.New("content deleted")

// Error для реализации интерфейса ошибки
func (au *AlreadyHaveThisURLError) Error() string {
	return "already have this url in db"
}
