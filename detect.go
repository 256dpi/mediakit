package mediakit

import (
	"bytes"
	"io"
	"net/http"

	"github.com/gabriel-vasile/mimetype"
)

// ImageTypes is a canonical list of well-known and modern image formats
// supported by mediakit.
func ImageTypes() []string {
	return []string{
		"image/gif",
		"image/heic",
		"image/heif",
		"image/jpeg",
		"image/jp2",
		"application/pdf",
		"image/png",
		"image/tiff",
		"image/webp",
	}
}

// AudioTypes is a canonical list of well-known and modern audio formats
// supported by mediakit.
func AudioTypes() []string {
	return []string{
		"audio/aac",
		"audio/aiff",
		"audio/flac",
		"audio/mpeg",
		"audio/x-m4a",
		"audio/wave",
	}
}

// VideoTypes is a canonical list of well-known and modern video formats
// supported by mediakit.
func VideoTypes() []string {
	return []string{
		"video/avi",
		"video/x-flv",
		"image/gif",
		"video/webm",
		"video/quicktime",
		"video/mpeg",
		"video/mp4",
		"video/webm",
	}
}

// ContainerTypes is a canonical list of additional container formats supported
// by mediakit that may contain audio, video or both.
func ContainerTypes() []string {
	return []string{
		"application/ogg",
		"video/x-ms-asf",
	}
}

// DetectBytes defines the maximum number of bytes used by Detect.
const DetectBytes = 3072

// Detect will attempt to detect a media type from the specified buffer.
// It delegates to http.DetectContentType and mimetype.Detect which together
// should detect a faire amount of media types and falls back to
// "application/octet-stream" if undetected.
func Detect(buf []byte) string {
	// use built-in detector
	typ := http.DetectContentType(buf)

	// use mimetype if not found
	if typ == "application/octet-stream" {
		typ = mimetype.Detect(buf).String()
	}

	return typ
}

// DetectStream will attempt to detect a media from the provided reader using
// Detect. It will read up to DetectBytes from the reader and return a new
// reader that will read from the read bytes and the remaining stream.
func DetectStream(stream io.Reader) (string, io.Reader, error) {
	// read from stream
	buf := make([]byte, DetectBytes)
	n, err := io.ReadFull(stream, buf)
	if err != nil {
		return "", nil, err
	}

	// resize buffer
	buf = buf[:n]

	// detect media type
	typ := Detect(buf)

	// assemble reader
	stream = io.MultiReader(bytes.NewReader(buf), stream)

	return typ, stream, nil
}
