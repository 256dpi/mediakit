package main

import (
	"flag"
	"os"
	"path/filepath"
	"time"

	"github.com/256dpi/mediakit"
	"github.com/256dpi/mediakit/ffmpeg"
	"github.com/256dpi/mediakit/vips"
	"github.com/256dpi/xo"
	"github.com/kr/pretty"
)

var mode = flag.String("mode", "video", "")

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
		err = mediakit.ConvertImage(nil, input, output, ffmpeg.ImageJPEG, mediakit.MaxWidth(300))
	case "audio":
		err = mediakit.ConvertAudio(nil, input, output, ffmpeg.AudioMP3VBRStandard, 48000, &mediakit.Progress{
			Rate: time.Second,
			Func: func(progress float64) {
				pretty.Println(progress)
			},
		})
	case "video":
		err = mediakit.ConvertVideo(nil, input, output, ffmpeg.VideoMP4H264AACFast, mediakit.MaxWidth(300), 30, 48000, &mediakit.Progress{
			Rate: time.Second,
			Func: func(progress float64) {
				pretty.Println(progress)
			},
		})
	case "extract":
		err = mediakit.ExtractImage(nil, input, temporary, output, 0.25, vips.JPGWeb, mediakit.MaxWidth(300))
	default:
		panic("unknown mode: " + *mode)
	}
	xo.PanicIf(err)

	// remove temporary
	err = os.Remove(outPath + ".tmp")
	xo.PanicIf(err)
}
