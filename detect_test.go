package mediakit

import (
	"io"
	"net/http"
	"testing"

	"github.com/gabriel-vasile/mimetype"
	"github.com/stretchr/testify/assert"

	"github.com/256dpi/mediakit/samples"
)

func TestDetect(t *testing.T) {
	for _, item := range []struct {
		sample string
		typ    string
	}{
		// image
		{
			sample: samples.ImageGIF,
			typ:    "image/gif",
		},
		{
			sample: samples.ImageHEIF,
			typ:    "image/heic",
		},
		{
			sample: samples.ImageJPEG,
			typ:    "image/jpeg",
		},
		{
			sample: samples.ImageJPEG2K,
			typ:    "image/jp2",
		},
		{
			sample: samples.ImagePDF,
			typ:    "application/pdf",
		},
		{
			sample: samples.ImagePNG,
			typ:    "image/png",
		},
		{
			sample: samples.ImageTIFF,
			typ:    "image/tiff",
		},
		{
			sample: samples.ImageWebP,
			typ:    "image/webp",
		},
		// audio
		{
			sample: samples.AudioAAC,
			typ:    "audio/aac",
		},
		{
			sample: samples.AudioAIFF,
			typ:    "audio/aiff",
		},
		{
			sample: samples.AudioFLAC,
			typ:    "audio/flac",
		},
		{
			sample: samples.AudioMPEG2,
			typ:    "application/octet-stream",
		},
		{
			sample: samples.AudioMPEG3,
			typ:    "audio/mpeg",
		},
		{
			sample: samples.AudioMPEG4,
			typ:    "audio/x-m4a",
		},
		{
			sample: samples.AudioOGG,
			typ:    "application/ogg",
		},
		{
			sample: samples.AudioWAV,
			typ:    "audio/wave",
		},
		{
			sample: samples.AudioWMA,
			typ:    "video/x-ms-asf",
		},
		// video
		{
			sample: samples.VideoAVI,
			typ:    "video/avi",
		},
		{
			sample: samples.VideoFLV,
			typ:    "video/x-flv",
		},
		{
			sample: samples.VideoGIF,
			typ:    "image/gif",
		},
		{
			sample: samples.VideoMKV,
			typ:    "video/webm",
		},
		{
			sample: samples.VideoMOV,
			typ:    "video/quicktime",
		},
		{
			sample: samples.VideoMPEG,
			typ:    "video/mpeg",
		},
		{
			sample: samples.VideoMPEG2,
			typ:    "video/mpeg",
		},
		{
			sample: samples.VideoMPEG4,
			typ:    "video/mp4",
		},
		{
			sample: samples.VideoWebM,
			typ:    "video/webm",
		},
		{
			sample: samples.VideoWMV,
			typ:    "video/x-ms-asf",
		},
	} {
		t.Run(item.sample, func(t *testing.T) {
			sample := samples.Load(item.sample)
			defer sample.Close()

			buf := make([]byte, DetectBytes)
			n, err := io.ReadFull(sample, buf)
			assert.NoError(t, err)
			buf = buf[:n]

			contentType := Detect(buf)
			assert.Equal(t, item.typ, contentType)
		})
	}
}

func BenchmarkHTTPDetect(b *testing.B) {
	sample := samples.Load(samples.AudioMPEG3)
	defer sample.Close()

	buf, err := io.ReadAll(sample)
	assert.NoError(b, err)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		http.DetectContentType(buf)
	}
}

func BenchmarkMimeTypeDetect(b *testing.B) {
	sample := samples.Load(samples.AudioMPEG3)
	defer sample.Close()

	buf, err := io.ReadAll(sample)
	assert.NoError(b, err)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		mimetype.Detect(buf)
	}
}
