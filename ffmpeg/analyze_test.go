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
				Duration: 105.81,
				Format: Format{
					Name:       "aac",
					LongName:   "raw ADTS AAC (Advanced Audio Coding)",
					ProbeScore: 51,
					Duration:   0,
				},
				Streams: []Stream{
					{
						CodecName:     "aac",
						CodecLongName: "AAC (Advanced Audio Coding)",
						CodecType:     "audio",
						BitRate:       145531,
						Duration:      0,
						SampleRate:    44100,
						Channels:      2,
					},
				},
			},
		},
		{
			sample: "sample.aiff",
			report: Report{
				Duration: 105.772948,
				Format: Format{
					Name:       "aiff",
					LongName:   "Audio IFF",
					ProbeScore: 100,
					Duration:   105.772948,
				},
				Streams: []Stream{
					{
						CodecName:     "pcm_s16be",
						CodecLongName: "PCM signed 16-bit big-endian",
						CodecType:     "audio",
						BitRate:       1411200,
						Duration:      105.772948,
						SampleRate:    44100,
						Channels:      2,
					},
				},
			},
		},
		{
			sample: "sample.flac",
			report: Report{
				Duration: 105.772948,
				Format: Format{
					Name:       "flac",
					LongName:   "raw FLAC",
					ProbeScore: 100,
					Duration:   105.772948,
				},
				Streams: []Stream{
					{
						CodecName:     "flac",
						CodecLongName: "FLAC (Free Lossless Audio Codec)",
						CodecType:     "audio",
						BitRate:       0,
						Duration:      105.772948,
						SampleRate:    44100,
						Channels:      2,
					},
				},
			},
		},
		{
			sample: "sample.m4a",
			report: Report{
				Duration: 105.797,
				Format: Format{
					Name:       "mov,mp4,m4a,3gp,3g2,mj2",
					LongName:   "QuickTime / MOV",
					ProbeScore: 100,
					Duration:   105.797,
				},
				Streams: []Stream{
					{
						CodecName:     "aac",
						CodecLongName: "AAC (Advanced Audio Coding)",
						CodecType:     "audio",
						BitRate:       130554,
						Duration:      105.772993,
						SampleRate:    44100,
						Channels:      2,
					},
				},
			},
		},
		{
			sample: "sample.mp2",
			report: Report{
				Duration: 105.79,
				Format: Format{
					Name:       "mp3",
					LongName:   "MP2/3 (MPEG audio layer 2/3)",
					ProbeScore: 51,
					Duration:   0,
				},
				Streams: []Stream{
					{
						CodecName:     "mp2",
						CodecLongName: "MP2 (MPEG audio layer 2)",
						CodecType:     "audio",
						BitRate:       384000,
						Duration:      0,
						SampleRate:    44100,
						Channels:      2,
					},
				},
			},
		},
		{
			sample: "sample.mp3",
			report: Report{
				Duration: 105.79,
				Format: Format{
					Name:       "mp3",
					LongName:   "MP2/3 (MPEG audio layer 2/3)",
					ProbeScore: 51,
					Duration:   0,
				},
				Streams: []Stream{
					{
						CodecName:     "mp3",
						CodecLongName: "MP3 (MPEG audio layer 3)",
						CodecType:     "audio",
						BitRate:       128000,
						Duration:      0,
						SampleRate:    44100,
						Channels:      2,
					},
				},
			},
		},
		{
			sample: "sample.ogg",
			report: Report{
				Duration: 105.77,
				Format: Format{
					Name:       "ogg",
					LongName:   "Ogg",
					ProbeScore: 100,
					Duration:   0,
				}, Streams: []Stream{
					{
						CodecName:     "vorbis",
						CodecLongName: "Vorbis",
						CodecType:     "audio",
						BitRate:       112000,
						Duration:      0,
						SampleRate:    44100,
						Channels:      2,
					},
				},
			},
		},
		{
			sample: "sample.wav",
			report: Report{
				Duration: 105.77,
				Format: Format{
					Name:       "wav",
					LongName:   "WAV / WAVE (Waveform Audio)",
					ProbeScore: 99,
					Duration:   0,
				},
				Streams: []Stream{
					{
						CodecName:     "pcm_s16le",
						CodecLongName: "PCM signed 16-bit little-endian",
						CodecType:     "audio",
						BitRate:       1411200,
						Duration:      0,
						SampleRate:    44100,
						Channels:      2,
					},
				},
			},
		},
		{
			sample: "sample.wma",
			report: Report{
				Duration: 105.789,
				Format: Format{
					Name:       "asf",
					LongName:   "ASF (Advanced / Active Streaming Format)",
					ProbeScore: 100,
					Duration:   105.789,
				},
				Streams: []Stream{
					{
						CodecName:     "wmav2",
						CodecLongName: "Windows Media Audio 2",
						CodecType:     "audio",
						BitRate:       128000,
						Duration:      105.789,
						SampleRate:    44100,
						Channels:      2,
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
					Name:       "hevc",
					LongName:   "raw HEVC video",
					ProbeScore: 51,
					Duration:   0,
				},
				Streams: []Stream{
					{
						CodecName:     "hevc",
						CodecLongName: "H.265 / HEVC (High Efficiency Video Coding)",
						CodecType:     "video",
						BitRate:       0,
						Duration:      0,
						Width:         1280,
						Height:        720,
					},
				},
			},
		},
		{
			sample: "sample.avi",
			report: Report{
				Duration: 28.236542,
				Format: Format{
					Name:       "avi",
					LongName:   "AVI (Audio Video Interleaved)",
					ProbeScore: 100,
					Duration:   28.236542,
				},
				Streams: []Stream{
					{
						CodecName:     "mpeg4",
						CodecLongName: "MPEG-4 part 2",
						CodecType:     "video",
						BitRate:       0,
						Duration:      28.236542,
						Width:         1280,
						Height:        720,
					},
				},
			},
		},
		{
			sample: "sample.mov",
			report: Report{
				Duration: 28.237,
				Format: Format{
					Name:       "mov,mp4,m4a,3gp,3g2,mj2",
					LongName:   "QuickTime / MOV",
					ProbeScore: 100,
					Duration:   28.237,
				},
				Streams: []Stream{
					{
						CodecName:     "h264",
						CodecLongName: "H.264 / AVC / MPEG-4 AVC / MPEG-4 part 10",
						CodecType:     "video",
						BitRate:       4937429,
						Duration:      28.236542,
						Width:         1280,
						Height:        720,
					},
				},
			},
		},
		{
			sample: "sample.mp4",
			report: Report{
				Duration: 28.237,
				Format: Format{
					Name:       "mov,mp4,m4a,3gp,3g2,mj2",
					LongName:   "QuickTime / MOV",
					ProbeScore: 100,
					Duration:   28.237,
				},
				Streams: []Stream{
					{
						CodecName:     "h264",
						CodecLongName: "H.264 / AVC / MPEG-4 AVC / MPEG-4 part 10",
						CodecType:     "video",
						BitRate:       4937429,
						Duration:      28.236542,
						Width:         1280,
						Height:        720,
					},
				},
			},
		},
		{
			sample: "sample.mpeg",
			report: Report{
				Duration: 28.23,
				Format: Format{
					Name:       "mpeg",
					LongName:   "MPEG-PS (MPEG-2 Program Stream)",
					ProbeScore: 26,
					Duration:   0,
				},
				Streams: []Stream{
					{
						CodecName:     "mpeg1video",
						CodecLongName: "MPEG-1 video",
						CodecType:     "video",
						BitRate:       104857200,
						Duration:      0,
						Width:         1280,
						Height:        720,
					},
				},
			},
		},
		{
			sample: "sample.mpg",
			report: Report{
				Duration: 28.27,
				Format: Format{
					Name:       "mpegvideo",
					LongName:   "raw MPEG video",
					ProbeScore: 51,
					Duration:   0,
				},
				Streams: []Stream{
					{
						CodecName:     "mpeg2video",
						CodecLongName: "MPEG-2 video",
						CodecType:     "video",
						BitRate:       0,
						Duration:      0,
						Width:         1280,
						Height:        720,
					},
				},
			},
		},
		{
			sample: "sample.webm",
			report: Report{
				Duration: 28.237,
				Format: Format{
					Name:       "matroska,webm",
					LongName:   "Matroska / WebM",
					ProbeScore: 100,
					Duration:   28.237,
				},
				Streams: []Stream{
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
			},
		},
		{
			sample: "sample.wmv",
			report: Report{
				Duration: 28.237,
				Format: Format{
					Name:       "asf",
					LongName:   "ASF (Advanced / Active Streaming Format)",
					ProbeScore: 100,
					Duration:   28.237,
				},
				Streams: []Stream{
					{
						CodecName:     "msmpeg4v3",
						CodecLongName: "MPEG-4 part 2 Microsoft variant version 3",
						CodecType:     "video",
						BitRate:       0,
						Duration:      28.237,
						Width:         1280,
						Height:        720,
					},
				},
			},
		},
	} {
		t.Run(item.sample, func(t *testing.T) {
			sample := loadSample(item.sample)
			defer sample.Close()

			report, err := Analyze(sample, AnalyzeOptions{
				Reset: func() error {
					_, err := sample.Seek(0, io.SeekStart)
					return err
				},
			})
			assert.NoError(t, err)
			assert.Equal(t, &item.report, report)
		})
	}
}

func TestAnalyzeError(t *testing.T) {
	report, err := Analyze(strings.NewReader("foo"), AnalyzeOptions{})
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

		_, err := Analyze(reader, AnalyzeOptions{})
		if err != nil {
			panic(err)
		}
	}
}
