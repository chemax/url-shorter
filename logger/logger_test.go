package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewLogger(t *testing.T) {
	lg, err := NewLogger()
	assert.Nil(t, err)
	assert.NotNil(t, lg)
	defer func(lg *Log) {
		err = lg.Shutdown()
	}(lg)
	lg.Debug("debug")
	lg.Info("info")
	lg.Warn("warn")
	lg.Error("error")

}
