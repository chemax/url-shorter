package config

import (
	"flag"
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	t.Run("err is nil", func(t *testing.T) {
		config, err := NewConfig()
		assert.Nil(t, err)
		assert.NotNil(t, config)
	})

}

func TestNetAddr_Set(t *testing.T) {
	flag.Parse()
	addr := &NetAddr{}

	err := addr.Set("localhost:8080")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if addr.Host != "localhost" || addr.Port != 8080 {
		t.Errorf("unexpected result: %s:%d", addr.Host, addr.Port)
	}

	err = addr.Set("invalid_address")
	if err == nil {
		t.Error("expected an error")
	}
}

func TestHTTPAddr_Set(t *testing.T) {
	httpAddr := &HTTPAddr{}

	err := httpAddr.Set("http://localhost:8080")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if httpAddr.Addr != "http://localhost:8080" {
		t.Errorf("unexpected result: %s", httpAddr.Addr)
	}

	err = httpAddr.Set("invalid_address")
	if err == nil {
		t.Error("expected an error")
	}
}

func TestConfig_GetNetAddr(t *testing.T) {
	cfgForTest := &Config{
		NetAddr: &NetAddr{Host: "localhost", Port: 8080},
	}

	result := cfgForTest.GetNetAddr()
	expected := "localhost:8080"

	if result != expected {
		t.Errorf("unexpected result: %s", result)
	}
}

func TestConfig_GetHTTPAddr(t *testing.T) {
	cfgForTest := &Config{
		HTTPAddr: &HTTPAddr{Addr: "http://localhost:8080"},
	}

	result := cfgForTest.GetHTTPAddr()
	expected := "http://localhost:8080"

	if result != expected {
		t.Errorf("unexpected result: %s", result)
	}
}

func TestInit(t *testing.T) {
	os.Setenv("SERVER_ADDRESS", "localhost:8080")
	os.Setenv("BASE_URL", "http://localhost:8080")

	defer func() {
		os.Unsetenv("SERVER_ADDRESS")
		os.Unsetenv("BASE_URL")
	}()

	//init()

	expectedNetAddr := &NetAddr{Host: "localhost", Port: 8080}
	if !reflect.DeepEqual(cfg.NetAddr, expectedNetAddr) {
		t.Errorf("unexpected NetAddr: %+v", cfg.NetAddr)
	}

	expectedHTTPAddr := &HTTPAddr{Addr: "http://localhost:8080"}
	if !reflect.DeepEqual(cfg.HTTPAddr, expectedHTTPAddr) {
		t.Errorf("unexpected HTTPAddr: %+v", cfg.HTTPAddr)
	}
}
