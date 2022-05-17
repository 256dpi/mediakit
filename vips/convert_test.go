package vips

import (
	"bytes"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvert(t *testing.T) {
	for i, item := range []struct {
		sample string
		opts   ConvertOptions
		report Report
	}{
		{
			sample: "sample.png",
			opts: ConvertOptions{
				Preset: JPGWeb,
				Width:  256,
			},
			report: Report{
				Width:  256,
				Height: 171,
				Bands:  3,
				Color:  "srgb",
				Format: "jpeg",
			},
		},
		{
			sample: "sample.jpg",
			opts: ConvertOptions{
				Preset: PNGWeb,
				Width:  512,
				Height: 256,
			},
			report: Report{
				Width:  384,
				Height: 256,
				Bands:  3,
				Color:  "srgb",
				Format: "png",
			},
		},
		{
			sample: "sample.gif",
			opts: ConvertOptions{
				Preset: JPGWeb,
				Width:  256,
				Height: 256,
				Crop:   true,
			},
			report: Report{
				Width:  256,
				Height: 256,
				Bands:  3,
				Color:  "srgb",
				Format: "jpeg",
			},
		},
		{
			sample: "sample.png",
			opts: ConvertOptions{
				Preset:      JPGWeb,
				Width:       256,
				KeepProfile: true,
				NoRotate:    true,
			},
			report: Report{
				Width:  256,
				Height: 171,
				Bands:  3,
				Color:  "srgb",
				Format: "jpeg",
			},
		},
	} {
		t.Run(strconv.Itoa(i)+"-"+item.sample, func(t *testing.T) {
			sample := loadSample(item.sample)
			defer sample.Close()

			var buf bytes.Buffer
			err := Convert(sample, &buf, item.opts)
			assert.NoError(t, err)

			report, err := Analyze(&buf)
			assert.NoError(t, err)
			assert.Equal(t, &item.report, report)
		})
	}
}

func TestConvertError(t *testing.T) {
	var buf bytes.Buffer
	err := Convert(strings.NewReader("foo"), &buf, ConvertOptions{
		Width: 1,
	})
	assert.Error(t, err)
	assert.Equal(t, "vipsforeignload: source is not in a known format", err.Error())
}
