package mediakit

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/256dpi/mediakit/ffmpeg"
	"github.com/256dpi/mediakit/samples"
	"github.com/256dpi/mediakit/vips"
)

func init() {
	ffmpeg.WarningsLogger = log.Default()
}

func TestConvertImage(t *testing.T) {
	input := samples.Buffer(samples.ImagePNG)
	output := makeBuffer(t.TempDir(), "output")

	err := ConvertImage(nil, input, output, vips.JPGWeb, KeepSize())
	assert.NoError(t, err)

	buf := make([]byte, DetectBytes)
	_, err = io.ReadFull(output, buf)
	assert.NoError(t, err)
	assert.Equal(t, "image/jpeg", Detect(buf))
}

func TestConvertAudio(t *testing.T) {
	input := samples.Buffer(samples.AudioWAV)
	output := makeBuffer(t.TempDir(), "output")

	var progress []float64
	err := ConvertAudio(nil, input, output, ffmpeg.AudioMP3VBRStandard, 48000, &Progress{
		Rate: time.Second,
		Func: func(f float64) {
			progress = append(progress, f)
		},
	})
	assert.NoError(t, err)
	assert.True(t, len(progress) >= 2)

	buf := make([]byte, DetectBytes)
	_, err = io.ReadFull(output, buf)
	assert.NoError(t, err)
	assert.Equal(t, "audio/mpeg", Detect(buf))
}

func TestConvertVideo(t *testing.T) {
	input := samples.Buffer(samples.VideoAVI)
	output := makeBuffer(t.TempDir(), "output")

	var progress []float64
	err := ConvertVideo(nil, input, output, ffmpeg.VideoMP4H264AACFast, MaxWidth(500), 30, 48000, &Progress{
		Rate: time.Second,
		Func: func(f float64) {
			progress = append(progress, f)
		},
	})
	assert.NoError(t, err)
	assert.True(t, len(progress) >= 2)

	buf := make([]byte, DetectBytes)
	_, err = io.ReadFull(output, buf)
	assert.NoError(t, err)
	assert.Equal(t, "video/mp4", Detect(buf))
}

func TestExtractImage(t *testing.T) {
	input := samples.Buffer(samples.VideoMOV)
	output := makeBuffer(t.TempDir(), "output")

	err := ExtractImage(nil, input, output, 0.25, vips.JPGWeb, KeepSize())
	assert.NoError(t, err)

	buf := make([]byte, DetectBytes)
	_, err = io.ReadFull(output, buf)
	assert.NoError(t, err)
	assert.Equal(t, "image/jpeg", Detect(buf))
}

func makeBuffer(dir string, name string) *os.File {
	file, err := os.Create(filepath.Join(dir, name))
	if err != nil {
		panic(err)
	}

	return file
}
