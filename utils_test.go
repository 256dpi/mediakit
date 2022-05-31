package mediakit

import (
	"log"
	"os"

	"github.com/256dpi/mediakit/ffmpeg"
)

func init() {
	ffmpeg.WarningsLogger = log.Default()
}

func loadSample(name string) *os.File {
	f, err := os.Open("./samples/" + name)
	if err != nil {
		panic(err)
	}

	return f
}
