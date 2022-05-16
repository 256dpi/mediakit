package mediakit

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFFProbe(t *testing.T) {
	for _, item := range []struct {
		ext string
		rep *Report
	}{
		// audio
		{
			ext: "aac",
			rep: &Report{
				Format: Format{
					Name:       "aac",
					LongName:   "raw ADTS AAC (Advanced Audio Coding)",
					ProbeScore: 51,
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
			ext: "aiff",
			rep: &Report{
				Format: Format{
					Name:       "aiff",
					LongName:   "Audio IFF",
					ProbeScore: 100,
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
			ext: "flac",
			rep: &Report{
				Format: Format{
					Name:       "flac",
					LongName:   "raw FLAC",
					ProbeScore: 100,
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
			ext: "m4a",
			rep: &Report{
				Format: Format{
					Name:       "mov,mp4,m4a,3gp,3g2,mj2",
					LongName:   "QuickTime / MOV",
					ProbeScore: 100,
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
			ext: "mp2",
			rep: &Report{
				Format: Format{
					Name:       "mp3",
					LongName:   "MP2/3 (MPEG audio layer 2/3)",
					ProbeScore: 51,
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
			ext: "mp3",
			rep: &Report{
				Format: Format{
					Name:       "mp3",
					LongName:   "MP2/3 (MPEG audio layer 2/3)",
					ProbeScore: 51,
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
			ext: "ogg",
			rep: &Report{
				Format: Format{
					Name:       "ogg",
					LongName:   "Ogg",
					ProbeScore: 100,
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
			ext: "wav",
			rep: &Report{
				Format: Format{
					Name:       "wav",
					LongName:   "WAV / WAVE (Waveform Audio)",
					ProbeScore: 99,
				},
				Streams: []Stream{
					{
						CodecName:     "pcm_s16le",
						CodecLongName: "PCM signed 16-bit little-endian",
						CodecType:     "audio",
						BitRate:       1411200,
						SampleRate:    44100,
						Channels:      2,
					},
				},
			},
		},
		{
			ext: "wma",
			rep: &Report{
				Format: Format{
					Name:       "asf",
					LongName:   "ASF (Advanced / Active Streaming Format)",
					ProbeScore: 100,
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
			ext: "hevc",
			rep: &Report{
				Format: Format{
					Name:       "hevc",
					LongName:   "raw HEVC video",
					ProbeScore: 51,
				},
				Streams: []Stream{
					{
						CodecName:     "hevc",
						CodecLongName: "H.265 / HEVC (High Efficiency Video Coding)",
						CodecType:     "video",
						BitRate:       0,
						Width:         1280,
						Height:        720,
					},
				},
			},
		},
		{
			ext: "avi",
			rep: &Report{
				Format: Format{
					Name:       "avi",
					LongName:   "AVI (Audio Video Interleaved)",
					ProbeScore: 100,
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
			ext: "mov",
			rep: &Report{
				Format: Format{
					Name:       "mov,mp4,m4a,3gp,3g2,mj2",
					LongName:   "QuickTime / MOV",
					ProbeScore: 100,
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
			ext: "mp4",
			rep: &Report{
				Format: Format{
					Name:       "mov,mp4,m4a,3gp,3g2,mj2",
					LongName:   "QuickTime / MOV",
					ProbeScore: 100,
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
			ext: "mpeg",
			rep: &Report{
				Format: Format{
					Name:       "mpeg",
					LongName:   "MPEG-PS (MPEG-2 Program Stream)",
					ProbeScore: 26,
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
			ext: "mpg",
			rep: &Report{
				Format: Format{
					Name:       "mpegvideo",
					LongName:   "raw MPEG video",
					ProbeScore: 51,
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
			ext: "webm",
			rep: &Report{
				Format: Format{
					Name:       "matroska,webm",
					LongName:   "Matroska / WebM",
					ProbeScore: 100,
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
			ext: "wmv",
			rep: &Report{
				Format: Format{
					Name:       "asf",
					LongName:   "ASF (Advanced / Active Streaming Format)",
					ProbeScore: 100,
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
		t.Run(item.ext, func(t *testing.T) {
			sample := loadSample(item.ext)
			defer sample.Close()

			report, err := FFProbe(sample)
			assert.NoError(t, err)
			assert.Equal(t, item.rep, report)
		})
	}
}

func TestFFProbeError(t *testing.T) {
	report, err := FFProbe(strings.NewReader("foo"))
	assert.Error(t, err)
	assert.Nil(t, report)
	assert.Equal(t, "invalid data found when processing input", err.Error())
}

func BenchmarkFFProbe(b *testing.B) {
	sample := loadSample("mp3")
	defer sample.Close()

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, sample)

	reader := bytes.NewReader(buf.Bytes())

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		reader.Reset(buf.Bytes())

		_, err := FFProbe(reader)
		if err != nil {
			panic(err)
		}
	}
}

func loadSample(ext string) io.ReadCloser {
	f, err := os.Open("./samples/sample." + ext)
	if err != nil {
		panic(err)
	}

	return f
}
