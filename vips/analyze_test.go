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
			sample: "sample.png",
			report: Report{
				Width:  1280,
				Height: 853,
				Bands:  3,
				Color:  "srgb",
				Format: "png",
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
