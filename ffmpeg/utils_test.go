package ffmpeg

import (
	"log"
	"os"
)

func init() {
	WarningsLogger = log.Default()
}

func loadSample(name string) *os.File {
	f, err := os.Open("../samples/" + name)
	if err != nil {
		panic(err)
	}

	return f
}
