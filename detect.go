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
