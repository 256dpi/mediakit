package mediakit

import (
	"net/http"

	"github.com/gabriel-vasile/mimetype"
)

// Detect will attempt to detect a content type from the specified buffer.
// It delegates to http.DetectContentType and mimetype.Detect which together
// should detect a faire amount of content types and falls back to
// "application/octet-stream" if undetected.
func Detect(buf []byte) string {
	// use built-in detector first and if not found use mimetype
	typ := http.DetectContentType(buf)
	if typ == "application/octet-stream" {
		typ = mimetype.Detect(buf).String()
	}

	// return if not found
	if typ == "application/octet-stream" {
		return typ
	}

	return typ
}
