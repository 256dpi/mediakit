package ffmpeg

import (
	"bytes"
	"io"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/256dpi/mediakit/samples"
)

func TestConvertAudio(t *testing.T) {
	for _, sample := range samples.Audio() {
		t.Run(sample, func(t *testing.T) {
			file := samples.Buffer(sample)
			defer file.Close()

			out := tempFile(t)
			err := Convert(nil, file, out, ConvertOptions{
				Preset: AudioMP3VBRStandard,
			})
			assert.NoError(t, err)

			rewind(out)
			report, err := Analyze(nil, out)
			assert.NoError(t, err)
			assert.True(t, report.Duration > 2 && report.Duration < 2.2)
			assert.Equal(t, &Report{
				Duration: report.Duration,
				Format: Format{
					Name:     "mp3",
					Duration: report.Duration,
				}, Streams: []Stream{
					{
						Type:       "audio",
						Codec:      "mp3",
						Duration:   report.Duration,
						Channels:   2,
						SampleRate: 44100,
					},
				},
			}, report)
		})
	}
}

func TestConvertVideo(t *testing.T) {
	for _, sample := range samples.Video() {
		t.Run(sample, func(t *testing.T) {
			file := samples.Buffer(sample)
			defer file.Close()

			out := tempFile(t)
			err := Convert(nil, file, out, ConvertOptions{
				Preset: VideoMP4H264AACFast,
			})
			assert.NoError(t, err)

			width, height := 800, 450
			if sample == samples.VideoMPEG4R {
				width, height = 450, 800
			}

			rewind(out)
			report, err := Analyze(nil, out)
			assert.NoError(t, err)
			assert.True(t, report.Duration >= 2 && report.Duration < 2.3)
			assert.True(t, report.Format.Duration >= 2 && report.Format.Duration < 2.3)
			assert.True(t, report.Streams[0].Duration >= 2 && report.Streams[0].Duration < 2.3)
			assert.True(t, report.Streams[1].Duration >= 2 && report.Streams[1].Duration < 2.3)
			assert.Equal(t, &Report{
				Duration: report.Duration,
				Format: Format{
					Name:     "mov,mp4,m4a,3gp,3g2,mj2",
					Duration: report.Format.Duration,
				}, Streams: []Stream{
					{
						Type:        "video",
						Codec:       "h264",
						Duration:    report.Streams[0].Duration,
						Width:       width,
						Height:      height,
						FrameRate:   25,
						PixelFormat: "yuv420p",
						ColorSpace:  "bt709",
					},
					{
						Type:       "audio",
						Codec:      "aac",
						Duration:   report.Streams[1].Duration,
						Channels:   2,
						SampleRate: 44100,
					},
				},
			}, report)
		})
	}
}

