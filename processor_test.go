package mediakit

import (
	"bytes"
	"io"
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
	err := p.ConvertImage(input, KeepSize(), func(output io.ReadSeeker) error {
		_, err := io.Copy(&buf, output)
		return err
	})
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
	err := p.ConvertAudio(input, func(f float64) {
		progress = append(progress, f)
	}, func(output io.ReadSeeker) error {
		_, err := io.Copy(&buf, output)
		return err
	})
	assert.NoError(t, err)
	assert.Equal(t, "audio/mpeg", Detect(buf.Bytes()))
	assert.True(t, len(progress) >= 2)
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
	err := p.ConvertVideo(input, MaxWidth(500), func(f float64) {
		progress = append(progress, f)
	}, func(output io.ReadSeeker) error {
		_, err := io.Copy(&buf, output)
		return err
	})
	assert.NoError(t, err)
	assert.Equal(t, "video/mp4", Detect(buf.Bytes()))
	assert.True(t, len(progress) >= 2)
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
	err := p.ExtractImage(input, 0.25, KeepSize(), func(output io.ReadSeeker) error {
		_, err := io.Copy(&buf, output)
		return err
	})
	assert.NoError(t, err)
	assert.Equal(t, "image/jpeg", Detect(buf.Bytes()))
}
