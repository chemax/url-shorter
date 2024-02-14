package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewLogger(t *testing.T) {
	lg, err := NewLogger()
	assert.Nil(t, err)
	assert.NotNil(t, lg)
	defer func(lg *log) {
		err = lg.Shutdown()
	}(lg)

}
