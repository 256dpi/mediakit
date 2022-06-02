package mediakit

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/256dpi/mediakit/ffmpeg"
	"github.com/256dpi/mediakit/samples"
	"github.com/256dpi/mediakit/vips"
)

func TestProcessorConvertImage(t *testing.T) {
	input := samples.Load(samples.ImagePNG)

	p := NewProcessor(t.TempDir())

	var buf bytes.Buffer
	err := p.ConvertImage(input, vips.JPGWeb, KeepSize(), func(output *os.File) error {
		_, err := io.Copy(&buf, output)
		return err
	})
	assert.NoError(t, err)
	assert.Equal(t, "image/jpeg", Detect(buf.Bytes()))
}

func TestProcessorConvertAudio(t *testing.T) {
	input := samples.Load(samples.AudioWAV)

	p := NewProcessor(t.TempDir())

	var buf bytes.Buffer
	var progress []float64
	err := p.ConvertAudio(input, ffmpeg.AudioMP3VBRStandard, func(f float64) {
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

	p := NewProcessor(t.TempDir())

	var buf bytes.Buffer
	var progress []float64
	err := p.ConvertVideo(input, ffmpeg.VideoMP4H264AACFast, MaxWidth(500), 30, func(f float64) {
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

	p := NewProcessor(t.TempDir())

	var buf bytes.Buffer
	err := p.ExtractImage(input, 0.25, vips.JPGWeb, KeepSize(), func(output *os.File) error {
		_, err := io.Copy(&buf, output)
		return err
	})
	assert.NoError(t, err)
	assert.Equal(t, "image/jpeg", Detect(buf.Bytes()))
}
