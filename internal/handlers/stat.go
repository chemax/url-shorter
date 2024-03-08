package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/chemax/url-shorter/models"
	"net"
	"net/http"
)

// DeleteUserURLsHandler ручка для удаления пользовательских урл пачкой
func (h *handlers) statHandler(res http.ResponseWriter, r *http.Request) {
	if !h.checkIP(r) { // если такая авторизация нужна по нескольким "направлениям", это стоит вынести в мидлварю, но это надо делать реальную систему правил, а у нас четкое ТЗ. Сойдет.
		h.Log.Warn("unauthorized access to stats!")
		res.WriteHeader(http.StatusForbidden)
		return
	}
	var err error
	defer func() {
		if err != nil {
			h.Log.Warn(err.Error())
			res.WriteHeader(http.StatusBadRequest)
		}
	}()
	stats, err := h.storage.GetStats()
	if err != nil {
		err = fmt.Errorf("request service stat error: %w", err)
	}
	data, err := json.Marshal(stats)
	if err != nil { // Ну не может моя структура маршалится с ошибкой, ну бред же, а не проверять нельзя - линтер забреет.
		h.Log.Error(err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	_, err = res.Write(data)
	if err != nil {
		h.Log.Warn("response write error: ", err.Error())
		err = nil
	}
}

func (h *handlers) checkIP(res *http.Request) bool {
	realIP := net.ParseIP(res.Header.Get(models.RealIP))
	return realIP != nil && h.TrustedSubnet.Contains(realIP)
}
