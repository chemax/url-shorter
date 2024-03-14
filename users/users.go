// Package users менеджер пользователей
package users

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/chemax/url-shorter/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/chemax/url-shorter/models"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

// Configer интерфейс конфиг-структуры
type Configer interface {
	SecretKey() string
	TokenExp() time.Duration
}

// Loggerer интерфейс логера
type Loggerer interface {
	Warn(args ...interface{})
	Debug(args ...interface{})
	Info(args ...interface{})
	Warnln(args ...interface{})
	Error(args ...interface{})
	Errorln(args ...interface{})
}

type users struct {
	SecretKey string
	TokenExp  time.Duration
	log       Loggerer
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

func (u *users) createUser() (userID, JWTToken string, err error) {
	u.log.Debug("Create new user")
	if !u.Use() {
		userID = uuid.New().String()
	} else {
		userID, err = u.CreateUser()
	}
	if err != nil {
		return "", "", err
	}
	token, err := usersManager.BuildJWTString(userID)
	return userID, token, err
}

func (u *users) writeUserToResponse(w http.ResponseWriter) (userID string, err error) {
	userID, token, err := usersManager.createUser()
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
			userID, err = u.writeUserToResponse(w)
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

// GetFromContextGRPC gets jwt token from context.
func GetFromContextGRPC(ctx context.Context) (string, error) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		values := md.Get("access-token")
		if len(values) > 0 {
			return values[0], nil
		}
	}

	return "", fmt.Errorf("no data")
}

func (u *users) JWTInterceptor(log *logger.Log) grpc.UnaryServerInterceptor {
	log.Debug("JWT interceptor enabled")
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		var tkn *jwt.Token
		var userID string
		var token string
		claims := &claimsStruct{}

		token, err := GetFromContextGRPC(ctx)
		if err != nil {
			userID, token, err = u.createUser() //TODO встроить в контекст нового пользователя
			if err != nil {
				usersManager.log.Error(err)
				return nil, err
			}
		}
		tkn, err = jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (any, error) {
			return []byte(u.SecretKey), nil
		})
		if err != nil || (tkn != nil && tkn.Valid) {
			userID, token, err = u.createUser() //TODO встроить в контекст нового пользователя
			if err != nil {
				usersManager.log.Error(err)
				return nil, err
			}
		}
		newCTX := context.WithValue(ctx, models.UserID, userID)
		newCTX = context.WithValue(newCTX, "access-token", token)
		resp, err := handler(newCTX, req)

		return resp, err
	}
}

// NewUser возвращает юзер менеджера
func NewUser(cfg Configer, log Loggerer, dbObj dataBaser) (*users, error) {
	usersManager.SecretKey = cfg.SecretKey()
	usersManager.TokenExp = cfg.TokenExp()
	usersManager.log = log
	usersManager.dataBaser = dbObj
	return usersManager, nil
}
