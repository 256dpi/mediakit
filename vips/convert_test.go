package vips

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/256dpi/mediakit/samples"
)

func TestConvert(t *testing.T) {
	presetConvertTest := func(t *testing.T, preset Preset, format string) {
		for _, sample := range samples.Images() {
			t.Run(sample, func(t *testing.T) {
				file := samples.Buffer(sample)
				defer file.Close()

				var buf bytes.Buffer
				err := Convert(nil, file, &buf, ConvertOptions{
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
					Pages:  1,
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

	t.Run("WebP", func(t *testing.T) {
		presetConvertTest(t, WebP, "webp")
	})
}

func TestConvertAnimations(t *testing.T) {
	for _, item := range []struct {
		name   string
		sample string
		opts   ConvertOptions
		report Report
	}{
		{
			sample: samples.AnimationGIF,
			opts: ConvertOptions{
				Preset:    WebP,
				Width:     256,
				MultiPage: true,
			},
			report: Report{
				Width:  256,
				Height: 144,
				Bands:  4,
				Color:  "srgb",
				Format: "webp",
				Pages:  10,
				Delay:  []int{200, 200, 200, 200, 200, 200, 200, 200, 200, 200},
			},
		},
		{
			sample: samples.AnimationWebP,
			opts: ConvertOptions{
				Preset:    GIFWeb,
				Width:     256,
				Height:    100,
				MultiPage: true,
			},
			report: Report{
				Width:  178,
				Height: 100,
				Bands:  4,
				Color:  "srgb",
				Format: "gif",
				Pages:  10,
				Delay:  []int{200, 200, 200, 200, 200, 200, 200, 200, 200, 200},
			},
		},
	} {
		t.Run(item.name, func(t *testing.T) {
			file := samples.Buffer(item.sample)
			defer file.Close()

			var buf bytes.Buffer
			err := Convert(nil, file, &buf, item.opts)
			assert.NoError(t, err)

			report, err := Analyze(nil, &buf)
			assert.NoError(t, err)
			assert.Equal(t, &item.report, report)
		})
	}
}

func TestConvertOptions(t *testing.T) {
	for _, item := range []struct {
		name   string
		sample string
		opts   ConvertOptions
		report Report
	}{
		{
			name:   "ResizeByWidth",
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
				Pages:  1,
			},
		},
		{
			name:   "ResizeBySize",
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
				Pages:  1,
			},
		},
		{
			name:   "CropyBySize",
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
				Pages:  1,
			},
		},
		{
			name:   "ResizeAndKeepMeta",
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
				Pages:  1,
			},
		},
	} {
		t.Run(item.name, func(t *testing.T) {
			file := samples.Load(item.sample)
			defer file.Close()

			var buf bytes.Buffer
			err := Convert(nil, file, &buf, item.opts)
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

func TestPipeline(t *testing.T) {
	file := samples.Load(samples.ImagePNG)
	defer file.Close()

	var buf bytes.Buffer
	err := Pipeline([][]string{
		{"resize", ".png", "0.2"},
		{"rotate", ".png", "90"},
	}, file, &buf)
	assert.NoError(t, err)
}
