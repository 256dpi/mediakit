package ffmpeg

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/256dpi/mediakit/samples"
)

func TestAnalyze(t *testing.T) {
	for _, item := range []struct {
		sample string
		report Report
	}{
		// audio
		{
			sample: samples.AudioAAC,
			report: Report{
				Duration: 2.127203,
				Format: Format{
					Name:     "aac",
					Duration: 2.127203,
				},
				Streams: []Stream{
					{
						Type:       "audio",
						Codec:      "aac",
						Duration:   2.127203,
						Channels:   2,
						SampleRate: 44100,
					},
				},
			},
		},
		{
			sample: samples.AudioAIFF,
			report: Report{
				Duration: 2.043356,
				Format: Format{
					Name:     "aiff",
					Duration: 2.043356,
				},
				Streams: []Stream{
					{
						Type:       "audio",
						Codec:      "pcm_s16be",
						Duration:   2.043356,
						Channels:   2,
						SampleRate: 44100,
					},
				},
			},
		},
		{
			sample: samples.AudioFLAC,
			report: Report{
				Duration: 2.115918,
				Format: Format{
					Name:     "flac",
					Duration: 2.115918,
				},
				Streams: []Stream{
					{
						Type:       "audio",
						Codec:      "flac",
						Duration:   2.115918,
						Channels:   2,
						SampleRate: 44100,
					},
				},
			},
		},
		{
			sample: samples.AudioMPEG2,
			report: Report{
				Duration: 2.115917,
				Format: Format{
					Name:     "mp3",
					Duration: 2.115917,
				},
				Streams: []Stream{
					{
						Type:       "audio",
						Codec:      "mp2",
						Duration:   2.115917,
						Channels:   2,
						SampleRate: 44100,
					},
				},
			},
		},
		{
			sample: samples.AudioMPEG3,
			report: Report{
				Duration: 2.123813,
				Format: Format{
					Name:     "mp3",
					Duration: 2.123813,
				},
				Streams: []Stream{
					{
						Type:       "audio",
						Codec:      "mp3",
						Duration:   2.123813,
						Channels:   2,
						SampleRate: 44100,
					},
				},
			},
		},
		{
			sample: samples.AudioMPEG4,
			report: Report{
				Duration: 2.116,
				Format: Format{
					Name:     "mov,mp4,m4a,3gp,3g2,mj2",
					Duration: 2.116,
				},
				Streams: []Stream{
					{
						Type:       "audio",
						Codec:      "aac",
						Duration:   2.115918,
						Channels:   2,
						SampleRate: 44100,
					},
				},
			},
		},
		{
			sample: samples.AudioOGG,
			report: Report{
				Duration: 2.115918,
				Format: Format{
					Name:     "ogg",
					Duration: 2.115918,
				}, Streams: []Stream{
					{
						Type:       "audio",
						Codec:      "vorbis",
						Duration:   2.115918,
						Channels:   2,
						SampleRate: 44100,
					},
				},
			},
		},
		{
			sample: samples.AudioWAV,
			report: Report{
				Duration: 2.043356,
				Format: Format{
					Name:     "wav",
					Duration: 2.043356,
				},
				Streams: []Stream{
					{
						Type:       "audio",
						Codec:      "pcm_s24le",
						Duration:   2.043356,
						Channels:   2,
						SampleRate: 44100,
					},
				},
			},
		},
		{
			sample: samples.AudioWMA,
			report: Report{
				Duration: 2.135,
				Format: Format{
					Name:     "asf",
					Duration: 2.135,
				},
				Streams: []Stream{
					{
						Type:       "audio",
						Codec:      "wmav2",
						Duration:   2.135,
						Channels:   2,
						SampleRate: 44100,
					},
				},
			},
		},
		// video
		{
			sample: samples.VideoAVI,
			report: Report{
				Duration: 2.136236,
				Format: Format{
					Name:     "avi",
					Duration: 2.136236,
				},
				Streams: []Stream{
					{
						Type:      "video",
						Codec:     "h264",
						Duration:  2.04,
						Width:     800,
						Height:    450,
						FrameRate: 25,
					},
					{
						Type:       "audio",
						Codec:      "aac",
						Duration:   2.136236,
						Channels:   2,
						SampleRate: 44100,
					},
				},
			},
		},
		{
			sample: samples.VideoFLV,
			report: Report{
				Duration: 2.069,
				Format: Format{
					Name:     "flv",
					Duration: 2.069,
				},
				Streams: []Stream{
					{
						Type:       "audio",
						Codec:      "mp3",
						Duration:   0,
						Channels:   2,
						SampleRate: 44100,
					},
					{
						Type:      "video",
						Codec:     "flv1",
						Duration:  0,
						Width:     800,
						Height:    450,
						FrameRate: 25,
					},
				},
			},
		},
		{
			sample: samples.VideoGIF,
			report: Report{
				Duration: 2,
				Format: Format{
					Name:     "gif",
					Duration: 2,
				},
				Streams: []Stream{
					{
						Type:      "video",
						Codec:     "gif",
						Duration:  2,
						Width:     800,
						Height:    450,
						FrameRate: 5,
					},
				},
			},
		},
		{
			sample: samples.VideoMKV,
			report: Report{
				Duration: 2.055,
				Format: Format{
					Name:     "matroska,webm",
					Duration: 2.055,
				},
				Streams: []Stream{
					{
						Type:      "video",
						Codec:     "hevc",
						Duration:  0,
						Width:     800,
						Height:    450,
						FrameRate: 25,
					},
					{
						Type:       "audio",
						Codec:      "ac3",
						Duration:   0,
						Channels:   2,
						SampleRate: 44100,
					},
				},
			},
		},
		{
			sample: samples.VideoMOV,
			report: Report{
				Duration: 2.044,
				Format: Format{
					Name:     "mov,mp4,m4a,3gp,3g2,mj2",
					Duration: 2.044,
				},
				Streams: []Stream{
					{
						Type:      "video",
						Codec:     "h264",
						Duration:  2.04,
						Width:     800,
						Height:    450,
						FrameRate: 25,
					},
					{
						Type:       "audio",
						Codec:      "aac",
						Duration:   2.042993,
						Channels:   2,
						SampleRate: 44100,
					},
				},
			},
		},
		{
			sample: samples.VideoMPEG,
			report: Report{
				Duration: 2.063678,
				Format: Format{
					Name:     "mpeg",
					Duration: 2.063678,
				},
				Streams: []Stream{
					{
						Type:      "video",
						Codec:     "mpeg1video",
						Duration:  2.04,
						Width:     800,
						Height:    450,
						FrameRate: 25,
					},
					{
						Type:       "audio",
						Codec:      "mp2",
						Duration:   2.063678,
						Channels:   2,
						SampleRate: 44100,
					},
				},
			},
		},
		{
			sample: samples.VideoMPEG2,
			report: Report{
				Duration: 2.063678,
				Format: Format{
					Name:     "mpeg",
					Duration: 2.063678,
				},
				Streams: []Stream{
					{
						Type:      "video",
						Codec:     "mpeg2video",
						Duration:  2.04,
						Width:     800,
						Height:    450,
						FrameRate: 25,
					},
					{
						Type:       "audio",
						Codec:      "mp2",
						Duration:   2.063678,
						Channels:   2,
						SampleRate: 44100,
					},
				},
			},
		},
		{
			sample: samples.VideoMPEG4,
			report: Report{
				Duration: 2.066578,
				Format: Format{
					Name:     "mov,mp4,m4a,3gp,3g2,mj2",
					Duration: 2.066578,
				},
				Streams: []Stream{
					{
						Type:      "video",
						Codec:     "h264",
						Duration:  2.04,
						Width:     800,
						Height:    450,
						FrameRate: 25,
					},
					{
						Type:       "audio",
						Codec:      "aac",
						Duration:   2.04,
						Channels:   2,
						SampleRate: 44100,
					},
				},
			},
		},
		{
			sample: samples.VideoWebM,
			report: Report{
				Duration: 2.05,
				Format: Format{
					Name:     "matroska,webm",
					Duration: 2.05,
				},
				Streams: []Stream{
					{
						Type:      "video",
						Codec:     "vp9",
						Duration:  0,
						Width:     800,
						Height:    450,
						FrameRate: 25,
					},
					{
						Type:       "audio",
						Codec:      "vorbis",
						Channels:   2,
						SampleRate: 44100,
					},
				},
			},
		},
		{
			sample: samples.VideoWMV,
			report: Report{
				Duration: 2.132,
				Format: Format{
					Name:     "asf",
					Duration: 2.132,
				},
				Streams: []Stream{
					{
						Type:      "video",
						Codec:     "wmv2",
						Duration:  2.086,
						Width:     800,
						Height:    450,
						FrameRate: 25,
					},
					{
						Type:       "audio",
						Codec:      "wmav2",
						Duration:   2.086,
						Channels:   2,
						SampleRate: 44100,
					},
				},
			},
		},
		// image
		{
			sample: samples.ImageGIF,
			report: Report{
				Duration: 0.1,
				Format: Format{
					Name:     "gif",
					Duration: 0.1,
				},
				Streams: []Stream{
					{
						Type:      "video",
						Codec:     "gif",
						Duration:  0.1,
						Width:     800,
						Height:    533,
						FrameRate: 10,
					},
				},
			},
		},
		{
			sample: samples.ImageJPEG,
			report: Report{
				Duration: 0.04,
				Format: Format{
					Name:     "image2",
					Duration: 0.04,
				},
				Streams: []Stream{
					{
						Type:      "video",
						Codec:     "mjpeg",
						Duration:  0.04,
						Width:     800,
						Height:    533,
						FrameRate: 25,
					},
				},
			},
		},
		{
			sample: samples.ImageJPEG2K,
			report: Report{
				Format: Format{
					Name: "j2k_pipe",
				},
				Streams: []Stream{
					{
						Type:      "video",
						Codec:     "jpeg2000",
						Width:     800,
						Height:    533,
						FrameRate: 25,
					},
				},
			},
		},
		{
			sample: samples.ImagePNG,
			report: Report{
				Format: Format{
					Name: "png_pipe",
				},
				Streams: []Stream{
					{
						Type:      "video",
						Codec:     "png",
						Width:     800,
						Height:    533,
						FrameRate: 25,
					},
				},
			},
		},
		{
			sample: samples.ImageTIFF,
			report: Report{
				Format: Format{
					Name: "tiff_pipe",
				},
				Streams: []Stream{
					{
						Type:      "video",
						Codec:     "tiff",
						Width:     800,
						Height:    533,
						FrameRate: 25,
					},
				},
			},
		},
		{
			sample: samples.ImageWebP,
			report: Report{
				Format: Format{
					Name: "webp_pipe",
				},
				Streams: []Stream{
					{
						Type:      "video",
						Codec:     "webp",
						Width:     800,
						Height:    533,
						FrameRate: 25,
					},
				},
			},
		},
	} {
		t.Run(item.sample, func(t *testing.T) {
			sample := samples.Buffer(item.sample)
			defer sample.Close()

			report, err := Analyze(sample)
			assert.NoError(t, err)
			assert.Equal(t, &item.report, report)
		})
	}
}

