package storage

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/chemax/url-shorter/interfaces"
	"github.com/chemax/url-shorter/util"
	"os"
	"sync"
)

var newLineBytes = []byte("\n")

type URL struct {
	URL    string `json:"url"`
	Code   string `json:"code"`
	UserID string `json:"userID"`
}
type URLManager struct {
	db       interfaces.DBInterface
	URLs     map[string]*URL
	URLMx    sync.RWMutex
	SavePath string
	logger   interfaces.LoggerInterface
}

var manager = &URLManager{URLs: make(map[string]*URL)}

func Init(cfg interfaces.ConfigInterface, logger interfaces.LoggerInterface, db interfaces.DBInterface) (*URLManager, error) {
	if cfg.GetDBUse() {
		manager.db = db
	} else {
		manager.db = nil
	}
	manager.SavePath = cfg.GetSavePath()
	manager.logger = logger
	err := manager.restore()
	if err != nil {
		return nil, fmt.Errorf("restore err: %w", err)
	}
	return manager, nil
}
func (u *URLManager) restore() error {
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
		parsedURL := &URL{}
		err := json.Unmarshal(scanner.Bytes(), parsedURL)
		if err != nil {
			u.logger.Error("restore error: ", err.Error())
			continue
		}
		u.URLs[parsedURL.Code] = parsedURL
		u.logger.Debugln("restored: ", scanner.Text())
	}
	return nil
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

// Переиспользование во все поля. Нет четких критериев, что делать при сбое в середине процесса, я решил не прерывать,
// а значит, я могу переиспользовать имеющийся функционал.
func (u *URLManager) BatchSave(arr []*util.URLStructForBatch, httpPrefix string) (responseArr []util.URLStructForBatchResponse, err error) {
	var errorArr []error
	for _, v := range arr {
		shortcode, err := u.AddNewURL(v.OriginalURL, "")
		if err != nil {
			errorArr = append(errorArr, err)
			continue
		}
		responseArr = append(responseArr, util.URLStructForBatchResponse{
			CorrelationID: v.CorrelationID,
			ShortURL:      fmt.Sprintf("%s/%s", httpPrefix, shortcode),
		})
	}
	if len(errorArr) > 0 {
		return responseArr, fmt.Errorf("add URLs list errors: %w", errors.Join(errorArr...))
	}
	return responseArr, nil
}

func (u *URLManager) dbGetURL(code string) (parsedURL string, err error) {
	parsedURL, err = u.db.Get(code)
	if err != nil {
		return "", fmt.Errorf("get url from db error: %w", err)
	}
	return parsedURL, nil
}

func (u *URLManager) dbAddNewURL(parsedURL, userID string) (code string, err error) {
	var loop int
	for {
		//TODO переделать на функции внутри postgresql?
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

func (u *URLManager) GetUserURLs(userID string) (URLs []util.URLStructUser, err error) {
	if u.db != nil {
		return u.db.GetAllURLs(userID)
	}
	u.URLMx.RLock()
	defer u.URLMx.RUnlock()

	for _, v := range u.URLs {
		if v.UserID == userID {
			URLs = append(URLs, util.URLStructUser{Shortcode: v.Code, URL: v.URL})
		}
	}
	return URLs, err
}

func (u *URLManager) GetURL(code string) (parsedURL string, err error) {
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

func (u *URLManager) DeleteListFor(forDelete []string, userID string) {
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

func (u *URLManager) AddNewURL(parsedURL string, userID string) (code string, err error) {
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
	u.URLs[code] = &URL{URL: parsedURL, Code: code, UserID: userID}
	u.saveToFile(code)
	return code, nil
}

func (u *URLManager) Ping() bool {
	if u.db == nil {
		return false
	}
	err := u.db.Ping()
	if err != nil {
		u.logger.Error(fmt.Errorf("ping db error: %w", err))
		return false
	}
	return true
}
