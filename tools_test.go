package mediakit

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/256dpi/mediakit/ffmpeg"
	"github.com/256dpi/mediakit/samples"
	"github.com/256dpi/mediakit/vips"
)

func TestConvertImage(t *testing.T) {
	input := samples.Buffer(samples.ImagePNG)
	output := makeBuffers(t.TempDir(), "output")[0]

	err := ConvertImage(input, output, vips.JPGWeb, KeepSize())
	assert.NoError(t, err)

	buf := make([]byte, DetectBytes)
	_, err = io.ReadFull(output, buf)
	assert.NoError(t, err)
	assert.Equal(t, "image/jpeg", Detect(buf))
}

func TestConvertAudio(t *testing.T) {
	input := samples.Buffer(samples.AudioWAV)
	output := makeBuffers(t.TempDir(), "output")[0]

	var progress []float64
	err := ConvertAudio(input, output, ffmpeg.AudioMP3VBRStandard, 48000, func(f float64) {
		progress = append(progress, f)
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
	output := makeBuffers(t.TempDir(), "output")[0]

	var progress []float64
	err := ConvertVideo(input, output, ffmpeg.VideoMP4H264AACFast, MaxWidth(500), 30, func(f float64) {
		progress = append(progress, f)
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
	buffers := makeBuffers(t.TempDir(), "temp", "output")

	err := ExtractImage(input, buffers[0], buffers[1], 0.25, vips.JPGWeb, KeepSize())
	assert.NoError(t, err)

	buf := make([]byte, DetectBytes)
	_, err = io.ReadFull(buffers[1], buf)
	assert.NoError(t, err)
	assert.Equal(t, "image/jpeg", Detect(buf))
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