func TestConvertImage(t *testing.T) {
	for _, sample := range []string{
		samples.ImageGIF,
		// samples.ImageHEIF,
		samples.ImageJPEG,
		samples.ImageJPEG2K,
		// samples.ImagePDF,
		samples.ImagePNG,
		samples.ImageTIFF,
		samples.ImageWebP,
	} {
		t.Run(sample, func(t *testing.T) {
			t.Run("JPG", func(t *testing.T) {
				sample := samples.Buffer(sample)
				defer sample.Close()

				out := tempFile(t)
				err := Convert(nil, sample, out, ConvertOptions{
					Preset: ImageJPEG,
				})
				assert.NoError(t, err)

				rewind(out)
				report, err := Analyze(nil, out)
				assert.NoError(t, err)
				assert.Equal(t, &Report{
					Duration: 0,
					Format: Format{
						Name:     "jpeg_pipe",
						Duration: 0,
					}, Streams: []Stream{
						{
							Type:        "video",
							Codec:       "mjpeg",
							Duration:    0,
							Width:       800,
							Height:      533,
							FrameRate:   25,
							PixelFormat: "yuvj444p",
							ColorSpace:  "bt470bg",
						},
					},
				}, report)
			})

			t.Run("PNG", func(t *testing.T) {
				sample := samples.Buffer(sample)
				defer sample.Close()

				out := tempFile(t)
				err := Convert(nil, sample, out, ConvertOptions{
					Preset: ImagePNG,
				})
				assert.NoError(t, err)

				rewind(out)
				report, err := Analyze(nil, out)
				assert.NoError(t, err)
				assert.Equal(t, &Report{
					Duration: 0,
					Format: Format{
						Name:     "png_pipe",
						Duration: 0,
					}, Streams: []Stream{
						{
							Type:        "video",
							Codec:       "png",
							Duration:    0,
							Width:       800,
							Height:      533,
							FrameRate:   25,
							PixelFormat: "rgb24",
							ColorSpace:  "gbr",
						},
					},
				}, report)
			})

			t.Run("WebP", func(t *testing.T) {
				sample := samples.Buffer(sample)
				defer sample.Close()

				out := tempFile(t)
				err := Convert(nil, sample, out, ConvertOptions{
					Preset: ImageWebP,
				})
				assert.NoError(t, err)

				rewind(out)
				report, err := Analyze(nil, out)
				assert.NoError(t, err)
				assert.Equal(t, &Report{
					Duration: 0,
					Format: Format{
						Name:     "webp_pipe",
						Duration: 0,
					}, Streams: []Stream{
						{
							Type:        "video",
							Codec:       "webp",
							Duration:    0,
							Width:       800,
							Height:      533,
							FrameRate:   25,
							PixelFormat: "yuv420p",
							ColorSpace:  "bt470bg",
						},
					},
				}, report)
			})
		})
	}
}

func TestConvertAnimation(t *testing.T) {
	for _, sample := range []string{
		samples.AnimationGIF,
		samples.VideoMPEG4,
	} {
		t.Run(sample, func(t *testing.T) {
			t.Run("GIF", func(t *testing.T) {
				buf := samples.Buffer(sample)
				defer buf.Close()

				out := tempFile(t)
				err := Convert(nil, buf, out, ConvertOptions{
					Preset: AnimationGIF,
				})
				assert.NoError(t, err)

				duration := 2.0
				height := 450
				frameRate := 5
				if sample == samples.VideoMPEG4 {
					duration = 2.04
					height = 450
					frameRate = 25
				}

				rewind(out)
				report, err := Analyze(nil, out)
				assert.NoError(t, err)
				assert.Equal(t, &Report{
					Duration: duration,
					Format: Format{
						Name:     "gif",
						Duration: duration,
					}, Streams: []Stream{
						{
							Type:        "video",
							Codec:       "gif",
							Duration:    duration,
							Width:       800,
							Height:      height,
							FrameRate:   FrameRate(frameRate),
							PixelFormat: "bgra",
						},
					},
				}, report)
			})

			t.Run("WebP", func(t *testing.T) {
				buf := samples.Buffer(sample)
				defer buf.Close()

				out := tempFile(t)
				err := Convert(nil, buf, out, ConvertOptions{
					Preset: AnimationWebP,
				})
				assert.NoError(t, err)

				rewind(out)
				report, err := Analyze(nil, out)
				assert.NoError(t, err)
				assert.Equal(t, &Report{
					Duration: 0,
					Format: Format{
						Name:     "webp_pipe",
						Duration: 0,
					}, Streams: []Stream{
						{
							Type:      "video",
							Codec:     "webp",
							FrameRate: 25,
						},
					},
				}, report)
			})
		})
	}
}

func TestConvertExtract(t *testing.T) {
	for _, sample := range samples.Video() {
		t.Run(sample, func(t *testing.T) {
			file := samples.Buffer(sample)
			defer file.Close()

			out := tempFile(t)
			err := Convert(nil, file, out, ConvertOptions{
				Preset: ImagePNG,
				Start:  1,
			})
			assert.NoError(t, err)

			width, height := 800, 450
			if sample == samples.VideoMPEG4R {
				width, height = 450, 800
			}

			rewind(out)
			report, err := Analyze(nil, out)
			assert.NoError(t, err)
			assert.Equal(t, &Report{
				Duration: 0,
				Format: Format{
					Name:     "png_pipe",
					Duration: 0,
				}, Streams: []Stream{
					{
						Type:        "video",
						Codec:       "png",
						Duration:    0,
						Width:       width,
						Height:      height,
						FrameRate:   25,
						PixelFormat: "rgb24",
						ColorSpace:  "gbr",
					},
				},
			}, report)
		})
	}
}

