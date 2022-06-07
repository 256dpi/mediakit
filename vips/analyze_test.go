package vips

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/256dpi/mediakit/samples"
)

func TestAnalyze(t *testing.T) {
	reports := map[string]*Report{
		samples.ImageGIF: {
			Width:  800,
			Height: 533,
			Bands:  4,
			Color:  "srgb",
			Format: "gif",
		},
		samples.ImageHEIF: {
			Width:  800,
			Height: 533,
			Bands:  3,
			Color:  "srgb",
			Format: "heif",
		},
		samples.ImageJPEG: {
			Width:  800,
			Height: 533,
			Bands:  3,
			Color:  "srgb",
			Format: "jpeg",
		},
		samples.ImageJPEG2K: {
			Width:  800,
			Height: 533,
			Bands:  3,
			Color:  "srgb",
			Format: "jp2k",
		},
		samples.ImagePDF: {
			Width:  800,
			Height: 533,
			Bands:  4,
			Color:  "srgb",
			Format: "pdf",
		},
		samples.ImagePNG: {
			Width:  800,
			Height: 533,
			Bands:  3,
			Color:  "srgb",
			Format: "png",
		},
		samples.ImageTIFF: {
			Width:  800,
			Height: 533,
			Bands:  4,
			Color:  "srgb",
			Format: "tiff",
		},
		samples.ImageWebP: {
			Width:  800,
			Height: 533,
			Bands:  3,
			Color:  "srgb",
			Format: "webp",
		},
	}

	for _, sample := range samples.Images() {
		t.Run(sample, func(t *testing.T) {
			file := samples.Load(sample)
			defer file.Close()

			report, err := Analyze(nil, file)
			assert.NoError(t, err)
			assert.Equal(t, reports[sample], report)
		})
	}
}
