package ffmpeg

import (
	"io"
	"log"
	"os"
	"testing"
)

func init() {
	WarningsLogger = log.Default()
}

func tempFile(t *testing.T) *os.File {
	f, err := os.CreateTemp(t.TempDir(), "ffmpeg-")
	if err != nil {
		panic(err)
	}
	return f
}

func rewind(f *os.File) {
	_, err := f.Seek(0, io.SeekStart)
	if err != nil {
		panic(err)
	}
}