func TestAnalyzePipe(t *testing.T) {
	sample := samples.Load(samples.AudioAAC)
	defer sample.Close()

	buf, err := io.ReadAll(sample)
	assert.NoError(t, err)

	report, err := Analyze(bytes.NewReader(buf))
	assert.NoError(t, err)
	assert.True(t, report.Duration > 2)
	assert.Equal(t, &Report{
		Duration: report.Duration,
		Format: Format{
			Name:     "aac",
			Duration: 0,
		},
		Streams: []Stream{
			{
				Type:       "audio",
				Codec:      "aac",
				Duration:   0,
				Channels:   2,
				SampleRate: 44100,
			},
		},
		DidParse: true,
	}, report)
}

func TestAnalyzeError(t *testing.T) {
	report, err := Analyze(strings.NewReader("foo"))
	assert.Error(t, err)
	assert.Nil(t, report)
	assert.Equal(t, "invalid data found when processing input", err.Error())
}

func BenchmarkAnalyze(b *testing.B) {
	sample := samples.Load(samples.AudioMPEG3)
	defer sample.Close()

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, sample)

	reader := bytes.NewReader(buf.Bytes())

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		reader.Reset(buf.Bytes())

		_, err := Analyze(reader)
		if err != nil {
			panic(err)
		}
	}
}
