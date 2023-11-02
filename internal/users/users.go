package users

import (
	"context"
	"github.com/chemax/url-shorter/interfaces"
	"github.com/chemax/url-shorter/internal/config"
	"github.com/chemax/url-shorter/internal/logger"
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
		uuid.New().String()
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
	http.SetCookie(w, &http.Cookie{
		Name:  "token",
		Value: token,
	})
	return userID, nil
}

func (u *Users) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var tkn *jwt.Token
		var userID string
		claims := &Claims{}
		c, err := r.Cookie("token")
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
		ctx := context.WithValue(r.Context(), "userID", userID)
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

func Init(cfg *config.Config, log *logger.Logger, dbObj databaseInterface) (*Users, error) {
	users.SecretKey = cfg.SecretKey
	users.TokenExp = cfg.TokenExp
	users.log = log
	users.databaseInterface = dbObj
	return users, nil
}