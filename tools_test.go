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

	rep, err := Analyze(nil, output)
	assert.NoError(t, err)
	assert.Equal(t, &Report{
		MediaType:  "image/jpeg",
		FileFormat: "jpeg",
		Width:      800,
		Height:     533,
	}, rep)
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

	rep, err := Analyze(nil, output)
	assert.NoError(t, err)
	assert.Equal(t, &Report{
		MediaType:  "audio/mpeg",
		FileFormat: "mp3",
		Streams:    []string{"audio"},
		Codecs:     []string{"mp3"},
		Duration:   2.089796,
		Channels:   2,
		SampleRate: 44100,
	}, rep)
}

func TestExtractAudio(t *testing.T) {
	input := samples.Buffer(samples.VideoMOV)
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

	rep, err := Analyze(nil, output)
	assert.NoError(t, err)
	assert.Equal(t, &Report{
		MediaType:  "audio/mpeg",
		FileFormat: "mp3",
		Streams:    []string{"audio"},
		Codecs:     []string{"mp3"},
		Duration:   2.089796,
		Channels:   2,
		SampleRate: 44100,
	}, rep)
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

	rep, err := Analyze(nil, output)
	assert.NoError(t, err)
	assert.Equal(t, &Report{
		MediaType:  "video/mp4",
		FileFormat: "mov,mp4,m4a,3gp,3g2,mj2",
		Width:      500,
		Height:     282,
		Streams:    []string{"video", "audio"},
		Codecs:     []string{"h264", "aac"},
		Duration:   2.12,
		Channels:   2,
		SampleRate: 44100,
		FrameRate:  25,
	}, rep)
}

func TestExtractAnimation(t *testing.T) {
	input := samples.Buffer(samples.VideoMPEG)
	output := makeBuffers(t.TempDir(), "output")[0]

	err := ConvertVideo(nil, input, output, ffmpeg.AnimationGIF, MaxWidth(500), 30, 48000, nil)
	assert.NoError(t, err)

	rep, err := Analyze(nil, output)
	assert.NoError(t, err)
	assert.Equal(t, &Report{
		MediaType:  "image/gif",
		FileFormat: "gif",
		Width:      500,
		Height:     281,
		Duration:   2.03,
		FrameRate:  30.04926108374384,
	}, rep)
}

func TestConvertAnimation(t *testing.T) {
	input := samples.Buffer(samples.AnimationGIF)
	output := makeBuffers(t.TempDir(), "output")[0]

	err := ConvertImage(nil, input, output, vips.WebP, MaxWidth(500))
	assert.NoError(t, err)

	rep, err := Analyze(nil, output)
	assert.NoError(t, err)
	assert.Equal(t, &Report{
		MediaType:  "image/webp",
		FileFormat: "webp",
		Width:      500,
		Height:     281,
		Duration:   2,
		FrameRate:  5,
	}, rep)
}

func TestExtractImage(t *testing.T) {
	input := samples.Buffer(samples.VideoMOV)
	buffers := makeBuffers(t.TempDir(), "temp", "output")

	err := ExtractImage(nil, input, buffers[0], buffers[1], 0.25, vips.JPGWeb, KeepSize())
	assert.NoError(t, err)

	rep, err := Analyze(nil, buffers[1])
	assert.NoError(t, err)
	assert.Equal(t, &Report{
		MediaType:  "image/jpeg",
		FileFormat: "jpeg",
		Width:      800,
		Height:     450,
	}, rep)
}

func TestCaptureScreenshot(t *testing.T) {
	output := makeBuffers(t.TempDir(), "output")[0]

	err := CaptureScreenshot(nil, "https://example.org", output, chromium.ScreenshotOptions{})
	assert.NoError(t, err)

	buf := make([]byte, DetectBytes)
	_, err = io.ReadFull(output, buf)
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
