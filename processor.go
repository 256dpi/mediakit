package mediakit

import (
	"io"
	"math"
	"os"
	"path/filepath"

	"github.com/256dpi/xo"
	"github.com/google/uuid"
	"github.com/samber/lo"

	"github.com/256dpi/mediakit/ffmpeg"
	"github.com/256dpi/mediakit/vips"
)

// Errors returned by the Processor.
var (
	ErrUnsupportedFormat = xo.BF("unsupported format")
	ErrUnsupportedStream = xo.BF("unsupported stream")
	ErrUnsupportedCodec  = xo.BF("unsupported codec")
)

// Config defines a Processor configuration.
type Config struct {
	// The directory to used for temporary files.
	Directory string

	// The supported formats and codecs.
	ImageFormats []string
	VideoFormats []string
	AudioFormats []string
	VideoCodecs  []string
	AudioCodecs  []string

	// The used presets.
	ImagePreset vips.Preset   // vips.JPGWeb
	AudioPreset ffmpeg.Preset // ffmpeg.AudioMP3VBRStandard
	VideoPreset ffmpeg.Preset // ffmpeg.VideoMP4H264AACFast

	// The maximum and targeted video frame rate.
	MaxFrameRate    float64 // 30
	TargetFrameRate float64 // 25
}

// Processor provides methods to convert streams and files as well as extract
// images from video streams or files.
type Processor struct {
	config Config
}

// NewProcessor creates and returns a new processor using the specified config.
func NewProcessor(config Config) *Processor {
	// ensure defaults
	if config.ImagePreset == 0 {
		config.ImagePreset = vips.JPGWeb
	}
	if config.AudioPreset == 0 {
		config.AudioPreset = ffmpeg.AudioMP3VBRStandard
	}
	if config.VideoPreset == 0 {
		config.VideoPreset = ffmpeg.VideoMP4H264AACFast
	}
	if config.MaxFrameRate == 0 {
		config.MaxFrameRate = 30
	}
	if config.TargetFrameRate == 0 {
		config.TargetFrameRate = 25
	}

	return &Processor{
		config: config,
	}
}

// ConvertImage will convert an image stream with the specified sizer applied.
func (p *Processor) ConvertImage(input io.Reader, sizer Sizer, fn func(output *os.File) error) error {
	return p.buffer(input, false, fn, func(input, _, output *os.File) error {
		return p.ConvertImageFile(input, output, sizer)
	})
}

// ConvertImageFile will convert an image file with the specified sizer applied.
func (p *Processor) ConvertImageFile(input, output *os.File, sizer Sizer) error {
	// analyze input
	report, err := vips.Analyze(input)
	if err != nil {
		return xo.W(err)
	}

	// check format
	if !lo.Contains(p.config.ImageFormats, report.Format) {
		return ErrUnsupportedFormat.WrapF(report.Format)
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
		Preset: p.config.ImagePreset,
		Width:  int(math.Round(size.W)),
		Height: int(math.Round(size.H)),
	}

	// convert image
	err = vips.Convert(input, output, opts)
	if err != nil {
		return xo.W(err)
	}

	return nil
}

// ConvertAudio will convert an audio stream.
func (p *Processor) ConvertAudio(input io.Reader, progress func(float64), fn func(output *os.File) error) error {
	return p.buffer(input, false, fn, func(input, _, output *os.File) error {
		return p.ConvertAudioFile(input, output, progress)
	})
}

