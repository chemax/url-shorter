package config

import (
	"errors"
	"flag"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

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

type NetAddr struct {
	host string
	port int
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

type Config struct {
	NetAddr  *NetAddr
	HTTPAddr *HTTPAddr
}

var (
	cfg = &Config{
		NetAddr:  &NetAddr{host: "localhost", port: 8080},
		HTTPAddr: &HTTPAddr{addr: "http://localhost:8080"},
	}
)

func (c *Config) GetNetAddr() string {
	return c.NetAddr.String()
}
func (c *Config) GetHTTPAddr() string {
	return c.HTTPAddr.String()
}

func init() {

	flag.Var(cfg.NetAddr, "a", "Net address host:port")
	flag.Var(cfg.HTTPAddr, "b", "http(s) address http://host:port")
	flag.Parse()

}

func Get() *Config {
	return cfg
}
