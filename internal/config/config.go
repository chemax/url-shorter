package config

import (
	"flag"
	"fmt"
	"github.com/chemax/url-shorter/util"
	"os"
)

var (
	cfg = &Config{
		NetAddr:  &NetAddr{Host: "localhost", Port: 8080},
		HTTPAddr: &HTTPAddr{Addr: "http://localhost:8080"},
		SavePath: &PathForSave{path: "/tmp/short-url-db.json"},
		DBConfig: &DBConfig{},
	}
)

func Init() (*Config, error) {
	flag.Var(cfg.NetAddr, "a", "Net address Host:Port")
	flag.Var(cfg.HTTPAddr, "b", "http(s) address http://host:port")
	flag.Var(cfg.SavePath, "f", "full path to file for save url's")
	flag.Var(cfg.DBConfig, "d", "DB connect string like \"postgres://username:password@localhost:5432/database_name\"")
	flag.Parse()

	if srvAddr, ok := os.LookupEnv(util.ServerAddressEnv); ok && srvAddr != "" {
		err := cfg.NetAddr.Set(srvAddr)
		if err != nil {
			return nil, fmt.Errorf("error setup server address: %w", err)
		}
	}
	if baseURL, ok := os.LookupEnv(util.BaseURLEnv); ok && baseURL != "" {
		err := cfg.HTTPAddr.Set(baseURL)
		if err != nil {
			return nil, fmt.Errorf("error setup base url: %w", err)
		}
	}
	if savePath, ok := os.LookupEnv(util.SavePath); ok && savePath != "" {
		err := cfg.SavePath.Set(savePath)
		if err != nil {
			return nil, fmt.Errorf("error setup save path: %w", err)
		}
	}
	if connectString, ok := os.LookupEnv(util.DBConnectString); ok && connectString != "" {
		err := cfg.DBConfig.Set(connectString)
		if err != nil {
			return nil, fmt.Errorf("error setup save path: %w", err)
		}
	}
	if cfg.DBConfig.String() == "" {
		return nil, fmt.Errorf("db connect string is empty")
	}
	return cfg, nil
}
