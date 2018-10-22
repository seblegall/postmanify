package postmanify

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConverter(t *testing.T) {
	conv := NewConverter(Config{})

	assert.NotNil(t, conv)
}
