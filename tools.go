package mediakit

import (
	"context"
	"io"
	"math"
	"os"
	"time"

	"github.com/256dpi/xo"

	"github.com/256dpi/mediakit/chromium"
	"github.com/256dpi/mediakit/ffmpeg"
	"github.com/256dpi/mediakit/vips"
)

// ErrMissingStream is returned if a required audio/video stream is missing.
var ErrMissingStream = xo.BF("missing stream")

// Progress describes a progress update receiver.
type Progress struct {
	Rate time.Duration
	Func func(float64)
}

// ConvertImage will convert an image using a preset and sizer.
func ConvertImage(ctx context.Context, input, output *os.File, preset vips.Preset, sizer Sizer) error {
	// analyze input
	report, err := vips.Analyze(ctx, input)
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
		Width:  report.Width,
		Height: report.Height,
	})

	// prepare options
	opts := vips.ConvertOptions{
		Preset: preset,
		Width:  size.Width,
		Height: size.Height,
	}

	// convert image
	err = vips.Convert(ctx, input, output, opts)
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
func ConvertAudio(ctx context.Context, input, output *os.File, preset ffmpeg.Preset, maxSampleRate int, progress *Progress) error {
	// analyze input
	report, err := ffmpeg.Analyze(ctx, input)
	if err != nil {
		return xo.W(err)
	}

	// check audio stream
	if !report.Has("audio") {
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
		opts.ProgressFunc = func(p ffmpeg.Progress) {
			progress.Func(math.Min(p.Duration/report.Duration, 1))
		}
		opts.ProgressRate = progress.Rate
	}

	// convert audio
	err = ffmpeg.Convert(ctx, input, output, opts)
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
func ConvertVideo(ctx context.Context, input, output *os.File, preset ffmpeg.Preset, sizer Sizer, maxFrameRate float64, maxSampleRate int, progress *Progress) error {
	// analyze input
	report, err := ffmpeg.Analyze(ctx, input)
	if err != nil {
		return xo.W(err)
	}

	// check video stream
	if !report.Has("video") {
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
		Width:  width,
		Height: height,
	})

	// get frame rate
	frameRate := report.FrameRate()
	if frameRate > maxFrameRate {
		frameRate = maxFrameRate
	}

	// get sample rate
	sampleRate := report.SampleRate()
	if sampleRate > maxSampleRate {
		sampleRate = maxSampleRate
	}

	// prepare options
	opts := ffmpeg.ConvertOptions{
		Preset:     preset,
		Width:      size.Width,
		Height:     size.Height,
		FrameRate:  frameRate,
		SampleRate: sampleRate,
	}

	// set progress
	if progress != nil {
		opts.ProgressFunc = func(p ffmpeg.Progress) {
			progress.Func(math.Min(p.Duration/report.Duration, 1))
		}
		opts.ProgressRate = progress.Rate
	}

	// convert video
	err = ffmpeg.Convert(ctx, input, output, opts)
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
func ExtractImage(ctx context.Context, input, temp, output *os.File, position float64, preset vips.Preset, sizer Sizer) error {
	// analyze input
	report, err := ffmpeg.Analyze(ctx, input)
	if err != nil {
		return xo.W(err)
	}

	// check video stream
	if !report.Has("video") {
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
	err = ffmpeg.Convert(ctx, input, temp, opts)
	if err != nil {
		return xo.W(err)
	}

	// convert image
	err = ConvertImage(ctx, temp, output, preset, sizer)
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

// CaptureScreenshot will capture a screenshot using a URL and options.
func CaptureScreenshot(ctx context.Context, url string, output *os.File, opts chromium.Options) error {
	// capture screenshot
	buf, err := chromium.Screenshot(ctx, url, opts)
	if err != nil {
		return err
	}

	// copy image
	_, err = output.Write(buf)
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
