package main

import (
	"flag"
	"os"
	"path/filepath"
	"time"

	"github.com/256dpi/mediakit"
	"github.com/256dpi/mediakit/ffmpeg"
	"github.com/kr/pretty"
)

var mode = flag.String("mode", "video", "")

func main() {
	// parse flags
	flag.Parse()

	// get path
	path, err := filepath.Abs(flag.Arg(0))
	if err != nil {
		panic(err)
	}

	// open input
	input, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer input.Close()

	// open output
	output, err := os.Create(path + ".out")
	if err != nil {
		panic(err)
	}
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
	default:
		panic("unknown mode: " + *mode)
	}
	if err != nil {
		panic(err)
	}
}
