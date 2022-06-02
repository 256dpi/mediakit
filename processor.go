package mediakit

import (
	"io"
	"os"
	"path/filepath"

	"github.com/256dpi/xo"
	"github.com/google/uuid"

	"github.com/256dpi/mediakit/ffmpeg"
	"github.com/256dpi/mediakit/vips"
)

// ErrMissingStream is returned if a required audio/video stream is missing.
var ErrMissingStream = xo.BF("missing stream")

// Config defines a Processor configuration.
type Config struct {
	// The directory to used for temporary files.
	Directory string

	// The used presets.
	ImagePreset vips.Preset   // vips.JPGWeb
	AudioPreset ffmpeg.Preset // ffmpeg.AudioMP3VBRStandard
	VideoPreset ffmpeg.Preset // ffmpeg.VideoMP4H264AACFast

	// The maximum video frame rate.
	MaxFrameRate float64 // 30
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

	return &Processor{
		config: config,
	}
}

// ConvertImage will convert an image stream with the specified sizer applied.
func (p *Processor) ConvertImage(input io.Reader, sizer Sizer, fn func(output *os.File) error) error {
	return p.buffer(input, false, fn, func(input, _, output *os.File) error {
		return ConvertImage(input, output, p.config.ImagePreset, sizer)
	})
}

// ConvertAudio will convert an audio stream.
func (p *Processor) ConvertAudio(input io.Reader, progress func(float64), fn func(output *os.File) error) error {
	return p.buffer(input, false, fn, func(input, _, output *os.File) error {
		return ConvertAudio(input, output, p.config.AudioPreset, progress)
	})
}

// ConvertVideo will convert a video stream with the specified sizer applied.
func (p *Processor) ConvertVideo(input io.Reader, sizer Sizer, progress func(float64), fn func(output *os.File) error) error {
	return p.buffer(input, false, fn, func(input, _, output *os.File) error {
		return ConvertVideo(input, output, p.config.VideoPreset, sizer, p.config.MaxFrameRate, progress)
	})
}

// ExtractImage will extract an image stream from a video stream at the provided
// position with the specified sizer applied.
func (p *Processor) ExtractImage(input io.Reader, pos float64, sizer Sizer, fn func(output *os.File) error) error {
	return p.buffer(input, true, fn, func(input, temp, output *os.File) error {
		return ExtractImage(input, temp, output, pos, p.config.ImagePreset, sizer)
	})
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
