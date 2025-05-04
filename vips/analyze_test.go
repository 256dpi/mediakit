package vips

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/256dpi/mediakit/samples"
)

func TestAnalyzeImages(t *testing.T) {
	reports := map[string]*Report{
		samples.ImageGIF: {
			Width:  800,
			Height: 533,
			Bands:  4,
			Color:  "srgb",
			Format: "gif",
			Pages:  1,
		},
		samples.ImageHEIF: {
			Width:  800,
			Height: 533,
			Bands:  3,
			Color:  "srgb",
			Format: "heif",
			Pages:  1,
		},
		samples.ImageJPEG: {
			Width:  800,
			Height: 533,
			Bands:  3,
			Color:  "srgb",
			Format: "jpeg",
			Pages:  1,
		},
		samples.ImageJPEG2K: {
			Width:  800,
			Height: 533,
			Bands:  3,
			Color:  "srgb",
			Format: "jp2k",
			Pages:  6, // TODO: Why?
		},
		samples.ImagePDF: {
			Width:  800,
			Height: 533,
			Bands:  4,
			Color:  "srgb",
			Format: "pdf",
			Pages:  1,
		},
		samples.ImagePNG: {
			Width:  800,
			Height: 533,
			Bands:  3,
			Color:  "srgb",
			Format: "png",
			Pages:  1,
		},
		samples.ImageTIFF: {
			Width:  800,
			Height: 533,
			Bands:  4,
			Color:  "srgb",
			Format: "tiff",
			Pages:  1,
		},
		samples.ImageWebP: {
			Width:  800,
			Height: 533,
			Bands:  3,
			Color:  "srgb",
			Format: "webp",
			Pages:  1,
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

func TestAnalyzePDF(t *testing.T) {
	file := samples.Load(samples.DocumentPDF)
	defer file.Close()

	report, err := Analyze(nil, file)
	assert.NoError(t, err)
	assert.Equal(t, &Report{
		Width:  595,
		Height: 842,
		Bands:  4,
		Color:  "srgb",
		Format: "pdf",
		Pages:  20,
	}, report)
}

func TestAnalyzeError(t *testing.T) {
	report, err := Analyze(nil, strings.NewReader("foo"))
	assert.Error(t, err)
	assert.Nil(t, report)
	assert.True(t, strings.Contains(err.Error(), "unable to load source") || strings.Contains(err.Error(), "exit status 255"))
}
