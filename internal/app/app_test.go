package app

import (
	"fmt"
	"github.com/chemax/url-shorter/models"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func TestApp(t *testing.T) {
	go func() {
		// Я не знаю как ещё это протестировать
		t := time.NewTicker(time.Second * 10)
		<-t.C
		os.Exit(0)
	}()
	t.Run("ServerAddressEnv", func(t *testing.T) {

		os.Setenv(models.ServerAddressEnv, "fdsfffwefew")
		err := Run()
		assert.Error(t, fmt.Errorf("need address in a form Host:Port"), err)
		os.Unsetenv(models.ServerAddressEnv)
	})
	t.Run("BaseURLEnv", func(t *testing.T) {
		os.Setenv(models.BaseURLEnv, "http://www.example.com:xxxx//////path")
		err := Run()
		assert.Error(t, err)
		os.Unsetenv(models.BaseURLEnv)
	})
}
