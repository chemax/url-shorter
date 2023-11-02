package users

import (
	"context"
	"fmt"
	"github.com/chemax/url-shorter/interfaces"
	"github.com/chemax/url-shorter/util"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"net/http"
	"time"
)

type Users struct {
	SecretKey string
	TokenExp  time.Duration
	log       interfaces.LoggerInterface
	databaseInterface
}

type Claims struct {
	jwt.RegisteredClaims
	UserID string
}

var users = &Users{}

type databaseInterface interface {
	CreateUser() (string, error)
	Use() bool
}

func (u *Users) createNewUser(w http.ResponseWriter) (userID string, err error) {
	u.log.Debug("Create new user")
	if !u.Use() {
		userID = uuid.New().String()
	} else {
		userID, err = u.CreateUser()
	}
	if err != nil {
		return "", err
	}
	token, err := users.BuildJWTString(userID)
	if err != nil {
		return "", err
	}
	myCookie := &http.Cookie{
		Name:  util.TokenCookieName,
		Value: token,
	}
	http.SetCookie(w, myCookie)
	w.Header().Add("Authorization", fmt.Sprintf("Bearer %s", token))
	return userID, nil
}

func (u *Users) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var tkn *jwt.Token
		var userID string
		claims := &Claims{}
		c, err := r.Cookie(util.TokenCookieName)
		if err == nil {
			tknStr := c.Value
			tkn, err = jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (any, error) {
				return []byte(u.SecretKey), nil
			})
		}
		if err != nil || (tkn != nil && !tkn.Valid) {
			userID, err = u.createNewUser(w)
			if err != nil {
				users.log.Error(err)
				w.WriteHeader(http.StatusBadRequest)
			}
		} else {
			userID = claims.UserID
		}
		ctx := context.WithValue(r.Context(), util.UserID, userID)
		next.ServeHTTP(w, r.WithContext(ctx))

	})
}

// BuildJWTString создаёт токен и возвращает его в виде строки.
func (u *Users) BuildJWTString(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
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

func Init(cfg interfaces.ConfigInterface, log interfaces.LoggerInterface, dbObj databaseInterface) (*Users, error) {
	users.SecretKey = cfg.SecretKey()
	users.TokenExp = cfg.TokenExp()
	users.log = log
	users.databaseInterface = dbObj
	return users, nil
}
