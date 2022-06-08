package mediakit

import (
	"io"
	"net/http"
	"testing"

	"github.com/gabriel-vasile/mimetype"
	"github.com/stretchr/testify/assert"

	"github.com/256dpi/mediakit/samples"
)

func TestDetectImages(t *testing.T) {
	for _, sample := range samples.Images() {
		t.Run(sample, func(t *testing.T) {
			typ := Detect(samples.Read(sample))
			assert.Contains(t, ImageTypes(), typ)
		})
	}
}

func TestDetectAudio(t *testing.T) {
	list := append(AudioTypes(), ContainerTypes()...)

	for _, sample := range samples.Audio() {
		t.Run(sample, func(t *testing.T) {
			typ := Detect(samples.Read(sample))
			assert.Contains(t, list, typ)
		})
	}
}

func TestDetectVideo(t *testing.T) {
	list := append(VideoTypes(), ContainerTypes()...)

	for _, sample := range samples.Video() {
		t.Run(sample, func(t *testing.T) {
			typ := Detect(samples.Read(sample))
			assert.Contains(t, list, typ)
		})
	}
}

func TestDetectStream(t *testing.T) {
	sample := samples.Load(samples.AudioMPEG3)
	defer sample.Close()

	typ, stream, err := DetectStream(sample)
	assert.NoError(t, err)
	assert.Equal(t, "audio/mpeg", typ)
	assert.NotNil(t, stream)

	buf, err := io.ReadAll(stream)
	assert.NoError(t, err)
	assert.Equal(t, samples.Read(samples.AudioMPEG3), buf)
}

func BenchmarkHTTPDetect(b *testing.B) {
	sample := samples.Read(samples.AudioMPEG3)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		http.DetectContentType(sample)
	}
}

func BenchmarkMimeTypeDetect(b *testing.B) {
	sample := samples.Read(samples.AudioMPEG3)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		mimetype.Detect(sample)
	}
}
