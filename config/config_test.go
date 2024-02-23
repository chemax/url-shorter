package config

import (
	"fmt"
	"os"
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

func TestConfig_GetNetAddr(t *testing.T) {
	cfgForTest := &Config{
		NetAddr: "localhost:8080",
	}

	result := cfgForTest.GetNetAddr()
	expected := "localhost:8080"

	if result != expected {
		t.Errorf("unexpected result: %s", result)
	}
}

func TestConfig_GetHTTPAddr(t *testing.T) {
	cfgForTest := &Config{
		HTTPAddr: "http://localhost:8080",
	}

	result := cfgForTest.GetHTTPAddr()
	expected := "http://localhost:8080"

	if result != expected {
		t.Errorf("unexpected result: %s", result)
	}
}

func TestInit(t *testing.T) {
	//проверяем правильность приоритетов конфигурации (какая же дурацкая задача)
	tFile, err := os.CreateTemp("", "config.json")
	assert.Nil(t, err)
	// Вызов ответственен за очистку
	defer func(name string) {
		err := tFile.Close()
		if err != nil {
			fmt.Println("error close tmp file", err.Error())
		}
		err = os.Remove(name)
		if err != nil {
			fmt.Println("error delete tmp file", err.Error())
		}
	}(tFile.Name())
	testJSONData := []byte(`{
	   "server_address": "localhost:9999",
	   "base_url": "http://localhost:54321",
	   "file_storage_path": "/some.db",
	   "database_dsn": "",
	   "enable_https": true
	}`)

	_, err = tFile.Write(testJSONData)
	assert.Nil(t, err)

	os.Setenv("SERVER_ADDRESS", "localhost:4444")
	os.Setenv("BASE_URL", "http://localhost:8080")
	os.Setenv("CONFIG", tFile.Name())

	//Если при наличии энв переменной полезет в несуществующий файл - чет не так.
	os.Args = []string{"somebinaryfile", "-c", "./config_test.json"}

	defer func() {
		os.Unsetenv("SERVER_ADDRESS")
		os.Unsetenv("BASE_URL")
		os.Unsetenv("CONFIG")
	}()

	cfgFirst, err := NewConfig()
	assert.Nil(t, err)
	assert.NotNil(t, cfgFirst)
	assert.Equal(t, "localhost:4444", cfgFirst.GetNetAddr())
	assert.Equal(t, "/some.db", cfgFirst.GetSavePath())
}
