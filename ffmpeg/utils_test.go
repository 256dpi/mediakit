package ffmpeg

import (
	"os"
)

func loadSample(ext string) *os.File {
	f, err := os.Open("../samples/sample." + ext)
	if err != nil {
		panic(err)
	}

	return f
}
