package main

import (
	"flag"
	"os"
	"path/filepath"

	"github.com/256dpi/mediakit/chromium"
)

var mode = flag.String("mode", "image", "")
var width = flag.Int64("width", 1920, "")
var height = flag.Int64("height", 1080, "")
var scale = flag.Float64("scale", 2, "")
var full = flag.Bool("full", false, "")
var pedantic = flag.Bool("pedantic", false, "")

func main() {
	// parse flags
	flag.Parse()
	if flag.NArg() != 2 {
		panic("usage: mk-capture <url> <output>")
	}

	// get URL and path
	inURL := flag.Arg(0)
	outPath, err := filepath.Abs(flag.Arg(1))
	if err != nil {
		panic(err)
	}

	// create output
	output, err := os.Create(outPath)
	if err != nil {
		panic(err)
	}
	defer output.Close()

	// convert input
	var buf []byte
	switch *mode {
	case "image":
		buf, err = chromium.CaptureScreenshot(nil, inURL, chromium.Options{
			Width:    *width,
			Height:   *height,
			Scale:    *scale,
			Full:     *full,
			Pedantic: *pedantic,
		})
	default:
		panic("unknown mode: " + *mode)
	}
	if err != nil {
		panic(err)
	}

	// copy output
	_, err = output.Write(buf)
	if err != nil {
		panic(err)
	}
}
