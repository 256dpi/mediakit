package mediakit

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/256dpi/mediakit/samples"
)

func TestTypes(t *testing.T) {
	for _, typ := range AnimationTypes() {
		assert.Contains(t, ImageTypes(), typ)
	}
}

func TestDetect(t *testing.T) {
	assert.Equal(t, "text/plain", Detect([]byte{}, true))
	assert.Equal(t, "application/octet-stream", Detect([]byte("\x00\x01"), true))
	assert.Equal(t, "text/plain; charset=utf-8", Detect([]byte("Hello world!"), true))
	assert.Equal(t, "text/plain", Detect([]byte("Hello world!"), false))
	assert.Equal(t, "image/jpeg", Detect(samples.Read(samples.ImageJPEG), true))
}

func TestDetectImages(t *testing.T) {
	for _, sample := range samples.Images() {
		t.Run(sample, func(t *testing.T) {
			typ := Detect(samples.Read(sample), false)
			assert.Contains(t, ImageTypes(), typ)
		})
	}
}

func TestDetectAudio(t *testing.T) {
	list := append(AudioTypes(), ContainerTypes()...)

	for _, sample := range samples.Audio() {
		t.Run(sample, func(t *testing.T) {
			typ := Detect(samples.Read(sample), false)
			assert.Contains(t, list, typ)
		})
	}
}

func TestDetectVideo(t *testing.T) {
	list := append(VideoTypes(), ContainerTypes()...)

	for _, sample := range samples.Video() {
		t.Run(sample, func(t *testing.T) {
			typ := Detect(samples.Read(sample), false)
			assert.Contains(t, list, typ)
		})
	}
}

func TestDetectStream(t *testing.T) {
	/* empty */

	typ, stream, err := DetectStream(bytes.NewReader([]byte{}), false)
	assert.NoError(t, err)
	assert.Equal(t, "text/plain", typ)
	assert.NotNil(t, stream)

	buf, err := io.ReadAll(stream)
	assert.NoError(t, err)
	assert.Equal(t, []byte{}, buf)

	/* bytes */

	typ, stream, err = DetectStream(bytes.NewReader([]byte("\x01\x02")), false)
	assert.NoError(t, err)
	assert.Equal(t, "application/octet-stream", typ)
	assert.NotNil(t, stream)

	buf, err = io.ReadAll(stream)
	assert.NoError(t, err)
	assert.Equal(t, []byte("\x01\x02"), buf)

	/* audio */

	sample := samples.Load(samples.AudioMPEG3)
	defer sample.Close()

	typ, stream, err = DetectStream(sample, false)
	assert.NoError(t, err)
	assert.Equal(t, "audio/mpeg", typ)
	assert.NotNil(t, stream)

	buf, err = io.ReadAll(stream)
	assert.NoError(t, err)
	assert.Equal(t, samples.Read(samples.AudioMPEG3), buf)
}

func BenchmarkDetect(b *testing.B) {
	sample := samples.Read(samples.AudioMPEG3)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		Detect(sample, false)
	}
}
