package chromium

import (
	"bytes"
	"image"
	"image/png"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScreenshot(t *testing.T) {
	/* without allocate */

	buf, err := Screenshot(nil, "https://www.chromium.org", Options{})
	assert.NoError(t, err)
	assert.NotEmpty(t, buf)

	img, err := png.Decode(bytes.NewReader(buf))
	assert.NoError(t, err)
	assert.Equal(t, image.Rectangle{
		Min: image.Point{},
		Max: image.Point{X: 1600, Y: 1200},
	}, img.Bounds())

	/* with allocate and options */

	ctx, cancel, err := Allocate()
	assert.NoError(t, err)
	assert.NotNil(t, ctx)
	defer cancel()

	buf, err = Screenshot(ctx, "https://www.chromium.org", Options{
		Width:  1920,
		Height: 1080,
		Scale:  2,
		Full:   true,
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, buf)

	img, err = png.Decode(bytes.NewReader(buf))
	assert.NoError(t, err)
	assert.Equal(t, image.Rectangle{
		Min: image.Point{},
		Max: image.Point{X: 3840, Y: 2160},
	}, img.Bounds())

	/* longer page */

	buf, err = Screenshot(ctx, "https://en.wikipedia.org/wiki/Chromium", Options{
		Width:  1280,
		Height: 720,
		Scale:  2,
		Full:   true,
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, buf)

	img, err = png.Decode(bytes.NewReader(buf))
	assert.NoError(t, err)
	assert.Equal(t, 2560, img.Bounds().Dx())
	assert.True(t, img.Bounds().Dy() > 25_000)
}
