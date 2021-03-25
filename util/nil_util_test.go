package util

import (
	"testing"

	"github.com/bmizerany/assert"
)

func TestIsNil(t *testing.T) {
	assert.Equal(t, true, IsNil(""))
	assert.Equal(t, false, IsNil(1))
	assert.Equal(t, true, IsNil(0))
	assert.Equal(t, false, IsNil(struct{}{}))
	assert.Equal(t, true, IsNil(nil))
}

func TestHasNil(t *testing.T) {
	assert.Equal(t, true, HasNil("", 1, nil))
	assert.Equal(t, false, HasNil("1", 1, struct {}{}))
	assert.Equal(t, true, HasNil())
}
