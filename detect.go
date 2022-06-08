package mediakit

import (
	"bytes"
	"io"
	"net/http"
	"strings"

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
		"audio/x-m4a",
		"audio/mpeg",
		"audio/ogg",
		"audio/wav",
	}
}

// VideoTypes is a canonical list of well-known and modern video formats
// supported by mediakit.
func VideoTypes() []string {
	return []string{
		"video/x-msvideo", // avi
		"video/x-flv",
		"image/gif",
		"video/x-matroska",
		"video/quicktime", // mov
		"video/mpeg",
		"video/mp4",
		"video/ogg",
		"video/webm",
	}
}

// ContainerTypes is a canonical list of additional container formats supported
// by mediakit that may contain audio, video or both.
func ContainerTypes() []string {
	return []string{
		"video/x-ms-asf", // wma, wmv
	}
}

// DetectBytes defines the maximum number of bytes used by Detect.
const DetectBytes = 3072

// Detect will attempt to detect a media type from the specified buffer.
// It delegates to mimetype.Detect and http.DetectContentType which together
// should detect a faire amount of media types and falls back to
// "application/octet-stream" if undetected.
func Detect(buf []byte, withParameters bool) string {
	// use built-in detector
	typ := mimetype.Detect(buf).String()

	// use native if not found
	if typ == "application/octet-stream" {
		typ = http.DetectContentType(buf)
	}

	// remove parameters if not requested and present
	if !withParameters && strings.Contains(typ, ";") {
		typ = typ[:strings.Index(typ, ";")]
	}

	return typ
}

// DetectStream will attempt to detect a media from the provided reader using
// Detect. It will read up to DetectBytes from the reader and return a new
// reader that will read from the read bytes and the remaining stream.
func DetectStream(stream io.Reader, withCharset bool) (string, io.Reader, error) {
	// read from stream
	buf := make([]byte, DetectBytes)
	n, err := io.ReadFull(stream, buf)
	if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
		return "", nil, err
	}

	// resize buffer
	buf = buf[:n]

	// detect media type
	typ := Detect(buf, withCharset)

	// assemble reader
	stream = io.MultiReader(bytes.NewReader(buf), stream)

	return typ, stream, nil
}
