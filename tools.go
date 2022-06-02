package mediakit

import (
	"io"
	"math"
	"os"

	"github.com/256dpi/xo"

	"github.com/256dpi/mediakit/ffmpeg"
	"github.com/256dpi/mediakit/vips"
)

// ErrMissingStream is returned if a required audio/video stream is missing.
var ErrMissingStream = xo.BF("missing stream")

// ConvertImage will convert an image using a preset and sizer.
func ConvertImage(input, output *os.File, preset vips.Preset, sizer Sizer) error {
	// analyze input
	report, err := vips.Analyze(input)
	if err != nil {
		return xo.W(err)
	}

	// rewind input
	_, err = input.Seek(0, io.SeekStart)
	if err != nil {
		return xo.W(err)
	}

	// apply sizer
	size := sizer(Size{
		W: float64(report.Width),
		H: float64(report.Height),
	})

	// prepare options
	opts := vips.ConvertOptions{
		Preset: preset,
		Width:  int(math.Round(size.W)),
		Height: int(math.Round(size.H)),
	}

	// convert image
	err = vips.Convert(input, output, opts)
	if err != nil {
		return xo.W(err)
	}

	// sync and rewind file
	err = syncAndRewind(output)
	if err != nil {
		return err
	}

	return nil
}

// ConvertAudio will convert audio using a preset.
func ConvertAudio(input, output *os.File, preset ffmpeg.Preset, maxSampleRate int, progress func(float64)) error {
	// analyze input
	report, err := ffmpeg.Analyze(input)
	if err != nil {
		return xo.W(err)
	}

	// check audio stream
	var ok bool
	for _, stream := range report.Streams {
		if stream.Type == "audio" {
			ok = true
		}
	}
	if !ok {
		return ErrMissingStream.Wrap()
	}

	// rewind input
	_, err = input.Seek(0, io.SeekStart)
	if err != nil {
		return xo.W(err)
	}

	// get frame rate
	sampleRate := report.SampleRate()
	if sampleRate > maxSampleRate {
		sampleRate = maxSampleRate
	}

	// prepare options
	opts := ffmpeg.ConvertOptions{
		Preset:     preset,
		SampleRate: sampleRate,
	}

	// set progress
	if progress != nil {
		opts.Progress = func(p ffmpeg.Progress) {
			progress(math.Min(p.Duration/report.Duration, 1))
		}
	}

	// convert audio
	err = ffmpeg.Convert(input, output, opts)
	if err != nil {
		return xo.W(err)
	}

	// sync and rewind file
	err = syncAndRewind(output)
	if err != nil {
		return err
	}

	return nil
}

// ConvertVideo will convert video using a preset, sizer and max frame rate.
func ConvertVideo(input, output *os.File, preset ffmpeg.Preset, sizer Sizer, maxFrameRate float64, progress func(float64)) error {
	// analyze input
	report, err := ffmpeg.Analyze(input)
	if err != nil {
		return xo.W(err)
	}

	// check video stream
	var ok bool
	for _, stream := range report.Streams {
		if stream.Type == "video" {
			ok = true
		}
	}
	if !ok {
		return ErrMissingStream.Wrap()
	}

	// rewind input
	_, err = input.Seek(0, io.SeekStart)
	if err != nil {
		return xo.W(err)
	}

	// get size
	width, height := report.Size()

	// apply sizer
	size := sizer(Size{
		W: float64(width),
		H: float64(height),
	})

	// get frame rate
	frameRate := report.FrameRate()
	if frameRate > maxFrameRate {
		frameRate = maxFrameRate
	}

	// prepare options
	opts := ffmpeg.ConvertOptions{
		Preset:    preset,
		Width:     int(math.Round(size.W)),
		Height:    int(math.Round(size.H)),
		FrameRate: frameRate,
	}

	// set progress
	if progress != nil {
		opts.Progress = func(p ffmpeg.Progress) {
			progress(math.Min(p.Duration/report.Duration, 1))
		}
	}

	// convert video
	err = ffmpeg.Convert(input, output, opts)
	if err != nil {
		return xo.W(err)
	}

	// sync and rewind file
	err = syncAndRewind(output)
	if err != nil {
		return err
	}

	return nil
}

// ExtractImage will extract an image using a position, preset and sizer.
func ExtractImage(input, temp, output *os.File, position float64, preset vips.Preset, sizer Sizer) error {
	// analyze input
	report, err := ffmpeg.Analyze(input)
	if err != nil {
		return xo.W(err)
	}

	// check audio stream
	var ok bool
	for _, stream := range report.Streams {
		if stream.Type == "video" {
			ok = true
		}
	}
	if !ok {
		return ErrMissingStream.Wrap()
	}

	// rewind input
	_, err = input.Seek(0, io.SeekStart)
	if err != nil {
		return xo.W(err)
	}

	// prepare options
	opts := ffmpeg.ConvertOptions{
		Preset: ffmpeg.ImagePNG, // lossless
		Start:  report.Duration * position,
	}

	// convert video
	err = ffmpeg.Convert(input, temp, opts)
	if err != nil {
		return xo.W(err)
	}

	// convert image
	err = ConvertImage(temp, output, preset, sizer)
	if err != nil {
		return err
	}

	// sync and rewind file
	err = syncAndRewind(output)
	if err != nil {
		return err
	}

	return nil
}

func syncAndRewind(file *os.File) error {
	// sync file
	err := file.Sync()
	if err != nil {
		return err
	}

	// rewind file
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}

	return nil
}
