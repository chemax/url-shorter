package storage

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/chemax/url-shorter/util"
)

var newLineBytes = []byte("\n")

// Configer интерфейс конфиг-структуры
type Configer interface {
	GetSavePath() string
	GetDBUse() bool
}

// Loggerer интерфейс логера
type Loggerer interface {
	Debugln(args ...interface{})
	Error(args ...interface{})
}

// DataBaser интерфейс для базы данных
type DataBaser interface {
	BatchDelete([]string, string)
	Ping() error
	SaveURL(code string, URL string, userID string) (string, error)
	Get(code string) (string, error)
	GetAllURLs(userID string) ([]util.URLWithShort, error)
	Use() bool
}

type singleURL struct {
	URL    string `json:"url"`
	Code   string `json:"code"`
	UserID string `json:"userID"`
}
type managerURL struct {
	db       DataBaser
	URLs     map[string]*singleURL
	URLMx    sync.RWMutex
	SavePath string
	log      Loggerer
}

var manager = &managerURL{URLs: make(map[string]*singleURL)}

// Init создает и возвращает структуру управления URL'ами
func Init(cfg Configer, logger Loggerer, db DataBaser) (*managerURL, error) {
	if cfg.GetDBUse() {
		manager.db = db
	} else {
		manager.db = nil
	}
	manager.SavePath = cfg.GetSavePath()
	manager.log = logger
	err := manager.restore()
	if err != nil {
		return nil, fmt.Errorf("restore err: %w", err)
	}
	return manager, nil
}

func (u *managerURL) restore() error {
	if u.db != nil {
		return nil
	}
	if u.SavePath == "" {
		return nil
	}
	file, err := os.OpenFile(u.SavePath, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return fmt.Errorf("error restore db: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		parsedURL := &singleURL{}
		err := json.Unmarshal(scanner.Bytes(), parsedURL)
		if err != nil {
			u.log.Error("restore error: ", err.Error())
			continue
		}
		u.URLs[parsedURL.Code] = parsedURL
		u.log.Debugln("restored: ", scanner.Text())
	}
	return nil
}
func (u *managerURL) saveToFile(code string) {
	if u.SavePath == "" {
		return
	}
	file, err := os.OpenFile(u.SavePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		u.log.Error(fmt.Sprintf("error open file [%s], error: %s", u.SavePath, err.Error()))
		return
	}
	defer file.Close()
	var data []byte
	data, err = json.Marshal(u.URLs[code])
	if err != nil {
		u.log.Error(fmt.Sprintf("unmarshal error: %s", err.Error()))
		return
	}
	data = append(data, newLineBytes...)
	_, err = file.Write(data)
	if err != nil {
		u.log.Error(fmt.Sprintf("error write to file [%s], error: %s", u.SavePath, err.Error()))
		return
	}

}

// BatchSave пакетное сохранение
func (u *managerURL) BatchSave(arr []*util.URLForBatch, httpPrefix string) (responseArr []util.URLForBatchResponse, err error) {
	var errorArr []error
	for _, v := range arr {
		shortcode, err := u.AddNewURL(v.OriginalURL, "")
		if err != nil {
			errorArr = append(errorArr, err)
			continue
		}
		responseArr = append(responseArr, util.URLForBatchResponse{
			CorrelationID: v.CorrelationID,
			ShortURL:      fmt.Sprintf("%s/%s", httpPrefix, shortcode),
		})
	}
	if len(errorArr) > 0 {
		return responseArr, fmt.Errorf("add URLs list errors: %w", errors.Join(errorArr...))
	}
	return responseArr, nil
}

func (u *managerURL) dbGetURL(code string) (parsedURL string, err error) {
	parsedURL, err = u.db.Get(code)
	if err != nil {
		return "", fmt.Errorf("get url from db error: %w", err)
	}
	return parsedURL, nil
}

func (u *managerURL) dbAddNewURL(parsedURL, userID string) (code string, err error) {
	var loop int
	for {
		code = util.RandStringRunes(util.CodeLength)
		dupCode, err := u.db.SaveURL(code, parsedURL, userID)
		if err != nil && !errors.Is(err, &util.AlreadyHaveThisURLError{}) {
			loop++
			if loop > util.CodeGenerateAttempts {
				code = ""
				return code, fmt.Errorf("can not found free code for short url")
			}
			continue
		}
		if errors.Is(err, &util.AlreadyHaveThisURLError{}) {
			return dupCode, err
		}
		return code, nil
	}
}

// GetUserURLs вернуть все URL пользователя
func (u *managerURL) GetUserURLs(userID string) (URLs []util.URLWithShort, err error) {
	if u.db != nil {
		return u.db.GetAllURLs(userID)
	}
	u.URLMx.RLock()
	defer u.URLMx.RUnlock()

	for _, v := range u.URLs {
		if v.UserID == userID {
			URLs = append(URLs, util.URLWithShort{Shortcode: v.Code, URL: v.URL})
		}
	}
	return URLs, err
}

// GetURL получить URL по коду
func (u *managerURL) GetURL(code string) (parsedURL string, err error) {
	if u.db != nil {
		return u.dbGetURL(code)
	}
	u.URLMx.RLock()
	defer u.URLMx.RUnlock()
	urlObj, ok := u.URLs[code]
	if !ok {
		return "", fmt.Errorf("requested url not found")
	}
	return urlObj.URL, nil
}

// DeleteListFor пакетное удаление
func (u *managerURL) DeleteListFor(forDelete []string, userID string) {
	if u.db != nil {
		u.db.BatchDelete(forDelete, userID)
		return
	}
	for _, v := range forDelete {
		u.URLMx.Lock()
		_, ok := u.URLs[v]
		if ok && u.URLs[v].UserID == userID {
			delete(u.URLs, v)
		}
	}
}

// AddNewURL сохранить URL
func (u *managerURL) AddNewURL(parsedURL string, userID string) (code string, err error) {
	if u.db != nil {
		return u.dbAddNewURL(parsedURL, userID)
	}
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
	u.URLs[code] = &singleURL{URL: parsedURL, Code: code, UserID: userID}
	u.saveToFile(code)
	return code, nil
}

// Ping пинг бд
func (u *managerURL) Ping() bool {
	if u.db == nil {
		return false
	}
	err := u.db.Ping()
	if err != nil {
		u.log.Error(fmt.Errorf("ping db error: %w", err))
		return false
	}
	return true
}
