package ffmpeg

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"

	"github.com/256dpi/mediakit/samples"
)

func TestAnalyzeAudio(t *testing.T) {
	for _, item := range []struct {
		sample string
		format string
		codec  string
	}{
		{
			sample: samples.AudioAAC,
			format: "aac",
			codec:  "aac",
		},
		{
			sample: samples.AudioAIFF,
			format: "aiff",
			codec:  "pcm_s16be",
		},
		{
			sample: samples.AudioFLAC,
			format: "flac",
			codec:  "flac",
		},
		{
			sample: samples.AudioMPEG3,
			format: "mp3",
			codec:  "mp3",
		},
		{
			sample: samples.AudioMPEG4,
			format: "mov,mp4,m4a,3gp,3g2,mj2",
			codec:  "aac",
		},
		{
			sample: samples.AudioOGG,
			format: "ogg",
			codec:  "vorbis",
		},
		{
			sample: samples.AudioWAV,
			format: "wav",
			codec:  "pcm_s24le",
		},
		{
			sample: samples.AudioWMA,
			format: "asf",
			codec:  "wmav2",
		},
	} {
		t.Run(item.sample, func(t *testing.T) {
			sample := samples.Buffer(item.sample)
			defer sample.Close()

			report, err := Analyze(nil, sample)
			assert.NoError(t, err)
			assert.True(t, report.Duration > 2)
			assert.True(t, report.Format.Duration > 2)
			assert.True(t, report.Streams[0].Duration > 2)
			assert.Equal(t, &Report{
				Duration: report.Duration,
				Format: Format{
					Name:     item.format,
					Duration: report.Format.Duration,
				},
				Streams: []Stream{
					{
						Type:       "audio",
						Codec:      item.codec,
						Duration:   report.Streams[0].Duration,
						Channels:   2,
						SampleRate: 44100,
					},
				},
			}, report)
		})
	}
}

func TestAnalyzeVideo(t *testing.T) {
	for _, item := range []struct {
		sample  string
		format  string
		vCodec  string
		aCodec  string
		rotated bool
		pixFmt  string
		colSpc  string
	}{
		{
			sample: samples.VideoAVI,
			format: "avi",
			vCodec: "h264",
			aCodec: "aac",
			pixFmt: "yuv420p",
			colSpc: "bt709",
		},
		{
			sample: samples.VideoFLV,
			format: "flv",
			vCodec: "flv1",
			aCodec: "mp3",
			pixFmt: "yuv420p",
		},
		{
			sample: samples.VideoMKV,
			format: "matroska,webm",
			vCodec: "hevc",
			aCodec: "ac3",
			pixFmt: "yuv420p",
			colSpc: "smpte170m",
		},
		{
			sample: samples.VideoMOV,
			format: "mov,mp4,m4a,3gp,3g2,mj2",
			vCodec: "h264",
			aCodec: "aac",
			pixFmt: "yuv420p",
			colSpc: "bt709",
		},
		{
			sample: samples.VideoMPEG,
			format: "mpeg",
			vCodec: "mpeg1video",
			aCodec: "mp2",
			pixFmt: "yuv420p",
		},
		{
			sample: samples.VideoMPEG2,
			format: "mpeg",
			vCodec: "mpeg2video",
			aCodec: "mp2",
			pixFmt: "yuv420p",
			colSpc: "bt709",
		},
		{
			sample: samples.VideoMPEG4,
			format: "mov,mp4,m4a,3gp,3g2,mj2",
			vCodec: "h264",
			aCodec: "aac",
			pixFmt: "yuv420p",
			colSpc: "bt709",
		},
		{
			sample:  samples.VideoMPEG4R,
			format:  "mov,mp4,m4a,3gp,3g2,mj2",
			vCodec:  "h264",
			aCodec:  "aac",
			rotated: true,
			pixFmt:  "yuv420p",
			colSpc:  "bt709",
		},
		{
			sample: samples.VideoOGG,
			format: "ogg",
			vCodec: "theora",
			aCodec: "flac",
			pixFmt: "yuv420p",
		},
		{
			sample: samples.VideoWebM,
			format: "matroska,webm",
			vCodec: "vp9",
			aCodec: "vorbis",
			pixFmt: "yuv420p",
			colSpc: "bt709",
		},
		{
			sample: samples.VideoWMV,
			format: "asf",
			vCodec: "wmv2",
			aCodec: "wmav2",
			pixFmt: "yuv420p",
		},
	} {
		t.Run(item.sample, func(t *testing.T) {
			sample := samples.Buffer(item.sample)
			defer sample.Close()

			report, err := Analyze(nil, sample)
			if item.format == "flv" {
				report.Streams = lo.Reverse(report.Streams)
			}
			assert.NoError(t, err)
			assert.True(t, report.Duration >= 2)
			assert.True(t, report.Format.Duration >= 2)
			if !lo.Contains([]string{"flv", "matroska,webm"}, item.format) {
				assert.True(t, report.Streams[0].Duration >= 2, report.Streams[0].Duration)
				assert.True(t, report.Streams[1].Duration >= 2, report.Streams[1].Duration)
			}

			width, height := 800, 450
			if item.rotated {
				width, height = 450, 800
			}

			frameRate := 25
			if item.sample == samples.VideoMPEG {
				frameRate = 50
			}

			assert.Equal(t, &Report{
				Duration: report.Duration,
				Format: Format{
					Name:     item.format,
					Duration: report.Format.Duration,
				},
				Streams: []Stream{
					{
						Type:        "video",
						Codec:       item.vCodec,
						Duration:    report.Streams[0].Duration,
						Width:       width,
						Height:      height,
						FrameRate:   FrameRate(frameRate),
						PixelFormat: item.pixFmt,
						ColorSpace:  item.colSpc,
					},
					{
						Type:       "audio",
						Codec:      item.aCodec,
						Duration:   report.Streams[1].Duration,
						Channels:   2,
						SampleRate: 44100,
					},
				},
			}, report)
		})
	}
}

