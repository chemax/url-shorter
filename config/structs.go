package config

/*
In this file only type structs and its methods
*/
import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/integrii/flaggy"
)

// Config содержит в себе весь конфиг включая подструктуры
type Config struct {
	NetAddr         string
	HTTPAddr        string
	PathSave        string
	DBConfig        string
	flagInitialized bool
	HTTPSEnabled    bool
	secretKey       string
	tokenExp        time.Duration
	ConfigPath      string
}

type tmpConfig struct {
	ServerAddress   string `json:"server_address"`
	BaseUrl         string `json:"base_url"`
	FileStoragePath string `json:"file_storage_path"`
	DatabaseDsn     string `json:"database_dsn"`
	EnableHttps     bool   `json:"enable_https"`
	ConfigPath      string `json:"-"`
}

//нафиг оверинжиниринг

// SecretKey возвращает секретный ключ
func (c *Config) SecretKey() string {
	return c.secretKey
}

// TokenExp возвращает срок жизни токена
func (c *Config) TokenExp() time.Duration {
	return c.tokenExp
}

// GetDBUse используется ли внешняя база данных или текстовая
func (c *Config) GetDBUse() bool {
	return c.DBConfig != ""
}

// GetSavePath путь к текстовой базе данных
func (c *Config) GetSavePath() string {
	return c.PathSave
}

// GetNetAddr получить сетевой адрес
func (c *Config) GetNetAddr() string {
	return c.NetAddr
}

// GetHTTPAddr возвращает http адрес сервиса (для генерации сокращенных ссылок)
func (c *Config) GetHTTPAddr() string {
	return c.HTTPAddr
}

func (c *Config) initFlags(tmpCfg *tmpConfig) {
	if c.flagInitialized {
		return
	}
	c.flagInitialized = true

	flaggy.String(&tmpCfg.ConfigPath, "c", "config", "config file path")

	flaggy.Bool(&tmpCfg.EnableHttps, "s", "", "enable https")
	flaggy.String(&tmpCfg.ServerAddress, "a", "", "Net address Host:Port")
	flaggy.String(&tmpCfg.BaseUrl, "b", "", "http(s) address http://host:port")
	flaggy.String(&tmpCfg.FileStoragePath, "f", "", "full path to file for save url's")
	flaggy.String(&tmpCfg.DatabaseDsn, "d", "", "DB connect string like \"postgres://username:password@localhost:5432/database_name\"")

	flaggy.Parse()
}

// SetFromTmpConfig вносит в конфиг параметры из временного конфига, который формируется из флагов, энва и файла
// для каждого вида конфигурации свой тмп, вносятся в порядке роста приоритета конфига
func (c *Config) SetFromTmpConfig(tmp *tmpConfig) {
	if tmp.ServerAddress != "" {
		cfg.NetAddr = tmp.ServerAddress
	}
	if tmp.BaseUrl != "" {
		cfg.HTTPAddr = tmp.BaseUrl
	}
	if tmp.FileStoragePath != "" {
		cfg.PathSave = tmp.FileStoragePath
	}
	if tmp.DatabaseDsn != "" {
		cfg.DBConfig = tmp.DatabaseDsn
	}
	cfg.HTTPSEnabled = tmp.EnableHttps // вот тут вопрос, что делать, ведь это инициализируется как фолс и будет перезатираться фолсом
}

func (c *Config) beautiPrint() {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(string(data))
}
