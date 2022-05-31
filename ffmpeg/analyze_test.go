package ffmpeg

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAnalyze(t *testing.T) {
	for _, item := range []struct {
		sample string
		report Report
	}{
		// audio
		{
			sample: "sample.aac",
			report: Report{
				Duration: 96.662622,
				Format: Format{
					Name:     "aac",
					Duration: 96.662622,
				},
				Streams: []Stream{
					{
						Type:       "audio",
						Codec:      "aac",
						Duration:   96.662622,
						Channels:   2,
						SampleRate: 44100,
					},
				},
			},
		},
		{
			sample: "sample.aiff",
			report: Report{
				Duration: 105.772948,
				Format: Format{
					Name:     "aiff",
					Duration: 105.772948,
				},
				Streams: []Stream{
					{
						Type:       "audio",
						Codec:      "pcm_s16be",
						Duration:   105.772948,
						Channels:   2,
						SampleRate: 44100,
					},
				},
			},
		},
		{
			sample: "sample.flac",
			report: Report{
				Duration: 105.772948,
				Format: Format{
					Name:     "flac",
					Duration: 105.772948,
				},
				Streams: []Stream{
					{
						Type:       "audio",
						Codec:      "flac",
						Duration:   105.772948,
						Channels:   2,
						SampleRate: 44100,
					},
				},
			},
		},
		{
			sample: "sample.m4a",
			report: Report{
				Duration: 105.797,
				Format: Format{
					Name:     "mov,mp4,m4a,3gp,3g2,mj2",
					Duration: 105.797,
				},
				Streams: []Stream{
					{
						Type:       "audio",
						Codec:      "aac",
						Duration:   105.772993,
						Channels:   2,
						SampleRate: 44100,
					},
				},
			},
		},
		{
			sample: "sample.mp2",
			report: Report{
				Duration: 105.795917,
				Format: Format{
					Name:     "mp3",
					Duration: 105.795917,
				},
				Streams: []Stream{
					{
						Type:       "audio",
						Codec:      "mp2",
						Duration:   105.795917,
						Channels:   2,
						SampleRate: 44100,
					},
				},
			},
		},
		{
			sample: "sample.mp3",
			report: Report{
				Duration: 105.822041,
				Format: Format{
					Name:     "mp3",
					Duration: 105.822041,
				},
				Streams: []Stream{
					{
						Type:       "audio",
						Codec:      "mp3",
						Duration:   105.822041,
						Channels:   2,
						SampleRate: 44100,
					},
				},
			},
		},
		{
			sample: "sample.ogg",
			report: Report{
				Duration: 105.772948,
				Format: Format{
					Name:     "ogg",
					Duration: 105.772948,
				}, Streams: []Stream{
					{
						Type:       "audio",
						Codec:      "vorbis",
						Duration:   105.772948,
						Channels:   2,
						SampleRate: 44100,
					},
				},
			},
		},
		{
			sample: "sample.wav",
			report: Report{
				Duration: 105.772948,
				Format: Format{
					Name:     "wav",
					Duration: 105.772948,
				},
				Streams: []Stream{
					{
						Type:       "audio",
						Codec:      "pcm_s16le",
						Duration:   105.772948,
						Channels:   2,
						SampleRate: 44100,
					},
				},
			},
		},
		{
			sample: "sample.wma",
			report: Report{
				Duration: 105.789,
				Format: Format{
					Name:     "asf",
					Duration: 105.789,
				},
				Streams: []Stream{
					{
						Type:       "audio",
						Codec:      "wmav2",
						Duration:   105.789,
						Channels:   2,
						SampleRate: 44100,
					},
				},
			},
		},
		// video
		{
			sample: "sample.hevc",
			report: Report{
				Duration: 28.23,
				Format: Format{
					Name:     "hevc",
					Duration: 0,
				},
				Streams: []Stream{
					{
						Type:      "video",
						Codec:     "hevc",
						Duration:  0,
						Width:     1280,
						Height:    720,
						FrameRate: 23.976023976023978,
					},
				},
				DidParse: true,
			},
		},
		{
			sample: "sample.avi",
			report: Report{
				Duration: 28.236542,
				Format: Format{
					Name:     "avi",
					Duration: 28.236542,
				},
				Streams: []Stream{
					{
						Type:      "video",
						Codec:     "mpeg4",
						Duration:  28.236542,
						Width:     1280,
						Height:    720,
						FrameRate: 23.976023976023978,
					},
				},
			},
		},
		{
			sample: "sample.mov",
			report: Report{
				Duration: 28.237,
				Format: Format{
					Name:     "mov,mp4,m4a,3gp,3g2,mj2",
					Duration: 28.237,
				},
				Streams: []Stream{
					{
						Type:      "video",
						Codec:     "h264",
						Duration:  28.236542,
						Width:     1280,
						Height:    720,
						FrameRate: 23.976023976023978,
					},
				},
			},
		},
		{
			sample: "sample.mp4",
			report: Report{
				Duration: 28.237,
				Format: Format{
					Name:     "mov,mp4,m4a,3gp,3g2,mj2",
					Duration: 28.237,
				},
				Streams: []Stream{
					{
						Type:      "video",
						Codec:     "h264",
						Duration:  28.236542,
						Width:     1280,
						Height:    720,
						FrameRate: 23.976023976023978,
					},
				},
			},
		},
		{
			sample: "sample.mpeg",
			report: Report{
				Duration: 28.236533,
				Format: Format{
					Name:     "mpeg",
					Duration: 28.236533,
				},
				Streams: []Stream{
					{
						Type:      "video",
						Codec:     "mpeg1video",
						Duration:  28.236533,
						Width:     1280,
						Height:    720,
						FrameRate: 23.976023976023978,
					},
				},
			},
		},
		{
			sample: "sample.mpg",
			report: Report{
				Duration: 28.27,
				Format: Format{
					Name:     "mpegvideo",
					Duration: 0,
				},
				Streams: []Stream{
					{
						Type:      "video",
						Codec:     "mpeg2video",
						Duration:  0,
						Width:     1280,
						Height:    720,
						FrameRate: 23.976023976023978,
					},
				},
				DidParse: true,
			},
		},
		{
			sample: "sample.webm",
			report: Report{
				Duration: 28.237,
				Format: Format{
					Name:     "matroska,webm",
					Duration: 28.237,
				},
				Streams: []Stream{
					{
						Type:      "video",
						Codec:     "vp9",
						Duration:  0,
						Width:     1280,
						Height:    720,
						FrameRate: 23.976023976023978,
					},
				},
			},
		},
		{
			sample: "sample.wmv",
			report: Report{
				Duration: 28.237,
				Format: Format{
					Name:     "asf",
					Duration: 28.237,
				},
				Streams: []Stream{
					{
						Type:      "video",
						Codec:     "msmpeg4v3",
						Duration:  28.237,
						Width:     1280,
						Height:    720,
						FrameRate: 23.976023976023978,
					},
				},
			},
		},
		// combined
		{
			sample: "combined_avc-aac.mov",
			report: Report{
				Duration: 7.2,
				Format: Format{
					Name:     "mov,mp4,m4a,3gp,3g2,mj2",
					Duration: 7.2,
				},
				Streams: []Stream{
					{
						Type:       "audio",
						Codec:      "aac",
						Duration:   7.2,
						Channels:   2,
						SampleRate: 48000,
					},
					{
						Type:      "video",
						Codec:     "h264",
						Duration:  7.2,
						Width:     1280,
						Height:    720,
						FrameRate: 25,
					},
					{
						Type:       "data",
						Codec:      "",
						Duration:   7.2,
						Channels:   0,
						SampleRate: 0,
						Width:      0,
						Height:     0,
					},
				},
			},
		},
		{
			sample: "combined_hevc-aac.mp4",
			report: Report{
				Duration: 7.189333,
				Format: Format{
					Name:     "mov,mp4,m4a,3gp,3g2,mj2",
					Duration: 7.189333,
				},
				Streams: []Stream{
					{
						Type:      "video",
						Codec:     "hevc",
						Duration:  7.16,
						Width:     1280,
						Height:    720,
						FrameRate: 25,
					},
					{
						Type:       "audio",
						Codec:      "aac",
						Duration:   7.16,
						Channels:   2,
						SampleRate: 48000,
					},
				},
			},
		},
		{
			sample: "combined_mpeg2.mpg",
			report: Report{
				Duration: 7.224,
				Format: Format{
					Name:     "mpeg",
					Duration: 7.224,
				}, Streams: []Stream{
					{
						Type:      "video",
						Codec:     "mpeg2video",
						Duration:  7.08,
						Width:     1280,
						Height:    720,
						FrameRate: 25,
					},
					{
						Type:       "audio",
						Codec:      "mp2",
						Duration:   7.224,
						Channels:   2,
						SampleRate: 48000,
					},
				},
			},
		},
		// image
		{
			sample: "sample.gif",
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
						Width:     1280,
						Height:    853,
						FrameRate: 10,
					},
				},
			},
		},
		{
			sample: "sample.jpg",
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
						Width:     1280,
						Height:    853,
						FrameRate: 25,
					},
				},
			},
		},
		{
			sample: "sample.png",
			report: Report{
				Format: Format{
					Name: "png_pipe",
				},
				Streams: []Stream{
					{
						Type:      "video",
						Codec:     "png",
						Width:     1280,
						Height:    853,
						FrameRate: 25,
					},
				},
			},
		},
	} {
		t.Run(item.sample, func(t *testing.T) {
			sample := loadSample(item.sample)
			defer sample.Close()

			report, err := Analyze(sample)
			assert.NoError(t, err)
			assert.Equal(t, &item.report, report)
		})
	}
}

func TestAnalyzePipe(t *testing.T) {
	sample := loadSample("sample.aac")
	defer sample.Close()

	buf, err := io.ReadAll(sample)
	assert.NoError(t, err)

	report, err := Analyze(bytes.NewReader(buf))
	assert.NoError(t, err)
	assert.Equal(t, &Report{
		Duration: 105.81,
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
	sample := loadSample("sample.mp3")
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
