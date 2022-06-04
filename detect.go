package mediakit

import (
	"net/http"

	"github.com/gabriel-vasile/mimetype"
)

// DetectBytes defines the maximum number of bytes used by Detect.
const DetectBytes = 3072

// Detect will attempt to detect a content type from the specified buffer.
// It delegates to http.DetectContentType and mimetype.Detect which together
// should detect a faire amount of content types and falls back to
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
