package mediakit

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFFMPEG(t *testing.T) {
	sample := loadSample("hevc")
	defer sample.Close()

	var out bytes.Buffer
	err := FFMPEG(sample, &out, FFMPEGOptions{
		Format:   "webm",
		Duration: 1,
	})
	assert.NoError(t, err)

	report, err := FFProbe(&out)
	assert.NoError(t, err)
	assert.Equal(t, &Report{
		Format: Format{
			Name:       "matroska,webm",
			LongName:   "Matroska / WebM",
			ProbeScore: 100,
			Duration:   1,
		}, Streams: []Stream{
			{
				CodecName:     "vp9",
				CodecLongName: "Google VP9",
				CodecType:     "video",
				BitRate:       0,
				Duration:      0,
				Width:         1280,
				Height:        720,
			},
		},
	}, report)
}

func TestFFMpegError(t *testing.T) {
	err := FFMPEG(strings.NewReader("foo"), io.Discard, FFMPEGOptions{})
	assert.Error(t, err)
	assert.Equal(t, "pipe:: invalid data found when processing input", err.Error())
}
