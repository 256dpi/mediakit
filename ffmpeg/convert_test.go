package ffmpeg

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertAudio(t *testing.T) {
	sample := loadSample("wav")
	defer sample.Close()

	var out bytes.Buffer
	err := Convert(sample, &out, ConvertOptions{
		Preset: AudioMP3VBRStandard,
	})
	assert.NoError(t, err)

	r := bytes.NewReader(out.Bytes())
	report, err := Analyze(r, AnalyzeOptions{
		Reset: func() error {
			_, err := r.Seek(0, io.SeekStart)
			return err
		},
	})
	assert.NoError(t, err)
	assert.Equal(t, &Report{
		Duration: 105.82,
		Format: Format{
			Name:       "mp3",
			LongName:   "MP2/3 (MPEG audio layer 2/3)",
			ProbeScore: 51,
			Duration:   0,
		}, Streams: []Stream{
			{
				CodecName:     "mp3",
				CodecLongName: "MP3 (MPEG audio layer 3)",
				CodecType:     "audio",
				BitRate:       163993,
				Duration:      0,
				SampleRate:    44100,
				Channels:      2,
			},
		},
	}, report)
}

func TestConvertVideo(t *testing.T) {
	sample := loadSample("mpeg")
	defer sample.Close()

	var out bytes.Buffer
	err := Convert(sample, &out, ConvertOptions{
		Preset:   VideoMP4H264Fast,
		Duration: 1,
	})
	assert.NoError(t, err)

	r := bytes.NewReader(out.Bytes())
	report, err := Analyze(r, AnalyzeOptions{
		Reset: func() error {
			_, err := r.Seek(0, io.SeekStart)
			return err
		},
	})
	assert.NoError(t, err)
	assert.Equal(t, &Report{
		Duration: 1,
		Format: Format{
			Name:       "mov,mp4,m4a,3gp,3g2,mj2",
			LongName:   "QuickTime / MOV",
			ProbeScore: 100,
			Duration:   1,
		},
		Streams: []Stream{
			{
				CodecName:     "h264",
				CodecLongName: "H.264 / AVC / MPEG-4 AVC / MPEG-4 part 10",
				CodecType:     "video",
				BitRate:       6068568,
				Duration:      1,
				Width:         1280,
				Height:        720,
			},
		},
	}, report)
}

func TestConvertError(t *testing.T) {
	err := Convert(strings.NewReader("foo"), io.Discard, ConvertOptions{
		Preset: AudioMP3VBRStandard,
	})
	assert.Error(t, err)
	assert.Equal(t, "pipe:: invalid data found when processing input", err.Error())
}
