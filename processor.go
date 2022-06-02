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

// Processor provides methods to convert streams and files as well as extract
// images from video streams or files.
type Processor struct {
	dir string
}

// NewProcessor creates and returns a new processor using the specified config.
func NewProcessor(directory string) *Processor {
	return &Processor{
		dir: directory,
	}
}

// ConvertImage will convert an image stream with the specified sizer applied.
func (p *Processor) ConvertImage(input io.Reader, preset vips.Preset, sizer Sizer, fn func(output *os.File) error) error {
	return p.buffer(input, false, fn, func(input, _, output *os.File) error {
		return ConvertImage(input, output, preset, sizer)
	})
}

// ConvertAudio will convert an audio stream.
func (p *Processor) ConvertAudio(input io.Reader, preset ffmpeg.Preset, maxSampleRate int, progress func(float64), fn func(output *os.File) error) error {
	return p.buffer(input, false, fn, func(input, _, output *os.File) error {
		return ConvertAudio(input, output, preset, maxSampleRate, progress)
	})
}

// ConvertVideo will convert a video stream with the specified sizer applied.
func (p *Processor) ConvertVideo(input io.Reader, preset ffmpeg.Preset, sizer Sizer, maxFrameRate float64, progress func(float64), fn func(output *os.File) error) error {
	return p.buffer(input, false, fn, func(input, _, output *os.File) error {
		return ConvertVideo(input, output, preset, sizer, maxFrameRate, progress)
	})
}

// ExtractImage will extract an image stream from a video stream at the provided
// position with the specified sizer applied.
func (p *Processor) ExtractImage(input io.Reader, pos float64, preset vips.Preset, sizer Sizer, fn func(output *os.File) error) error {
	return p.buffer(input, true, fn, func(input, temp, output *os.File) error {
		return ExtractImage(input, temp, output, pos, preset, sizer)
	})
}

func (p *Processor) buffer(input io.Reader, temp bool, output func(*os.File) error, fn func(input *os.File, temp *os.File, output *os.File) error) error {
	// prepare paths
	id := uuid.New().String()
	inputPath := filepath.Join(p.dir, id+"-input")
	tempPath := filepath.Join(p.dir, id+"-temp")
	outputPath := filepath.Join(p.dir, id+"-output")

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

	// yield output
	err = output(outputFile)
	if err != nil {
		return xo.W(err)
	}

	return nil
}
