package app

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestApp(t *testing.T) {
	//ctrl := gomock.NewController(t)
	//defer ctrl.Finish()
	t.Run("error db", func(t *testing.T) {
		err := Run()
		assert.Equal(t, err.Error(), "error init config: db connect string is empty")
	})
}
