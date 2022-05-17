package ffmpeg

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvert(t *testing.T) {
	for _, item := range []struct {
		sample  string
		options ConvertOptions
		report  Report
	}{
		// samples
		{
			sample: "sample.wav",
			options: ConvertOptions{
				Preset: AudioMP3VBRStandard,
			},
			report: Report{
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
			},
		},
		{
			sample: "sample.mpeg",
			options: ConvertOptions{
				Preset:   VideoMP4H264AACFast,
				Duration: 1,
				Width:    1024,
				Height:   -1,
			},
			report: Report{
				Duration: 1.001,
				Format: Format{
					Name:       "mov,mp4,m4a,3gp,3g2,mj2",
					LongName:   "QuickTime / MOV",
					ProbeScore: 100,
					Duration:   1.001,
				},
				Streams: []Stream{
					{
						CodecName:     "h264",
						CodecLongName: "H.264 / AVC / MPEG-4 AVC / MPEG-4 part 10",
						CodecType:     "video",
						BitRate:       3887608,
						Duration:      1.001,
						Width:         1024,
						Height:        576,
						FrameRate:     23.976023976023978,
					},
				},
			},
		},
		{
			sample: "sample.jpg",
			options: ConvertOptions{
				Preset: ImagePNG,
				Width:  640,
				Height: -1,
			},
			report: Report{
				Duration: 0.04,
				Format: Format{
					Name:       "png_pipe",
					LongName:   "piped png sequence",
					ProbeScore: 99,
					Duration:   0,
				},
				Streams: []Stream{
					{
						CodecName:     "png",
						CodecLongName: "PNG (Portable Network Graphics) image",
						CodecType:     "video",
						BitRate:       0,
						Duration:      0,
						Width:         640,
						Height:        427,
						FrameRate:     25,
					},
				},
			},
		},
		// combined
		{
			sample: "combined_avc-aac.mov",
			options: ConvertOptions{
				Preset:   VideoMP4H264AACFast,
				Duration: 1,
			},
			report: Report{
				Duration: 1.08,
				Format: Format{
					Name:       "mov,mp4,m4a,3gp,3g2,mj2",
					LongName:   "QuickTime / MOV",
					ProbeScore: 100,
					Duration:   1.08,
				},
				Streams: []Stream{
					{
						CodecName:     "h264",
						CodecLongName: "H.264 / AVC / MPEG-4 AVC / MPEG-4 part 10",
						CodecType:     "video",
						BitRate:       394504,
						Duration:      1,
						Width:         1280,
						Height:        720,
						FrameRate:     25,
					},
					{
						CodecName:     "aac",
						CodecLongName: "AAC (Advanced Audio Coding)",
						CodecType:     "audio",
						BitRate:       25104,
						Duration:      1.08,
						SampleRate:    48000,
						Channels:      2,
					},
					{
						CodecType: "data",
						Duration:  1,
					},
				},
			},
		},
		{
			sample: "combined_hevc-aac.mp4",
			options: ConvertOptions{
				Preset:   VideoMP4H264AACFast,
				Duration: 1,
			},
			report: Report{
				Duration: 1.08,
				Format: Format{
					Name:       "mov,mp4,m4a,3gp,3g2,mj2",
					LongName:   "QuickTime / MOV",
					ProbeScore: 100,
					Duration:   1.08,
				},
				Streams: []Stream{
					{
						CodecName:     "h264",
						CodecLongName: "H.264 / AVC / MPEG-4 AVC / MPEG-4 part 10",
						CodecType:     "video",
						BitRate:       387656,
						Duration:      1,
						Width:         1280,
						Height:        720,
						FrameRate:     25,
					},
					{
						CodecName:     "aac",
						CodecLongName: "AAC (Advanced Audio Coding)",
						CodecType:     "audio",
						BitRate:       25259,
						Duration:      1.08,
						SampleRate:    48000,
						Channels:      2,
					},
				},
			},
		},
		{
			sample: "combined_mpeg2.mpg",
			options: ConvertOptions{
				Preset:   VideoMP4H264AACFast,
				Duration: 1,
			},
			report: Report{
				Duration: 1.08,
				Format: Format{
					Name:       "mov,mp4,m4a,3gp,3g2,mj2",
					LongName:   "QuickTime / MOV",
					ProbeScore: 100,
					Duration:   1.08,
				},
				Streams: []Stream{
					{
						CodecName:     "h264",
						CodecLongName: "H.264 / AVC / MPEG-4 AVC / MPEG-4 part 10",
						CodecType:     "video",
						BitRate:       397792,
						Duration:      1,
						Width:         1280,
						Height:        720,
						FrameRate:     25,
					},
					{
						CodecName:     "aac",
						CodecLongName: "AAC (Advanced Audio Coding)",
						CodecType:     "audio",
						BitRate:       35681,
						Duration:      1.08,
						SampleRate:    48000,
						Channels:      2,
					},
				},
			},
		},
		// images
		{
			sample: "sample.mp4",
			options: ConvertOptions{
				Preset: ImageJPEG,
				Start:  5,
			},
			report: Report{
				Duration: 0.04,
				Format: Format{
					Name:       "jpeg_pipe",
					LongName:   "piped jpeg sequence",
					ProbeScore: 51,
					Duration:   0,
				}, Streams: []Stream{
					{
						CodecName:     "mjpeg",
						CodecLongName: "Motion JPEG",
						CodecType:     "video",
						Width:         1280,
						Height:        720,
						FrameRate:     25,
					},
				},
			},
		},
		{
			sample: "sample.mp4",
			options: ConvertOptions{
				Preset: ImagePNG,
				Start:  5,
			},
			report: Report{
				Duration: 0.04,
				Format: Format{
					Name:       "png_pipe",
					LongName:   "piped png sequence",
					ProbeScore: 99,
					Duration:   0,
				}, Streams: []Stream{
					{
						CodecName:     "png",
						CodecLongName: "PNG (Portable Network Graphics) image",
						CodecType:     "video",
						Width:         1280,
						Height:        720,
						FrameRate:     25,
					},
				},
			},
		},
	} {
		t.Run(item.sample, func(t *testing.T) {
			sample := loadSample(item.sample)
			defer sample.Close()

			var out bytes.Buffer
			err := Convert(sample, &out, item.options)
			assert.NoError(t, err)

			r := bytes.NewReader(out.Bytes())
			report, err := Analyze(r, AnalyzeOptions{
				Reset: func() error {
					_, err := r.Seek(0, io.SeekStart)
					return err
				},
			})
			assert.NoError(t, err)
			assert.Equal(t, &item.report, report)
		})
	}
}

