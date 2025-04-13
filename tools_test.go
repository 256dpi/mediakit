package mediakit

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/256dpi/mediakit/chromium"
	"github.com/256dpi/mediakit/ffmpeg"
	"github.com/256dpi/mediakit/samples"
	"github.com/256dpi/mediakit/vips"
)

func init() {
	ffmpeg.WarningsLogger = log.Default()
}

func TestConvertImage(t *testing.T) {
	input := samples.Buffer(samples.ImagePNG)
	output := makeBuffers(t.TempDir(), "output")[0]

	err := ConvertImage(nil, input, output, vips.JPGWeb, KeepSize())
	assert.NoError(t, err)

	buf := make([]byte, DetectBytes)
	_, err = io.ReadFull(output, buf)
	assert.NoError(t, err)
	assert.Equal(t, "image/jpeg", Detect(buf, false))
}

func TestConvertAudio(t *testing.T) {
	input := samples.Buffer(samples.AudioWAV)
	output := makeBuffers(t.TempDir(), "output")[0]

	var progress []float64
	err := ConvertAudio(nil, input, output, ffmpeg.AudioMP3VBRStandard, 48000, &Progress{
		Rate: time.Second,
		Func: func(f float64) {
			progress = append(progress, f)
		},
	})
	assert.NoError(t, err)
	assert.True(t, len(progress) >= 1)

	buf := make([]byte, DetectBytes)
	_, err = io.ReadFull(output, buf)
	assert.NoError(t, err)
	assert.Equal(t, "audio/mpeg", Detect(buf, false))
}

func TestConvertVideo(t *testing.T) {
	input := samples.Buffer(samples.VideoAVI)
	output := makeBuffers(t.TempDir(), "output")[0]

	var progress []float64
	err := ConvertVideo(nil, input, output, ffmpeg.VideoMP4H264AACFast, MaxWidth(500), 30, 48000, &Progress{
		Rate: time.Second,
		Func: func(f float64) {
			progress = append(progress, f)
		},
	})
	assert.NoError(t, err)
	assert.True(t, len(progress) >= 1)

	buf := make([]byte, DetectBytes)
	_, err = io.ReadFull(output, buf)
	assert.NoError(t, err)
	assert.Equal(t, "video/mp4", Detect(buf, false))
}

func TestExtractImage(t *testing.T) {
	input := samples.Buffer(samples.VideoMOV)
	buffers := makeBuffers(t.TempDir(), "temp", "output")

	err := ExtractImage(nil, input, buffers[0], buffers[1], 0.25, vips.JPGWeb, KeepSize())
	assert.NoError(t, err)

	buf := make([]byte, DetectBytes)
	_, err = io.ReadFull(buffers[1], buf)
	assert.NoError(t, err)
	assert.Equal(t, "image/jpeg", Detect(buf, false))
}

func TestCaptureScreenshot(t *testing.T) {
	output := makeBuffers(t.TempDir(), "output")[0]

	err := CaptureScreenshot(nil, "https://www.chromium.org", output, chromium.ScreenshotOptions{})
	assert.NoError(t, err)

	buf := make([]byte, DetectBytes)
	_, err = io.ReadFull(output, buf)
	assert.NoError(t, err)
	assert.Equal(t, "image/png", Detect(buf, false))
}

func makeBuffers(dir string, names ...string) []*os.File {
	var list []*os.File
	for _, name := range names {
		file, err := os.Create(filepath.Join(dir, name))
		if err != nil {
			panic(err)
		}
		list = append(list, file)
	}
	return list
}
