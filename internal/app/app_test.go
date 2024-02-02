package app

import (
	"fmt"
	"os"
	"testing"

	"github.com/chemax/url-shorter/util"
	"github.com/stretchr/testify/assert"
)

func TestApp(t *testing.T) {
	t.Run("ServerAddressEnv", func(t *testing.T) {
		os.Setenv(util.ServerAddressEnv, "fdsfffwefew")
		err := Run()
		assert.Error(t, fmt.Errorf("need address in a form Host:Port"), err)
		os.Unsetenv(util.ServerAddressEnv)
	})
	t.Run("BaseURLEnv", func(t *testing.T) {
		os.Setenv(util.BaseURLEnv, "http://www.example.com:xxxx//////path")
		err := Run()
		assert.Error(t, err)
		os.Unsetenv(util.BaseURLEnv)
	})
}
