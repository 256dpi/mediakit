package mediakit

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProcessorConvertImage(t *testing.T) {
	input := loadSample("sample.png")

	p := NewProcessor(Config{
		Directory:    t.TempDir(),
		FS:           os.DirFS("/"),
		ImageFormats: []string{"png"},
	})

	var buf bytes.Buffer
	err := p.ConvertImage(input, &buf, KeepSize())
	assert.NoError(t, err)
	assert.Equal(t, "image/jpeg", Detect(buf.Bytes()))
}

func TestProcessorConvertAudio(t *testing.T) {
	input := loadSample("sample.wav")

	p := NewProcessor(Config{
		Directory:    t.TempDir(),
		FS:           os.DirFS("/"),
		AudioFormats: []string{"wav"},
		AudioCodecs:  []string{"pcm_s16le"},
	})

	var buf bytes.Buffer
	var progress []float64
	err := p.ConvertAudio(input, &buf, func(f float64) {
		progress = append(progress, f)
	})
	assert.NoError(t, err)
	assert.Equal(t, "audio/mpeg", Detect(buf.Bytes()))
	assert.Equal(t, []float64{0, 0.6, 1}, round(progress))
}

func TestProcessorConvertVideo(t *testing.T) {
	input := loadSample("sample.avi")

	p := NewProcessor(Config{
		Directory:    t.TempDir(),
		FS:           os.DirFS("/"),
		VideoFormats: []string{"avi"},
		VideoCodecs:  []string{"mpeg4"},
	})

	var buf bytes.Buffer
	var progress []float64
	err := p.ConvertVideo(input, &buf, MaxWidth(500), func(f float64) {
		progress = append(progress, f)
	})
	assert.NoError(t, err)
	assert.Equal(t, "video/mp4", Detect(buf.Bytes()))
	assert.Equal(t, []float64{0, 0.5, 1}, round(progress))
}

func TestProcessorExtractImage(t *testing.T) {
	input := loadSample("sample.avi")

	p := NewProcessor(Config{
		Directory:    t.TempDir(),
		FS:           os.DirFS("/"),
		ImageFormats: []string{"png"},
		VideoFormats: []string{"avi"},
		VideoCodecs:  []string{"mpeg4"},
	})

	var buf bytes.Buffer
	err := p.ExtractImage(input, &buf, 0.25, KeepSize())
	assert.NoError(t, err)
	assert.Equal(t, "image/jpeg", Detect(buf.Bytes()))
}
