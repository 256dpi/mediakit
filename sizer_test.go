package mediakit

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSizer(t *testing.T) {
	// keep size
	assert.Equal(t, Size{W: 42, H: 84}, KeepSize()(Size{W: 42, H: 84}))

	// max width
	assert.Equal(t, Size{W: 42, H: 84}, MaxWidth(50)(Size{W: 42, H: 84}))
	assert.Equal(t, Size{W: 42, H: 84}, MaxWidth(42)(Size{W: 42, H: 84}))
	assert.Equal(t, Size{W: 20, H: 40}, MaxWidth(20)(Size{W: 42, H: 84}))
	assert.Equal(t, Size{W: 84, H: 42}, MaxWidth(100)(Size{W: 84, H: 42}))
	assert.Equal(t, Size{W: 84, H: 42}, MaxWidth(84)(Size{W: 84, H: 42}))
	assert.Equal(t, Size{W: 50, H: 25}, MaxWidth(50)(Size{W: 84, H: 42}))

	// max height
	assert.Equal(t, Size{W: 42, H: 84}, MaxHeight(100)(Size{W: 42, H: 84}))
	assert.Equal(t, Size{W: 42, H: 84}, MaxHeight(84)(Size{W: 42, H: 84}))
	assert.Equal(t, Size{W: 25, H: 50}, MaxHeight(50)(Size{W: 42, H: 84}))
	assert.Equal(t, Size{W: 84, H: 42}, MaxHeight(50)(Size{W: 84, H: 42}))
	assert.Equal(t, Size{W: 84, H: 42}, MaxHeight(42)(Size{W: 84, H: 42}))
	assert.Equal(t, Size{W: 40, H: 20}, MaxHeight(20)(Size{W: 84, H: 42}))

	// max area
	assert.Equal(t, Size{W: 42, H: 84}, MaxArea(50*100)(Size{W: 42, H: 84}))
	assert.Equal(t, Size{W: 42, H: 84}, MaxArea(42*84)(Size{W: 42, H: 84}))
	assert.Equal(t, Size{W: 25, H: 50}, MaxArea(25*50)(Size{W: 42, H: 84}))
	assert.Equal(t, Size{W: 84, H: 42}, MaxArea(50*100)(Size{W: 84, H: 42}))
	assert.Equal(t, Size{W: 84, H: 42}, MaxArea(42*84)(Size{W: 84, H: 42}))
	assert.Equal(t, Size{W: 50, H: 25}, MaxArea(25*50)(Size{W: 84, H: 42}))

	// max size
	assert.Equal(t, Size{W: 42, H: 84}, MaxSize(Size{W: 50, H: 100})(Size{W: 42, H: 84}))
	assert.Equal(t, Size{W: 42, H: 84}, MaxSize(Size{W: 42, H: 84})(Size{W: 42, H: 84}))
	assert.Equal(t, Size{W: 20, H: 40}, MaxSize(Size{W: 20, H: 84})(Size{W: 42, H: 84}))
	assert.Equal(t, Size{W: 25, H: 50}, MaxSize(Size{W: 42, H: 50})(Size{W: 42, H: 84}))
}
