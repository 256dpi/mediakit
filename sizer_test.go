package mediakit

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSizer(t *testing.T) {
	// keep size
	assert.Equal(t, Size{Width: 42, Height: 84}, KeepSize()(Size{Width: 42, Height: 84}))

	// max width
	assert.Equal(t, Size{Width: 42, Height: 84}, MaxWidth(50)(Size{Width: 42, Height: 84}))
	assert.Equal(t, Size{Width: 42, Height: 84}, MaxWidth(42)(Size{Width: 42, Height: 84}))
	assert.Equal(t, Size{Width: 20, Height: 40}, MaxWidth(20)(Size{Width: 42, Height: 84}))
	assert.Equal(t, Size{Width: 84, Height: 42}, MaxWidth(100)(Size{Width: 84, Height: 42}))
	assert.Equal(t, Size{Width: 84, Height: 42}, MaxWidth(84)(Size{Width: 84, Height: 42}))
	assert.Equal(t, Size{Width: 50, Height: 25}, MaxWidth(50)(Size{Width: 84, Height: 42}))

	// max height
	assert.Equal(t, Size{Width: 42, Height: 84}, MaxHeight(100)(Size{Width: 42, Height: 84}))
	assert.Equal(t, Size{Width: 42, Height: 84}, MaxHeight(84)(Size{Width: 42, Height: 84}))
	assert.Equal(t, Size{Width: 25, Height: 50}, MaxHeight(50)(Size{Width: 42, Height: 84}))
	assert.Equal(t, Size{Width: 84, Height: 42}, MaxHeight(50)(Size{Width: 84, Height: 42}))
	assert.Equal(t, Size{Width: 84, Height: 42}, MaxHeight(42)(Size{Width: 84, Height: 42}))
	assert.Equal(t, Size{Width: 40, Height: 20}, MaxHeight(20)(Size{Width: 84, Height: 42}))

	// max area
	assert.Equal(t, Size{Width: 42, Height: 84}, MaxArea(50*100)(Size{Width: 42, Height: 84}))
	assert.Equal(t, Size{Width: 42, Height: 84}, MaxArea(42*84)(Size{Width: 42, Height: 84}))
	assert.Equal(t, Size{Width: 25, Height: 50}, MaxArea(25*50)(Size{Width: 42, Height: 84}))
	assert.Equal(t, Size{Width: 84, Height: 42}, MaxArea(50*100)(Size{Width: 84, Height: 42}))
	assert.Equal(t, Size{Width: 84, Height: 42}, MaxArea(42*84)(Size{Width: 84, Height: 42}))
	assert.Equal(t, Size{Width: 50, Height: 25}, MaxArea(25*50)(Size{Width: 84, Height: 42}))

	// max size
	assert.Equal(t, Size{Width: 42, Height: 84}, MaxSize(Size{Width: 50, Height: 100})(Size{Width: 42, Height: 84}))
	assert.Equal(t, Size{Width: 42, Height: 84}, MaxSize(Size{Width: 42, Height: 84})(Size{Width: 42, Height: 84}))
	assert.Equal(t, Size{Width: 20, Height: 40}, MaxSize(Size{Width: 20, Height: 84})(Size{Width: 42, Height: 84}))
	assert.Equal(t, Size{Width: 25, Height: 50}, MaxSize(Size{Width: 42, Height: 50})(Size{Width: 42, Height: 84}))
}
