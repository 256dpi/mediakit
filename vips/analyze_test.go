package vips

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/256dpi/mediakit/samples"
)

func TestAnalyze(t *testing.T) {
	for i, item := range []struct {
		sample string
		report Report
	}{
		{
			sample: samples.ImageGIF,
			report: Report{
				Width:  800,
				Height: 533,
				Bands:  4,
				Color:  "srgb",
				Format: "gif",
			},
		},
		{
			sample: samples.ImageHEIF,
			report: Report{
				Width:  800,
				Height: 533,
				Bands:  3,
				Color:  "srgb",
				Format: "heif",
			},
		},
		{
			sample: samples.ImageJPEG,
			report: Report{
				Width:  800,
				Height: 533,
				Bands:  3,
				Color:  "srgb",
				Format: "jpeg",
			},
		},
		{
			sample: samples.ImageJPEG2K,
			report: Report{
				Width:  800,
				Height: 533,
				Bands:  3,
				Color:  "srgb",
				Format: "jp2k",
			},
		},
		{
			sample: samples.ImagePDF,
			report: Report{
				Width:  0,
				Height: 533,
				Bands:  4,
				Color:  "srgb",
				Format: "pdf",
			},
		},
		{
			sample: samples.ImagePNG,
			report: Report{
				Width:  800,
				Height: 533,
				Bands:  3,
				Color:  "srgb",
				Format: "png",
			},
		},
		{
			sample: samples.ImageTIFF,
			report: Report{
				Width:  800,
				Height: 533,
				Bands:  4,
				Color:  "srgb",
				Format: "tiff",
			},
		},
		{
			sample: samples.ImageWebP,
			report: Report{
				Width:  800,
				Height: 533,
				Bands:  3,
				Color:  "srgb",
				Format: "webp",
			},
		},
	} {
		t.Run(strconv.Itoa(i)+"-"+item.sample, func(t *testing.T) {
			sample := samples.Load(item.sample)
			defer sample.Close()

			report, err := Analyze(sample)
			assert.NoError(t, err)
			assert.Equal(t, &item.report, report)
		})
	}
}
