package ffmpeg

import (
	"io"
	"os"
)

func loadSample(ext string) io.ReadCloser {
	f, err := os.Open("../samples/sample." + ext)
	if err != nil {
		panic(err)
	}

	return f
}
