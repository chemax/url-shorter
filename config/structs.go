package config

/*
In this file only type structs and its methods
*/
import (
	"errors"
	"flag"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// Config содержит в себе весь конфиг включая подструктуры
type Config struct {
	NetAddr         *NetAddr
	HTTPAddr        *HTTPAddr
	PathSave        *PathForSave
	DBConfig        *DBConfig
	flagInitialized bool
	secretKey       string
	tokenExp        time.Duration
}

// DBConfig конфиг базы данных
type DBConfig struct {
	connectString string
}

// Set сеттер коннект стринга базы данных
func (p *DBConfig) Set(s string) error {
	p.connectString = s
	return nil
}

// String стрингует коннект стринг базы данных (для реализации интерфейса стрингера)
func (p *DBConfig) String() string {
	return p.connectString
}

// PathForSave путь сохранения "текстовой" базы сокращений
type PathForSave struct {
	path string
}

// Set implements pkg "flag" interface
func (p *PathForSave) Set(s string) error {
	p.path = s
	return nil
}

// String стрингует коннект стринг файлового хранилища (для реализации интерфейса стрингера)
func (p *PathForSave) String() string {
	return p.path
}

// HTTPAddr хранит http адрес приложения, никак формально не связан с сетевым адресом, используется для генерации сокращенных ссылок.
type HTTPAddr struct {
	Addr string `json:"addr"`
}

// String для реализации интерфейса стрингера
func (h *HTTPAddr) String() string {
	return h.Addr
}

// Set implements pkg "flag" interface
func (h *HTTPAddr) Set(s string) error {
	_, err := url.ParseRequestURI(s)
	if err != nil {
		return fmt.Errorf("parse base http addr error: %w", err)
	}
	s = strings.TrimSuffix(s, "/")
	h.Addr = s
	return nil
}

// NetAddr хранит хост и порт которые биндит приложение для http сервера
type NetAddr struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

// String возвращает хост:порт для бинда http-сервером
func (a *NetAddr) String() string {
	return fmt.Sprintf("%s:%d", a.Host, a.Port)
}

// Set implements pkg "flag" interface
func (a *NetAddr) Set(s string) error {
	hp := strings.Split(s, ":")
	if len(hp) != 2 {
		return errors.New("need address in a form Host:Port")
	}
	port, err := strconv.Atoi(hp[1])
	if err != nil {
		return err
	}
	a.Host = hp[0]
	a.Port = port
	return nil
}

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
	return c.DBConfig.connectString != ""
}

// GetSavePath путь к текстовой базе данных
func (c *Config) GetSavePath() string {
	return c.PathSave.path
}

// GetNetAddr получить сетевой адрес
func (c *Config) GetNetAddr() string {
	return c.NetAddr.String()
}

// GetHTTPAddr возвращает http адрес сервиса (для генерации сокращенных ссылок)
func (c *Config) GetHTTPAddr() string {
	return c.HTTPAddr.String()
}

func (c *Config) initFlags() {
	if c.flagInitialized {
		return
	}
	c.flagInitialized = true
	flag.Var(cfg.NetAddr, "a", "Net address Host:Port")
	flag.Var(cfg.HTTPAddr, "b", "http(s) address http://host:port")
	flag.Var(cfg.PathSave, "f", "full path to file for save url's")
	flag.Var(cfg.DBConfig, "d", "DB connect string like \"postgres://username:password@localhost:5432/database_name\"")
}
