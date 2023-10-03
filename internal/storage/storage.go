package storage

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/chemax/url-shorter/util"
	"os"
	"sync"
)

var newLineBytes = []byte("\n")

type URL struct {
	URL  string `json:"url"`
	Code string `json:"code"`
}
type URLManager struct {
	URLs     map[string]*URL
	URLMx    sync.RWMutex
	SavePath string
	logger   util.LoggerInterface
}

var manager = &URLManager{URLs: make(map[string]*URL)}

func Get(savePath string, logger util.LoggerInterface) *URLManager {
	manager.SavePath = savePath
	manager.logger = logger
	once := sync.Once{}
	once.Do(manager.restore)
	return manager
}
func (u *URLManager) restore() {
	if u.SavePath == "" {
		return
	}
	file, err := os.OpenFile(u.SavePath, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		u.logger.Error("error restore db", err.Error())
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		parsedURL := &URL{}
		err := json.Unmarshal(scanner.Bytes(), parsedURL)
		if err != nil {
			u.logger.Error("restore error: ", err.Error())
			continue
		}
		u.URLs[parsedURL.Code] = parsedURL
		u.logger.Debug("restored: ", scanner.Text())
	}
}

func (u *URLManager) GetURL(code string) (parsedURL string, err error) {
	u.URLMx.RLock()
	defer u.URLMx.RUnlock()
	urlObj, ok := u.URLs[code]
	if !ok {
		return "", fmt.Errorf("requested url not found")
	}
	return urlObj.URL, nil
}

func (u *URLManager) saveToFile(code string) {
	if u.SavePath == "" {
		return
	}
	file, err := os.OpenFile(u.SavePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		u.logger.Error(fmt.Sprintf("error open file [%s], error: %s", u.SavePath, err.Error()))
		return
	}
	defer file.Close()
	var data []byte
	data, err = json.Marshal(u.URLs[code])
	if err != nil {
		u.logger.Error(fmt.Sprintf("unmarshal error: %s", err.Error()))
		return
	}
	data = append(data, newLineBytes...)
	_, err = file.Write(data)
	if err != nil {
		u.logger.Error(fmt.Sprintf("error write to file [%s], error: %s", u.SavePath, err.Error()))
		return
	}

}

func (u *URLManager) AddNewURL(parsedURL string) (code string, err error) {
	var ok = true
	var loop int
	u.URLMx.Lock()
	defer u.URLMx.Unlock()
	for ok {
		code = util.RandStringRunes(util.CodeLength)
		_, ok = u.URLs[code]
		loop++
		if loop > util.CodeGenerateAttempts {
			code = ""
			return code, fmt.Errorf("can not found free code for short url")
		}
	}
	u.URLs[code] = &URL{URL: parsedURL, Code: code}
	u.saveToFile(code)
	return code, err
}
