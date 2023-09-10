package config

import (
	"errors"
	"flag"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
)

type Config struct {
	NetAddr  *NetAddr
	HTTPAddr *HTTPAddr
}

const serverAddressEnv = "SERVER_ADDRESS"
const baseUrlEnv = "BASE_URL"

var (
	cfg = &Config{
		NetAddr:  &NetAddr{host: "localhost", port: 8080},
		HTTPAddr: &HTTPAddr{addr: "http://localhost:8080"},
	}
)

type NetAddr struct {
	host string
	port int
}

type HTTPAddr struct {
	addr string
}

func (h HTTPAddr) String() string {
	return h.addr
}

func (h *HTTPAddr) Set(s string) error {
	_, err := url.ParseRequestURI(s)
	if err != nil {
		return err
	}
	s = strings.TrimSuffix(s, "/")
	h.addr = s
	return nil
}

func (a NetAddr) String() string {
	return fmt.Sprintf("%s:%d", a.host, a.port)
}

func (a *NetAddr) Set(s string) error {
	hp := strings.Split(s, ":")
	if len(hp) != 2 {
		return errors.New("need address in a form host:port")
	}
	port, err := strconv.Atoi(hp[1])
	if err != nil {
		return err
	}
	a.host = hp[0]
	a.port = port
	return nil
}

func (c *Config) GetNetAddr() string {
	return c.NetAddr.String()
}
func (c *Config) GetHTTPAddr() string {
	return c.HTTPAddr.String()
}

func MustConfig() {
	if srvAddr, ok := os.LookupEnv(serverAddressEnv); ok {
		err := cfg.NetAddr.Set(srvAddr)
		if err != nil {
			panic(err)
		}
	} else {
		flag.Var(cfg.NetAddr, "a", "Net address host:port")
	}
	if baseUrl, ok := os.LookupEnv(baseUrlEnv); ok {
		err := cfg.HTTPAddr.Set(baseUrl)
		if err != nil {
			panic(err)
		}
	} else {
		flag.Var(cfg.HTTPAddr, "b", "http(s) address http://host:port")
	}
	flag.Parse()
}

func Get() *Config {
	var once sync.Once
	once.Do(MustConfig)
	return cfg
}
