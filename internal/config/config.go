package config

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/chemax/url-shorter/util"
)

var (
	cfg = &Config{
		NetAddr:   &NetAddr{Host: "localhost", Port: 8080},
		HTTPAddr:  &HTTPAddr{Addr: "http://localhost:8080"},
		PathSave:  &PathForSave{path: "/tmp/short-url-db.json"},
		DBConfig:  &DBConfig{connectString: ""},
		tokenExp:  time.Hour * 3,
		secretKey: "XXXXXX",
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
		err := cfg.PathSave.Set(savePath)
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
