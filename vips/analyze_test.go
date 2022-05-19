package vips

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAnalyze(t *testing.T) {
	for i, item := range []struct {
		sample string
		report Report
	}{
		{
			sample: "sample.gif",
			report: Report{
				Width:  1280,
				Height: 853,
				Bands:  3,
				Color:  "srgb",
				Format: "gif",
			},
		},
		{
			sample: "sample.heif",
			report: Report{
				Width:  8736,
				Height: 5856,
				Bands:  3,
				Color:  "srgb",
				Format: "heif",
			},
		},
		{
			sample: "sample.jpg",
			report: Report{
				Width:  1280,
				Height: 853,
				Bands:  3,
				Color:  "srgb",
				Format: "jpeg",
			},
		},
		{
			sample: "sample.pdf",
			report: Report{
				Width:  0,
				Height: 792,
				Bands:  4,
				Color:  "srgb",
				Format: "pdf",
			},
		},
		{
			sample: "sample.png",
			report: Report{
				Width:  1280,
				Height: 853,
				Bands:  3,
				Color:  "srgb",
				Format: "png",
			},
		},
		{
			sample: "sample.tiff",
			report: Report{
				Width:  1280,
				Height: 853,
				Bands:  3,
				Color:  "srgb",
				Format: "tiff",
			},
		},
		{
			sample: "sample.webp",
			report: Report{
				Width:  4275,
				Height: 2451,
				Bands:  3,
				Color:  "srgb",
				Format: "webp",
			},
		},
	} {
		t.Run(strconv.Itoa(i)+"-"+item.sample, func(t *testing.T) {
			sample := loadSample(item.sample)
			defer sample.Close()

			report, err := Analyze(sample)
			assert.NoError(t, err)
			assert.Equal(t, &item.report, report)
		})
	}
}
