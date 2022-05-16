package ffmpeg

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTranscode(t *testing.T) {
	sample := loadSample("hevc")
	defer sample.Close()

	var out bytes.Buffer
	err := Transcode(sample, &out, TranscodeOptions{
		Format:   "webm",
		Duration: 1,
	})
	assert.NoError(t, err)

	report, err := Analyze(&out)
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

func TestTranscodeError(t *testing.T) {
	err := Transcode(strings.NewReader("foo"), io.Discard, TranscodeOptions{})
	assert.Error(t, err)
	assert.Equal(t, "pipe:: invalid data found when processing input", err.Error())
}
