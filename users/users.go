package users

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/chemax/url-shorter/models"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

// configer интерфейс конфиг-структуры
type configer interface {
	SecretKey() string
	TokenExp() time.Duration
}

// loggerer интерфейс логера
type loggerer interface {
	Debugln(args ...interface{})
	Error(args ...interface{})
}

type users struct {
	SecretKey string
	TokenExp  time.Duration
	log       loggerer
	dataBaser
}

type claimsStruct struct {
	jwt.RegisteredClaims
	UserID string
}

var usersManager = &users{}

type dataBaser interface {
	CreateUser() (string, error)
	Use() bool
}

func (u *users) createNewUser(w http.ResponseWriter) (userID string, err error) {
	u.log.Debugln("Create new user")
	if !u.Use() {
		userID = uuid.New().String()
	} else {
		userID, err = u.CreateUser()
	}
	if err != nil {
		return "", err
	}
	token, err := usersManager.BuildJWTString(userID)
	if err != nil {
		return "", err
	}
	myCookie := &http.Cookie{
		Name:  models.TokenCookieName,
		Value: token,
	}
	http.SetCookie(w, myCookie)
	w.Header().Add("Authorization", fmt.Sprintf("Bearer %s", token))
	return userID, nil
}

// Middleware для аутентификации и авторизации
func (u *users) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var tkn *jwt.Token
		var userID string
		claims := &claimsStruct{}
		c, err := r.Cookie(models.TokenCookieName)
		if err == nil {
			tknStr := c.Value
			tkn, err = jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (any, error) {
				return []byte(u.SecretKey), nil
			})
		}
		if err != nil && r.URL.String() == "/api/user/urls" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if err != nil || (tkn != nil && !tkn.Valid) {
			userID, err = u.createNewUser(w)
			if err != nil {
				usersManager.log.Error(err)
				w.WriteHeader(http.StatusBadRequest)
			}
		} else {
			userID = claims.UserID
		}
		ctx := context.WithValue(r.Context(), models.UserID, userID)
		next.ServeHTTP(w, r.WithContext(ctx))

	})
}

// BuildJWTString создаёт токен и возвращает его в виде строки.
func (u *users) BuildJWTString(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claimsStruct{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(u.TokenExp)),
		},
		UserID: userID,
	})

	tokenString, err := token.SignedString([]byte(u.SecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// NewUser возвращает юзер менеджера
func NewUser(cfg configer, log loggerer, dbObj dataBaser) (*users, error) {
	usersManager.SecretKey = cfg.SecretKey()
	usersManager.TokenExp = cfg.TokenExp()
	usersManager.log = log
	usersManager.dataBaser = dbObj
	return usersManager, nil
}
