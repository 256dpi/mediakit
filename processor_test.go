package mediakit

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/256dpi/mediakit/samples"
)

func TestProcessorConvertImage(t *testing.T) {
	input := samples.Load(samples.ImagePNG)

	p := NewProcessor(Config{
		Directory: t.TempDir(),
	})

	var buf bytes.Buffer
	err := p.ConvertImage(input, KeepSize(), func(output *os.File) error {
		_, err := io.Copy(&buf, output)
		return err
	})
	assert.NoError(t, err)
	assert.Equal(t, "image/jpeg", Detect(buf.Bytes()))
}

func TestProcessorConvertAudio(t *testing.T) {
	input := samples.Load(samples.AudioWAV)

	p := NewProcessor(Config{
		Directory: t.TempDir(),
	})

	var buf bytes.Buffer
	var progress []float64
	err := p.ConvertAudio(input, func(f float64) {
		progress = append(progress, f)
	}, func(output *os.File) error {
		_, err := io.Copy(&buf, output)
		return err
	})
	assert.NoError(t, err)
	assert.Equal(t, "audio/mpeg", Detect(buf.Bytes()))
	assert.True(t, len(progress) >= 2)
}

func TestProcessorConvertVideo(t *testing.T) {
	input := samples.Load(samples.VideoAVI)

	p := NewProcessor(Config{
		Directory: t.TempDir(),
	})

	var buf bytes.Buffer
	var progress []float64
	err := p.ConvertVideo(input, MaxWidth(500), func(f float64) {
		progress = append(progress, f)
	}, func(output *os.File) error {
		_, err := io.Copy(&buf, output)
		return err
	})
	assert.NoError(t, err)
	assert.Equal(t, "video/mp4", Detect(buf.Bytes()))
	assert.True(t, len(progress) >= 2)
}

func TestProcessorExtractImage(t *testing.T) {
	input := samples.Load(samples.VideoMOV)

	p := NewProcessor(Config{
		Directory: t.TempDir(),
	})

	var buf bytes.Buffer
	err := p.ExtractImage(input, 0.25, KeepSize(), func(output *os.File) error {
		_, err := io.Copy(&buf, output)
		return err
	})
	assert.NoError(t, err)
	assert.Equal(t, "image/jpeg", Detect(buf.Bytes()))
}
