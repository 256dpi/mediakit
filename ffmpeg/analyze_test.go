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
			sample: samples.AudioMPEG2,
			format: "mp3",
			codec:  "mp2",
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

			report, err := Analyze(sample)
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
		sample string
		format string
		vCodec string
		aCodec string
	}{
		{
			sample: samples.VideoAVI,
			format: "avi",
			vCodec: "h264",
			aCodec: "aac",
		},
		{
			sample: samples.VideoFLV,
			format: "flv",
			vCodec: "flv1",
			aCodec: "mp3",
		},
		{
			sample: samples.VideoGIF,
			format: "gif",
			vCodec: "gif",
			aCodec: "",
		},
		{
			sample: samples.VideoMKV,
			format: "matroska,webm",
			vCodec: "hevc",
			aCodec: "ac3",
		},
		{
			sample: samples.VideoMOV,
			format: "mov,mp4,m4a,3gp,3g2,mj2",
			vCodec: "h264",
			aCodec: "aac",
		},
		{
			sample: samples.VideoMPEG,
			format: "mpeg",
			vCodec: "mpeg1video",
			aCodec: "mp2",
		},
		{
			sample: samples.VideoMPEG2,
			format: "mpeg",
			vCodec: "mpeg2video",
			aCodec: "mp2",
		},
		{
			sample: samples.VideoMPEG4,
			format: "mov,mp4,m4a,3gp,3g2,mj2",
			vCodec: "h264",
			aCodec: "aac",
		},
		{
			sample: samples.VideoOGG,
			format: "ogg",
			vCodec: "theora",
			aCodec: "flac",
		},
		{
			sample: samples.VideoWebM,
			format: "matroska,webm",
			vCodec: "vp9",
			aCodec: "vorbis",
		},
		{
			sample: samples.VideoWMV,
			format: "asf",
			vCodec: "wmv2",
			aCodec: "wmav2",
		},
	} {
		t.Run(item.sample, func(t *testing.T) {
			sample := samples.Buffer(item.sample)
			defer sample.Close()

			report, err := Analyze(sample)
			if item.format == "flv" {
				report.Streams = lo.Reverse(report.Streams)
			}
			if item.format == "gif" {
				report.Streams = append(report.Streams, Stream{
					Type:       "audio",
					Codec:      "",
					Duration:   2.1,
					Channels:   2,
					SampleRate: 44100,
				})
			}
			assert.NoError(t, err)
			assert.True(t, report.Duration >= 2)
			assert.True(t, report.Format.Duration >= 2)
			if !lo.Contains([]string{"flv", "matroska,webm"}, item.format) {
				assert.True(t, report.Streams[0].Duration >= 2, report.Streams[0].Duration)
				assert.True(t, report.Streams[1].Duration >= 2, report.Streams[1].Duration)
			}
			assert.Equal(t, &Report{
				Duration: report.Duration,
				Format: Format{
					Name:     item.format,
					Duration: report.Format.Duration,
				},
				Streams: []Stream{
					{
						Type:      "video",
						Codec:     item.vCodec,
						Duration:  report.Streams[0].Duration,
						Width:     800,
						Height:    450,
						FrameRate: FrameRate(lo.Ternary(item.format == "gif", 5, 25)),
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

func TestAnalyzeImage(t *testing.T) {
	for _, item := range []struct {
		sample string
		format string
		codec  string
	}{
		{
			sample: samples.ImageGIF,
			format: "gif",
			codec:  "gif",
		},
		{
			sample: samples.ImageJPEG,
			format: "image2",
			codec:  "mjpeg",
		},
		{
			sample: samples.ImageJPEG2K,
			format: "j2k_pipe",
			codec:  "jpeg2000",
		},
		{
			sample: samples.ImagePNG,
			format: "png_pipe",
			codec:  "png",
		},
		{
			sample: samples.ImageTIFF,
			format: "tiff_pipe",
			codec:  "tiff",
		},
		{
			sample: samples.ImageWebP,
			format: "webp_pipe",
			codec:  "webp",
		},
	} {
		t.Run(item.sample, func(t *testing.T) {
			sample := samples.Buffer(item.sample)
			defer sample.Close()

			report, err := Analyze(sample)
			assert.NoError(t, err)
			assert.Equal(t, &Report{
				Duration: report.Duration,
				Format: Format{
					Name:     item.format,
					Duration: report.Format.Duration,
				},
				Streams: []Stream{
					{
						Type:      "video",
						Codec:     item.codec,
						Duration:  report.Streams[0].Duration,
						Width:     800,
						Height:    533,
						FrameRate: report.Streams[0].FrameRate,
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
