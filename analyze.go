package mediakit

import (
	"context"
	"io"
	"os"

	"github.com/256dpi/xo"
	"github.com/samber/lo"

	"github.com/256dpi/mediakit/ffmpeg"
	"github.com/256dpi/mediakit/vips"
)

// Report describes a file analysis.
type Report struct {
	// generic
	MediaType  string `json:"mediaType"`
	FileFormat string `json:"fileFormat"`

	// image/video
	Width  int `json:"width"`
	Height int `json:"height"`

	// audio/video
	Streams  []string `json:"streams"`
	Codecs   []string `json:"codecs"`
	Duration float64  `json:"duration"`

	// audio
	Channels   int `json:"channels"`
	SampleRate int `json:"sampleRate"`

	// video
	FrameRate float64 `json:"frameRate"`
}

// Analyze will analyze the provided file and return a report.
func Analyze(ctx context.Context, input *os.File) (*Report, error) {
	// detect media stream
	mediaType, _, err := DetectStream(input, false)
	if err != nil {
		return nil, xo.W(err)
	}

	// rewind input
	_, err = input.Seek(0, io.SeekStart)
	if err != nil {
		return nil, xo.W(err)
	}

	// analyze video and audio
	if lo.Contains(AudioTypes(), mediaType) || lo.Contains(VideoTypes(), mediaType) || lo.Contains(ContainerTypes(), mediaType) {
		rep, err := ffmpeg.Analyze(ctx, input)
		if err != nil {
			return nil, xo.W(err)
		}

		// get size
		width, height := rep.Size()

		// get codes and channels
		var streams []string
		var codecs []string
		var channels int
		for _, stream := range rep.Streams {
			if stream.Type != "data" {
				streams = append(streams, stream.Type)
				codecs = append(codecs, stream.Codec)
				if stream.Channels > channels {
					channels = stream.Channels
				}
			}
		}

		return &Report{
			MediaType:  mediaType,
			FileFormat: rep.Format.Name,
			Width:      width,
			Height:     height,
			Streams:    streams,
			Codecs:     codecs,
			Duration:   rep.Duration,
			Channels:   channels,
			SampleRate: rep.SampleRate(),
			FrameRate:  rep.FrameRate(),
		}, nil
	} else if lo.Contains(ImageTypes(), mediaType) {
		rep, err := vips.Analyze(ctx, input)
		if err != nil {
			return nil, xo.W(err)
		}

		return &Report{
			MediaType:  mediaType,
			FileFormat: rep.Format,
			Width:      rep.Width,
			Height:     rep.Height,
		}, nil
	}

	return &Report{
		MediaType: mediaType,
	}, nil
}
