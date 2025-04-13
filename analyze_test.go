package mediakit

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/256dpi/mediakit/samples"
)

func TestAnalyze(t *testing.T) {
	for _, item := range []struct {
		sample string
		report Report
	}{
		// image
		{
			sample: samples.ImageGIF,
			report: Report{
				MediaType:  "image/gif",
				FileFormat: "gif",
				Width:      800,
				Height:     533,
				Streams:    []string{"video"},
				Codecs:     []string{"gif"},
				Duration:   0.1,
				FrameRate:  100,
			},
		},
		{
			sample: samples.ImageHEIF,
			report: Report{
				MediaType:  "image/heic",
				FileFormat: "heif",
				Width:      800,
				Height:     533,
			},
		},
		{
			sample: samples.ImageJPEG,
			report: Report{
				MediaType:  "image/jpeg",
				FileFormat: "jpeg",
				Width:      800,
				Height:     533,
			},
		},
		{
			sample: samples.ImageJPEG2K,
			report: Report{
				MediaType:  "image/jp2",
				FileFormat: "jp2k",
				Width:      800,
				Height:     533,
			},
		},
		{
			sample: samples.ImagePDF,
			report: Report{
				MediaType:  "application/pdf",
				FileFormat: "pdf",
				Width:      800,
				Height:     533,
			},
		},
		{
			sample: samples.ImagePNG,
			report: Report{
				MediaType:  "image/png",
				FileFormat: "png",
				Width:      800,
				Height:     533,
			},
		},
		{
			sample: samples.ImageTIFF,
			report: Report{
				MediaType:  "image/tiff",
				FileFormat: "tiff",
				Width:      800,
				Height:     533,
			},
		},
		{
			sample: samples.ImageWebP,
			report: Report{
				MediaType:  "image/webp",
				FileFormat: "webp",
				Width:      800,
				Height:     533,
			},
		},
		// audio
		{
			sample: samples.AudioAAC,
			report: Report{
				MediaType:  "audio/aac",
				FileFormat: "aac",
				Streams:    []string{"audio"},
				Codecs:     []string{"aac"},
				Duration:   2.127203,
				Channels:   2,
				SampleRate: 44100,
			},
		},
		{
			sample: samples.AudioAIFF,
			report: Report{
				MediaType:  "audio/aiff",
				FileFormat: "aiff",
				Streams:    []string{"audio"},
				Codecs:     []string{"pcm_s16be"},
				Duration:   2.043356,
				Channels:   2,
				SampleRate: 44100,
			},
		},
		{
			sample: samples.AudioFLAC,
			report: Report{
				MediaType:  "audio/flac",
				FileFormat: "flac",
				Streams:    []string{"audio"},
				Codecs:     []string{"flac"},
				Duration:   2.115918,
				Channels:   2,
				SampleRate: 44100,
			},
		},
		{
			sample: samples.AudioMPEG3,
			report: Report{
				MediaType:  "audio/mpeg",
				FileFormat: "mp3",
				Streams:    []string{"audio"},
				Codecs:     []string{"mp3"},
				Duration:   2.123813,
				Channels:   2,
				SampleRate: 44100,
			},
		},
		{
			sample: samples.AudioMPEG4,
			report: Report{
				MediaType:  "audio/x-m4a",
				FileFormat: "mov,mp4,m4a,3gp,3g2,mj2",
				Streams:    []string{"audio"},
				Codecs:     []string{"aac"},
				Duration:   2.115918,
				Channels:   2,
				SampleRate: 44100,
			},
		},
		{
			sample: samples.AudioOGG,
			report: Report{
				MediaType:  "audio/ogg",
				FileFormat: "ogg",
				Streams:    []string{"audio"},
				Codecs:     []string{"vorbis"},
				Duration:   2.115918,
				Channels:   2,
				SampleRate: 44100,
			},
		},
		{
			sample: samples.AudioWAV,
			report: Report{
				MediaType:  "audio/wav",
				FileFormat: "wav",
				Streams:    []string{"audio"},
				Codecs:     []string{"pcm_s24le"},
				Duration:   2.043356,
				Channels:   2,
				SampleRate: 44100,
			},
		},
		{
			sample: samples.AudioWMA,
			report: Report{
				MediaType:  "video/x-ms-asf",
				FileFormat: "asf",
				Streams:    []string{"audio"},
				Codecs:     []string{"wmav2"},
				Duration:   2.135,
				Channels:   2,
				SampleRate: 44100,
			},
		},
		// video
		{
			sample: samples.VideoAVI,
			report: Report{
				MediaType:  "video/x-msvideo",
				FileFormat: "avi",
				Width:      800,
				Height:     450,
				Streams:    []string{"video", "audio"},
				Codecs:     []string{"h264", "aac"},
				Duration:   2.136236,
				Channels:   2,
				SampleRate: 44100,
				FrameRate:  25,
			},
		},
		{
			sample: samples.VideoFLV,
			report: Report{
				MediaType:  "video/x-flv",
				FileFormat: "flv",
				Width:      800,
				Height:     450,
				Streams:    []string{"audio", "video"},
				Codecs:     []string{"mp3", "flv1"},
				Duration:   2.069,
				Channels:   2,
				SampleRate: 44100,
				FrameRate:  25,
			},
		},
		{
			sample: samples.VideoGIF,
			report: Report{
				MediaType:  "image/gif",
				FileFormat: "gif",
				Width:      800,
				Height:     450,
				Streams:    []string{"video"},
				Codecs:     []string{"gif"},
				Duration:   2,
				Channels:   0,
				SampleRate: 0,
				FrameRate:  5,
			},
		},
		{
			sample: samples.VideoMKV,
			report: Report{
				MediaType:  "video/x-matroska",
				FileFormat: "matroska,webm",
				Width:      800,
				Height:     450,
				Streams:    []string{"video", "audio"},
				Codecs:     []string{"hevc", "ac3"},
				Duration:   2.055,
				Channels:   2,
				SampleRate: 44100,
				FrameRate:  25,
			},
		},
		{
			sample: samples.VideoMOV,
			report: Report{
				MediaType:  "video/quicktime",
				FileFormat: "mov,mp4,m4a,3gp,3g2,mj2",
				Width:      800,
				Height:     450,
				Streams:    []string{"video", "audio"},
				Codecs:     []string{"h264", "aac"},
				Duration:   2.042993,
				Channels:   2,
				SampleRate: 44100,
				FrameRate:  25,
			},
		},
		{
			sample: samples.VideoMPEG,
			report: Report{
				MediaType:  "video/mpeg",
				FileFormat: "mpeg",
				Width:      800,
				Height:     450,
				Streams:    []string{"video", "audio"},
				Codecs:     []string{"mpeg1video", "mp2"},
				Duration:   2.063678,
				Channels:   2,
				SampleRate: 44100,
				FrameRate:  50,
			},
		},
		{
			sample: samples.VideoMPEG2,
			report: Report{
				MediaType:  "video/mpeg",
				FileFormat: "mpeg",
				Width:      800,
				Height:     450,
				Streams:    []string{"video", "audio"},
				Codecs:     []string{"mpeg2video", "mp2"},
				Duration:   2.063678,
				Channels:   2,
				SampleRate: 44100,
				FrameRate:  25,
			},
		},
		{
			sample: samples.VideoMPEG4,
			report: Report{
				MediaType:  "video/mp4",
				FileFormat: "mov,mp4,m4a,3gp,3g2,mj2",
				Width:      800,
				Height:     450,
				Streams:    []string{"video", "audio"},
				Codecs:     []string{"h264", "aac"},
				Duration:   2.04,
				Channels:   2,
				SampleRate: 44100,
				FrameRate:  25,
			},
		},
		{
			sample: samples.VideoMPEG4R,
			report: Report{
				MediaType:  "video/mp4",
				FileFormat: "mov,mp4,m4a,3gp,3g2,mj2",
				Width:      450,
				Height:     800,
				Streams:    []string{"video", "audio"},
				Codecs:     []string{"h264", "aac"},
				Duration:   2.04,
				Channels:   2,
				SampleRate: 44100,
				FrameRate:  25,
			},
		},
		{
			sample: samples.VideoOGG,
			report: Report{
				MediaType:  "video/ogg",
				FileFormat: "ogg",
				Width:      800,
				Height:     450,
				Streams:    []string{"video", "audio"},
				Codecs:     []string{"theora", "flac"},
				Duration:   2.043356,
				Channels:   2,
				SampleRate: 44100,
				FrameRate:  25,
			},
		},
		{
			sample: samples.VideoWebM,
			report: Report{
				MediaType:  "video/webm",
				FileFormat: "matroska,webm",
				Width:      800,
				Height:     450,
				Streams:    []string{"video", "audio"},
				Codecs:     []string{"vp9", "vorbis"},
				Duration:   2.05,
				Channels:   2,
				SampleRate: 44100,
				FrameRate:  25,
			},
		},
		{
			sample: samples.VideoWMV,
			report: Report{
				MediaType:  "video/x-ms-asf",
				FileFormat: "asf",
				Width:      800,
				Height:     450,
				Streams:    []string{"video", "audio"},
				Codecs:     []string{"wmv2", "wmav2"},
				Duration:   2.132,
				Channels:   2,
				SampleRate: 44100,
				FrameRate:  25,
			},
		},
		// document
		{
			sample: samples.DocumentPDF,
			report: Report{
				MediaType:  "application/pdf",
				FileFormat: "pdf",
				Width:      595,
				Height:     842,
			},
		},
	} {
		sample := samples.Buffer(item.sample)
		defer sample.Close()

		report, err := Analyze(nil, sample)
		assert.NoError(t, err)
		assert.Equal(t, &item.report, report, item.sample)
	}
}