func TestConvertOptions(t *testing.T) {
	for i, item := range []struct {
		sample string
		opts   ConvertOptions
		report Report
	}{
		// audio
		{
			sample: samples.AudioAIFF,
			opts: ConvertOptions{
				Preset:     AudioMP3VBRStandard,
				Start:      1,
				Duration:   0.5,
				SampleRate: 16000,
			},
			report: Report{
				Format: Format{
					Name: "mp3",
				}, Streams: []Stream{
					{
						Type:       "audio",
						Codec:      "mp3",
						Channels:   2,
						SampleRate: 16000,
					},
				},
				DidScan: true,
			},
		},
		// video
		{
			sample: samples.VideoMOV,
			opts: ConvertOptions{
				Preset:     VideoMP4H264AACFast,
				Start:      1,
				Duration:   0.5,
				Width:      256,
				Height:     -1,
				FrameRate:  10,
				SampleRate: 16000,
			},
			report: Report{
				Format: Format{
					Name: "mov,mp4,m4a,3gp,3g2,mj2",
				},
				Streams: []Stream{
					{
						Type:        "video",
						Codec:       "h264",
						Width:       256,
						Height:      144,
						FrameRate:   10,
						PixelFormat: "yuv420p",
						ColorSpace:  "bt709",
					},
					{
						Type:       "audio",
						Codec:      "aac",
						Channels:   2,
						SampleRate: 16000,
					},
				},
			},
		},
	} {
		t.Run(strconv.Itoa(i)+"-"+item.sample, func(t *testing.T) {
			sample := samples.Buffer(item.sample)
			defer sample.Close()

			var buf bytes.Buffer
			err := Convert(nil, sample, &buf, item.opts)
			assert.NoError(t, err)

			report, err := Analyze(nil, bytes.NewReader(buf.Bytes()))
			assert.NoError(t, err)
			assert.True(t, report.Duration > 0)
			report.Duration = 0
			if item.opts.Preset != AudioMP3VBRStandard {
				assert.True(t, report.Format.Duration > 0)
				report.Format.Duration = 0
				for i, stream := range report.Streams {
					assert.True(t, stream.Duration > 0)
					report.Streams[i].Duration = 0
				}
			}
			assert.Equal(t, &item.report, report)
		})
	}
}

func TestConvertPipe(t *testing.T) {
	sample := samples.Load(samples.VideoMPEG4)
	defer sample.Close()

	buf, err := io.ReadAll(sample)
	assert.NoError(t, err)

	var out bytes.Buffer
	r := bytes.NewReader(buf)
	err = Convert(nil, r, &out, ConvertOptions{
		Preset: VideoMP4H264AACFast,
	})
	assert.NoError(t, err)
}

func TestConvertProgress(t *testing.T) {
	sample := samples.Load(samples.VideoMPEG2)
	defer sample.Close()

	out := tempFile(t)
	var progress []Progress
	err := Convert(nil, sample, out, ConvertOptions{
		Preset:   VideoMP4H264AACFast,
		Duration: 1,
		ProgressFunc: func(p Progress) {
			progress = append(progress, p)
		},
		ProgressRate: time.Second,
	})
	assert.NoError(t, err)
	assert.True(t, len(progress) >= 1)
	assert.True(t, progress[len(progress)-1].Duration > 0)
	assert.True(t, progress[len(progress)-1].Size > 36)
}

func TestConvertError(t *testing.T) {
	err := Convert(nil, strings.NewReader("foo"), io.Discard, ConvertOptions{
		Preset: AudioMP3VBRStandard,
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid data found when processing input")
}
