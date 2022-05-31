package mediakit

import (
	"io"
	"net/http"
	"testing"

	"github.com/gabriel-vasile/mimetype"
	"github.com/stretchr/testify/assert"
)

func TestDetect(t *testing.T) {
	for _, item := range []struct {
		sample string
		typ    string
	}{
		{
			sample: "sample.aac",
			typ:    "audio/aac",
		},
		{
			sample: "sample.aiff",
			typ:    "audio/aiff",
		},
		{
			sample: "sample.avi",
			typ:    "video/avi",
		},
		{
			sample: "sample.flac",
			typ:    "audio/flac",
		},
		{
			sample: "sample.gif",
			typ:    "image/gif",
		},
		{
			sample: "sample.hevc",
			typ:    "application/octet-stream",
		},
		{
			sample: "sample.jpg",
			typ:    "image/jpeg",
		},
		{
			sample: "sample.m4a",
			typ:    "audio/x-m4a",
		},
		{
			sample: "sample.mov",
			typ:    "video/quicktime",
		},
		{
			sample: "sample.mp2",
			typ:    "application/octet-stream",
		},
		{
			sample: "sample.mp3",
			typ:    "audio/mpeg",
		},
		{
			sample: "sample.mp4",
			typ:    "video/mp4",
		},
		{
			sample: "sample.mpeg",
			typ:    "video/mpeg",
		},
		{
			sample: "sample.mpg",
			typ:    "video/mpeg",
		},
		{
			sample: "sample.ogg",
			typ:    "application/ogg",
		},
		{
			sample: "sample.png",
			typ:    "image/png",
		},
		{
			sample: "sample.wav",
			typ:    "audio/wave",
		},
		{
			sample: "sample.webm",
			typ:    "video/webm",
		},
		{
			sample: "sample.wma",
			typ:    "video/x-ms-asf",
		},
		{
			sample: "sample.wmv",
			typ:    "video/x-ms-asf",
		},
	} {
		t.Run(item.sample, func(t *testing.T) {
			sample := loadSample(item.sample)
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
	sample := loadSample("sample.mp3")
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
	sample := loadSample("sample.mp3")
	defer sample.Close()

	buf, err := io.ReadAll(sample)
	assert.NoError(b, err)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		mimetype.Detect(buf)
	}
}
