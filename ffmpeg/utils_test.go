package ffmpeg

import (
	"log"
	"os"
	"runtime"
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

var isDarwin = runtime.GOOS == "darwin"

func osInt(macos, linux int64) int64 {
	if isDarwin {
		return macos
	}
	return linux
}

func osFloat(macos, linux float64) float64 {
	if isDarwin {
		return macos
	}
	return linux
}
