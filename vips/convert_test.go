package vips

import (
	"bytes"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/256dpi/mediakit/samples"
)

func TestConvert(t *testing.T) {
	presetConvertTest := func(t *testing.T, preset Preset, format string) {
		for _, name := range samples.Images() {
			t.Run(name, func(t *testing.T) {
				sample := samples.Load(name)
				defer sample.Close()

				var buf bytes.Buffer
				err := Convert(nil, sample, &buf, ConvertOptions{
					Preset: preset,
					Width:  256,
					Height: 256,
					Crop:   true,
				})
				assert.NoError(t, err)

				report, err := Analyze(nil, &buf)
				assert.NoError(t, err)
				assert.Equal(t, &Report{
					Width:  256,
					Height: 256,
					Bands:  report.Bands, // may be 3 or 4
					Color:  "srgb",
					Format: format,
				}, report)
			})
		}
	}

	t.Run("JPG", func(t *testing.T) {
		presetConvertTest(t, JPGWeb, "jpeg")
	})

	t.Run("PNG", func(t *testing.T) {
		presetConvertTest(t, PNGWeb, "png")
	})
}

func TestConvertOptions(t *testing.T) {
	for i, item := range []struct {
		sample string
		opts   ConvertOptions
		report Report
	}{
		{
			sample: samples.ImageJPEG,
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
			sample: samples.ImageJPEG,
			opts: ConvertOptions{
				Preset: JPGWeb,
				Width:  512,
				Height: 256,
			},
			report: Report{
				Width:  384,
				Height: 256,
				Bands:  3,
				Color:  "srgb",
				Format: "jpeg",
			},
		},
		{
			sample: samples.ImageJPEG,
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
			sample: samples.ImageJPEG,
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
			sample := samples.Load(item.sample)
			defer sample.Close()

			var buf bytes.Buffer
			err := Convert(nil, sample, &buf, item.opts)
			assert.NoError(t, err)

			report, err := Analyze(nil, &buf)
			assert.NoError(t, err)
			assert.Equal(t, &item.report, report)
		})
	}
}

func TestConvertError(t *testing.T) {
	var buf bytes.Buffer
	err := Convert(nil, strings.NewReader("foo"), &buf, ConvertOptions{
		Preset: JPGWeb,
		Width:  1,
	})
	assert.Error(t, err)
	assert.Equal(t, "vipsforeignload: source is not in a known format", err.Error())
}
