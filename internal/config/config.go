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
		DBConfig: &DBConfig{connectString: ""},
	}
)

func Init() (*Config, error) {
	cfg.initFlags()
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
	if connectString, ok := os.LookupEnv(util.DBConnectString); ok {
		err := cfg.DBConfig.Set(connectString)
		if err != nil {
			return nil, fmt.Errorf("error setup save path: %w", err)
		}
	}
	return cfg, nil
}