// ConvertAudioFile will convert an audio file.
func (p *Processor) ConvertAudioFile(input, output *os.File, progress func(float64)) error {
	// analyze input
	report, err := ffmpeg.Analyze(input)
	if err != nil {
		return xo.W(err)
	}

	// check format, stream type and stream codec
	if !lo.Contains(p.config.AudioFormats, report.Format.Name) {
		return ErrUnsupportedFormat.WrapF(report.Format.Name)
	}
	for _, stream := range report.Streams {
		if stream.Type != "audio" {
			return ErrUnsupportedStream.WrapF(stream.Type)
		} else if !lo.Contains(p.config.AudioCodecs, stream.Codec) {
			return ErrUnsupportedCodec.WrapF(stream.Codec)
		}
	}

	// rewind input
	_, err = input.Seek(0, io.SeekStart)
	if err != nil {
		return xo.W(err)
	}

	// prepare options
	opts := ffmpeg.ConvertOptions{
		Preset: p.config.AudioPreset,
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

	return nil
}

// ConvertVideo will convert a video stream with the specified sizer applied.
func (p *Processor) ConvertVideo(input io.Reader, sizer Sizer, progress func(float64), fn func(output *os.File) error) error {
	return p.buffer(input, false, fn, func(input, _, output *os.File) error {
		return p.ConvertVideoFile(input, output, sizer, progress)
	})
}

// ConvertVideoFile will convert a video file with the specified sizer applied.
func (p *Processor) ConvertVideoFile(input, output *os.File, sizer Sizer, progress func(float64)) error {
	// analyze input
	report, err := ffmpeg.Analyze(input)
	if err != nil {
		return xo.W(err)
	}

	// check format, stream type and stream codec
	if !lo.Contains(p.config.VideoFormats, report.Format.Name) {
		return ErrUnsupportedFormat.WrapF(report.Format.Name)
	}
	for _, stream := range report.Streams {
		if stream.Type != "video" {
			return ErrUnsupportedStream.WrapF(stream.Type)
		} else if !lo.Contains(p.config.VideoCodecs, stream.Codec) {
			return ErrUnsupportedCodec.WrapF(stream.Codec)
		}
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
	if frameRate > p.config.MaxFrameRate {
		frameRate = p.config.TargetFrameRate
	}

	// prepare options
	opts := ffmpeg.ConvertOptions{
		Preset:    p.config.VideoPreset,
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

	return nil
}

// ExtractImage will extract an image stream from a video stream at the provided
// position with the specified sizer applied.
func (p *Processor) ExtractImage(input io.Reader, pos float64, sizer Sizer, fn func(output *os.File) error) error {
	return p.buffer(input, true, fn, func(input, temp, output *os.File) error {
		return p.ExtractImageFile(input, temp, output, pos, sizer)
	})
}

// ExtractImageFile will extract an image file from a video file at the
// provided position with the specified sizer applied.
func (p *Processor) ExtractImageFile(input, temp, output *os.File, pos float64, sizer Sizer) error {
	// analyze input
	report, err := ffmpeg.Analyze(input)
	if err != nil {
		return xo.W(err)
	}

	// check format, stream type and stream codec
	if !lo.Contains(p.config.VideoFormats, report.Format.Name) {
		return ErrUnsupportedFormat.WrapF(report.Format.Name)
	}
	for _, stream := range report.Streams {
		if stream.Type != "video" {
			return ErrUnsupportedStream.WrapF(stream.Type)
		} else if !lo.Contains(p.config.VideoCodecs, stream.Codec) {
			return ErrUnsupportedCodec.WrapF(stream.Codec)
		}
	}

	// rewind input
	_, err = input.Seek(0, io.SeekStart)
	if err != nil {
		return xo.W(err)
	}

	// prepare options
	opts := ffmpeg.ConvertOptions{
		Preset: ffmpeg.ImagePNG, // lossless
		Start:  report.Duration * pos,
	}

	// convert video
	err = ffmpeg.Convert(input, temp, opts)
	if err != nil {
		return xo.W(err)
	}

	// convert image
	err = p.ConvertImageFile(temp, output, sizer)
	if err != nil {
		return err
	}

	return nil
}

func (p *Processor) buffer(input io.Reader, temp bool, output func(*os.File) error, fn func(input *os.File, temp *os.File, output *os.File) error) error {
	// prepare paths
	id := uuid.New().String()
	inputPath := filepath.Join(p.config.Directory, id+"-input")
	tempPath := filepath.Join(p.config.Directory, id+"-temp")
	outputPath := filepath.Join(p.config.Directory, id+"-output")

	// create input file
	inputFile, err := os.Create(inputPath)
	if err != nil {
		return xo.W(err)
	}
	defer func() {
		_ = inputFile.Close()
		_ = os.Remove(inputPath)
	}()

	// copy input
	_, err = io.Copy(inputFile, input)
	if err != nil {
		return xo.W(err)
	}
	err = inputFile.Sync()
	if err != nil {
		return xo.W(err)
	}

	// rewind input
	_, err = inputFile.Seek(0, io.SeekStart)
	if err != nil {
		return xo.W(err)
	}

	// create temp file
	var tempFile *os.File
	if temp {
		tempFile, err = os.Create(tempPath)
		if err != nil {
			return xo.W(err)
		}
		defer func() {
			_ = tempFile.Close()
			_ = os.Remove(tempPath)
		}()
	}

	// create output file
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return xo.W(err)
	}
	defer func() {
		_ = outputFile.Close()
		_ = os.Remove(outputPath)
	}()

	// yield
	err = fn(inputFile, tempFile, outputFile)
	if err != nil {
		return err
	}

	// sync output
	err = outputFile.Sync()
	if err != nil {
		return xo.W(err)
	}

	// rewind output
	_, err = outputFile.Seek(0, io.SeekStart)
	if err != nil {
		return xo.W(err)
	}

	// yield output
	err = output(outputFile)
	if err != nil {
		return xo.W(err)
	}

	return nil
}
