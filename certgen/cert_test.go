package certgen

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewCert(t *testing.T) {
	c1, c2 := NewCert()
	assert.NotNil(t, c1)
	assert.NotNil(t, c2)
}
