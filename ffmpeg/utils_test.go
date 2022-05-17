package ffmpeg

import (
	"log"
	"os"
)

func init() {
	WarningsLogger = log.Default()
}

func loadSample(ext string) *os.File {
	f, err := os.Open("../samples/" + ext)
	if err != nil {
		panic(err)
	}

	return f
}