func TestAnalyzeAnimation(t *testing.T) {
	for _, item := range []struct {
		sample string
		format string
		vCodec string
		pixFmt string
		colSpc string
	}{
		{
			sample: samples.AnimationGIF,
			format: "gif",
			vCodec: "gif",
			pixFmt: "bgra",
		},
		{
			sample: samples.AnimationWebP,
			format: "webp_pipe",
			vCodec: "webp",
		},
	} {
		t.Run(item.sample, func(t *testing.T) {
			sample := samples.Buffer(item.sample)
			defer sample.Close()

			// Note: ffmpeg does not properly detect WebP animations

			report, err := Analyze(nil, sample)
			assert.NoError(t, err)
			assert.Equal(t, report.Duration >= 2, item.sample != samples.AnimationWebP)
			assert.Equal(t, report.Format.Duration >= 2, item.sample != samples.AnimationWebP)

			width, height := 800, 450
			if item.sample == samples.AnimationWebP {
				width, height = 0, 0
			}

			frameRate := 25
			if item.sample == samples.AnimationGIF {
				frameRate = 5
			}

			assert.Equal(t, &Report{
				Duration: report.Duration,
				Format: Format{
					Name:     item.format,
					Duration: report.Format.Duration,
				},
				Streams: []Stream{
					{
						Type:        "video",
						Codec:       item.vCodec,
						Duration:    report.Streams[0].Duration,
						Width:       width,
						Height:      height,
						FrameRate:   FrameRate(frameRate),
						PixelFormat: item.pixFmt,
						ColorSpace:  item.colSpc,
					},
				},
			}, report)
		})
	}
}

func TestAnalyzeImage(t *testing.T) {
	for _, item := range []struct {
		sample string
		format string
		codec  string
		pixFmt string
		colSpc string
	}{
		{
			sample: samples.ImageGIF,
			format: "gif",
			codec:  "gif",
			pixFmt: "bgra",
		},
		{
			sample: samples.ImageJPEG,
			format: "jpeg_pipe",
			codec:  "mjpeg",
			pixFmt: "yuvj444p",
			colSpc: "bt470bg",
		},
		{
			sample: samples.ImageJPEG2K,
			format: "j2k_pipe",
			codec:  "jpeg2000",
			pixFmt: "rgb24",
		},
		{
			sample: samples.ImagePNG,
			format: "png_pipe",
			codec:  "png",
			pixFmt: "rgb24",
			colSpc: "gbr",
		},
		{
			sample: samples.ImageTIFF,
			format: "tiff_pipe",
			codec:  "tiff",
			pixFmt: "rgba",
		},
		{
			sample: samples.ImageWebP,
			format: "webp_pipe",
			codec:  "webp",
			pixFmt: "yuv420p",
			colSpc: "bt470bg",
		},
	} {
		t.Run(item.sample, func(t *testing.T) {
			sample := samples.Buffer(item.sample)
			defer sample.Close()

			report, err := Analyze(nil, sample)
			assert.NoError(t, err)
			assert.Equal(t, &Report{
				Duration: report.Duration,
				Format: Format{
					Name:     item.format,
					Duration: report.Format.Duration,
				},
				Streams: []Stream{
					{
						Type:        "video",
						Codec:       item.codec,
						Duration:    report.Streams[0].Duration,
						Width:       800,
						Height:      533,
						FrameRate:   report.Streams[0].FrameRate,
						PixelFormat: item.pixFmt,
						ColorSpace:  item.colSpc,
					},
				},
			}, report)
		})
	}
}

func TestAnalyzePipe(t *testing.T) {
	sample := samples.Load(samples.AudioAAC)
	defer sample.Close()

	buf, err := io.ReadAll(sample)
	assert.NoError(t, err)

	report, err := Analyze(nil, bytes.NewReader(buf))
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
		DidScan: true,
	}, report)
}

func TestAnalyzeError(t *testing.T) {
	report, err := Analyze(nil, strings.NewReader("foo"))
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

		_, err := Analyze(nil, reader)
		if err != nil {
			panic(err)
		}
	}
}