func TestConvertPipe(t *testing.T) {
	sample := loadSample("combined_hevc-aac.mp4")
	defer sample.Close()

	buf, err := io.ReadAll(sample)
	assert.NoError(t, err)

	var out bytes.Buffer
	r := bytes.NewReader(buf)
	err = Convert(r, &out, ConvertOptions{
		Preset: VideoMP4H264AACFast,
	})
	assert.NoError(t, err)
}

func TestConvertProgress(t *testing.T) {
	sample := loadSample("sample.mpeg")
	defer sample.Close()

	var out bytes.Buffer
	var progress []Progress
	err := Convert(sample, &out, ConvertOptions{
		Preset:   VideoMP4H264AACFast,
		Duration: 1,
		Progress: func(p Progress) {
			progress = append(progress, p)
		},
	})
	assert.NoError(t, err)
	assert.Len(t, progress, 2)
	assert.Equal(t, Progress{
		Duration: 0,
		Size:     36,
	}, progress[0])
	assert.Equal(t, Progress{
		Duration: 0.875917,
		Size:     775897,
	}, progress[1])
}

func TestConvertError(t *testing.T) {
	err := Convert(strings.NewReader("foo"), io.Discard, ConvertOptions{
		Preset: AudioMP3VBRStandard,
	})
	assert.Error(t, err)
	assert.Equal(t, "pipe:: invalid data found when processing input", err.Error())
}
