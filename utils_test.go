package mediakit

import (
	"log"
	"math"
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

func round(f []float64) []float64 {
	for i, n := range f {
		f[i] = math.Round(n*10) / 10
	}
	return f
}
