package config

import (
	"errors"
	"flag"
	"fmt"
	"github.com/chemax/url-shorter/util"
	"net/url"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	NetAddr  *NetAddr
	HTTPAddr *HTTPAddr
	SavePath *PathForSave
}

var (
	cfg = &Config{
		NetAddr:  &NetAddr{Host: "localhost", Port: 8080},
		HTTPAddr: &HTTPAddr{Addr: "http://localhost:8080"},
		SavePath: &PathForSave{path: "/tmp/short-url-db.json"},
	}
)

type PathForSave struct {
	path string
}

func (p *PathForSave) Set(s string) error {
	p.path = s
	return nil
}

func (p PathForSave) String() string {
	return p.path
}

type NetAddr struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

type HTTPAddr struct {
	Addr string `json:"addr"`
}

func (h HTTPAddr) String() string {
	return h.Addr
}

func (h *HTTPAddr) Set(s string) error {
	_, err := url.ParseRequestURI(s)
	if err != nil {
		return err
	}
	s = strings.TrimSuffix(s, "/")
	h.Addr = s
	return nil
}

func (a NetAddr) String() string {
	return fmt.Sprintf("%s:%d", a.Host, a.Port)
}

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

func (c *Config) GetNetAddr() string {
	return c.NetAddr.String()
}
func (c *Config) GetHTTPAddr() string {
	return c.HTTPAddr.String()
}

func Init() (*Config, error){
	flag.Var(cfg.NetAddr, "a", "Net address Host:Port")
	flag.Var(cfg.HTTPAddr, "b", "http(s) address http://host:port")
	flag.Var(cfg.SavePath, "f", "full path to file for save url's")
	flag.Parse()

	if srvAddr, ok := os.LookupEnv(util.ServerAddressEnv); ok && srvAddr != "" {
		err := cfg.NetAddr.Set(srvAddr)
		if err != nil {
			return nil, fmt.Errorf("error setup server address: %w",err)
		}
	}
	if baseURL, ok := os.LookupEnv(util.BaseURLEnv); ok && baseURL != "" {
		err := cfg.HTTPAddr.Set(baseURL)
		if err != nil {
			return nil, fmt.Errorf("error setup base url: %w",err)
		}
	}
	if savePath, ok := os.LookupEnv(util.SavePath); ok && savePath != "" {
		err := cfg.SavePath.Set(savePath)
		if err != nil {
			return nil, fmt.Errorf("error setup save path: %w",err)
		}
	}
	return cfg, nil
}


