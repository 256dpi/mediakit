package vips

import (
	"bytes"
	"image"
	"image/jpeg"
	"image/png"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestThumbnail(t *testing.T) {
	for i, item := range []struct {
		sample string
		opts   ThumbnailOptions
		size   image.Point
	}{
		{
			sample: "sample.png",
			opts: ThumbnailOptions{
				Preset: JPGWeb,
				Width:  128,
			},
			size: image.Pt(128, 86),
		},
		{
			sample: "sample.jpg",
			opts: ThumbnailOptions{
				Preset: PNGWeb,
				Width:  128,
			},
			size: image.Pt(128, 85),
		},
		{
			sample: "sample.gif",
			opts: ThumbnailOptions{
				Preset: JPGWeb,
				Width:  128,
			},
			size: image.Pt(128, 86),
		},
	} {
		t.Run(strconv.Itoa(i)+"-"+item.sample, func(t *testing.T) {
			sample := loadSample(item.sample)

			var buf bytes.Buffer
			err := Thumbnail(sample, &buf, item.opts)
			assert.NoError(t, err)

			if item.opts.Preset == JPGWeb {
				img, err := jpeg.Decode(&buf)
				assert.NoError(t, err)
				assert.Equal(t, item.size, img.Bounds().Size())
			} else if item.opts.Preset == PNGWeb {
				img, err := png.Decode(&buf)
				assert.NoError(t, err)
				assert.Equal(t, item.size, img.Bounds().Size())
			}
		})
	}
}

func TestThumbnailError(t *testing.T) {
	var buf bytes.Buffer
	err := Thumbnail(strings.NewReader("foo"), &buf, ThumbnailOptions{})
	assert.Error(t, err)
	assert.Equal(t, "vipsforeignload: source is not in a known format", err.Error())
}
