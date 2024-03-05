// Package config служит хранилищем конфигурации приложения, получает данные из окружения, из аргументов и имеет дефолтные значения для некоторых параметров.
// реализует паттерн singleton.
package config

import (
	"os"
	"time"

	"github.com/chemax/url-shorter/models"
)

var (
	cfg = &Config{
		NetAddr:       "localhost:8080",
		HTTPAddr:      "http://localhost:8080",
		PathSave:      "/tmp/short-url-db.json",
		DBConfig:      "",
		tokenExp:      time.Hour * 3,
		secretKey:     "XXXXXX",
		TrustedSubnet: "",
	}
)

// NewConfig инициализация конфига
// мне не нравится работа ради работы даже в учебных проектах.
// такой тупой фигней как 100500 вариантов конфигурации сервиса никто в здравом уме делать не будет.
// это понижает надежность сервиса в разы. 100500 способов отстрелить себе ногу, привет Си.
func NewConfig() (*Config, error) {
	//defer cfg.beautiPrint()
	tmpConfigForFlags := &tmpConfig{} //сюда мы запишем данные из флагов, чтобы не затереть их потом конфигом из JSON
	cfg.initFlags(tmpConfigForFlags)  // прихраним их и спарсим жсончи
	// нам надо определить есть ли у нас конфигФайл. Сначала проверяем энв.
	cfgPath := os.Getenv(models.CONFIG)
	if cfgPath == "" {
		cfgPath = tmpConfigForFlags.ConfigPath
	}
	// выглядит корявенько, конечно.
	if cfgPath != "" {
		fileJSONData, err := os.ReadFile(cfgPath)
		if err != nil {
			return nil, err
		}
		err = cfg.parseJSON(fileJSONData)
		if err != nil {
			return nil, err
		}
		// Данные из файла подтянуты, теперь можно флаги записать.
	}
	cfg.SetFromTmpConfig(tmpConfigForFlags)
	// а теперь по классике, там где есть энв переменная, мы её считываем.
	if srvAddr, ok := os.LookupEnv(models.ServerAddressEnv); ok && srvAddr != "" {
		cfg.NetAddr = srvAddr
	}
	if baseURL, ok := os.LookupEnv(models.BaseURLEnv); ok && baseURL != "" {
		cfg.HTTPAddr = baseURL
	}
	if savePath, ok := os.LookupEnv(models.SavePath); ok && savePath != "" {
		cfg.PathSave = savePath
	}
	if connectString, ok := os.LookupEnv(models.DBConnectString); ok {
		cfg.DBConfig = connectString
	}
	if trustedSubnet, ok := os.LookupEnv(models.TrustedSubnet); ok {
		cfg.TrustedSubnet = trustedSubnet
	}
	if _, ok := os.LookupEnv(models.HTTPSEnabled); ok {
		cfg.HTTPSEnabled = true
	}
	return cfg, nil
}
