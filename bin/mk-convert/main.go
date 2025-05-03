package main

import (
	"flag"
	"os"
	"path/filepath"
	"time"

	"github.com/256dpi/xo"
	"github.com/kr/pretty"

	"github.com/256dpi/mediakit"
	"github.com/256dpi/mediakit/ffmpeg"
	"github.com/256dpi/mediakit/vips"
)

var mode = flag.String("mode", "video", "")
var preset = flag.Int("preset", 0, "")
var width = flag.Int("width", 300, "")

func main() {
	// parse flags
	flag.Parse()
	if flag.NArg() != 2 {
		panic("usage: mk-convert <input> <output>")
	}

	// get paths
	inPath, err := filepath.Abs(flag.Arg(0))
	xo.PanicIf(err)
	outPath, err := filepath.Abs(flag.Arg(1))
	xo.PanicIf(err)

	// open input
	input, err := os.Open(inPath)
	xo.PanicIf(err)
	defer input.Close()

	// create temporary file
	temporary, err := os.Create(outPath + ".tmp")
	xo.PanicIf(err)
	defer temporary.Close()

	// create output
	output, err := os.Create(outPath)
	xo.PanicIf(err)
	defer output.Close()

	// convert input
	switch *mode {
	case "image":
		if *preset == 0 {
			*preset = int(vips.JPGWeb)
		}
		err = mediakit.ConvertImage(nil, input, output, vips.Preset(*preset), mediakit.MaxWidth(*width))
	case "audio":
		if *preset == 0 {
			*preset = int(ffmpeg.AudioMP3VBRStandard)
		}
		err = mediakit.ConvertAudio(nil, input, output, ffmpeg.Preset(*preset), 48000, &mediakit.Progress{
			Rate: time.Second,
			Func: func(progress float64) {
				pretty.Println(progress)
			},
		})
	case "video":
		if *preset == 0 {
			*preset = int(ffmpeg.VideoMP4H264AACFast)
		}
		err = mediakit.ConvertVideo(nil, input, output, ffmpeg.Preset(*preset), mediakit.MaxWidth(*width), 30, 48000, &mediakit.Progress{
			Rate: time.Second,
			Func: func(progress float64) {
				pretty.Println(progress)
			},
		})
	case "extract":
		if *preset == 0 {
			*preset = int(vips.JPGWeb)
		}
		err = mediakit.ExtractImage(nil, input, temporary, output, 0.25, vips.Preset(*preset), mediakit.MaxWidth(*width))
	default:
		panic("unknown mode: " + *mode)
	}
	xo.PanicIf(err)

	// remove temporary
	err = os.Remove(outPath + ".tmp")
	xo.PanicIf(err)
}
