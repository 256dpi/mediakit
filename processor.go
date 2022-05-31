package mediakit

import (
	"io"
	"io/fs"
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
	// The directory to use for temporary files.
	Directory string

	// The file system to use for temporary files.
	FS fs.FS

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
func (p *Processor) ConvertImage(source io.Reader, sink io.Writer, sizer Sizer) error {
	return p.buffer(source, sink, false, func(source, _, sink *os.File) error {
		return p.ConvertImageFile(source, sink, sizer)
	})
}

// ConvertImageFile will convert an image file with the specified sizer applied.
func (p *Processor) ConvertImageFile(source, sink *os.File, sizer Sizer) error {
	// analyze source
	report, err := vips.Analyze(source)
	if err != nil {
		return xo.W(err)
	}

	// check format
	if !lo.Contains(p.config.ImageFormats, report.Format) {
		return ErrUnsupportedFormat.WrapF(report.Format)
	}

	// rewind source
	_, err = source.Seek(0, io.SeekStart)
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
	err = vips.Convert(source, sink, opts)
	if err != nil {
		return xo.W(err)
	}

	return nil
}

// ConvertAudio will convert an audio stream.
func (p *Processor) ConvertAudio(source io.Reader, sink io.Writer, progress func(float64)) error {
	return p.buffer(source, sink, false, func(source, _, sink *os.File) error {
		return p.ConvertAudioFile(source, sink, progress)
	})
}

// ConvertAudioFile will convert an audio file.
func (p *Processor) ConvertAudioFile(source, sink *os.File, progress func(float64)) error {
	// analyze source
	report, err := ffmpeg.Analyze(source)
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

	// rewind source
	_, err = source.Seek(0, io.SeekStart)
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
	err = ffmpeg.Convert(source, sink, opts)
	if err != nil {
		return xo.W(err)
	}

	return nil
}

// ConvertVideo will convert a video stream with the specified sizer applied.
func (p *Processor) ConvertVideo(source io.Reader, sink io.Writer, sizer Sizer, progress func(float64)) error {
	return p.buffer(source, sink, false, func(source, _, sink *os.File) error {
		return p.ConvertVideoFile(source, sink, sizer, progress)
	})
}

// ConvertVideoFile will convert a video file with the specified sizer applied.
func (p *Processor) ConvertVideoFile(source, sink *os.File, sizer Sizer, progress func(float64)) error {
	// analyze source
	report, err := ffmpeg.Analyze(source)
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

	// rewind source
	_, err = source.Seek(0, io.SeekStart)
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
	err = ffmpeg.Convert(source, sink, opts)
	if err != nil {
		return xo.W(err)
	}

	return nil
}

// ExtractImage will extract an image stream from a video stream at the provided
// position with the specified sizer applied.
func (p *Processor) ExtractImage(source io.Reader, sink io.Writer, pos float64, sizer Sizer) error {
	return p.buffer(source, sink, true, func(source, temp, sink *os.File) error {
		return p.ExtractImageFile(source, temp, sink, pos, sizer)
	})
}

// ExtractImageFile will extract an image file from a video file at the
// provided position with the specified sizer applied.
func (p *Processor) ExtractImageFile(source, temp, sink *os.File, pos float64, sizer Sizer) error {
	// analyze source
	report, err := ffmpeg.Analyze(source)
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

	// rewind source
	_, err = source.Seek(0, io.SeekStart)
	if err != nil {
		return xo.W(err)
	}

	// prepare options
	opts := ffmpeg.ConvertOptions{
		Preset: ffmpeg.ImagePNG, // lossless
		Start:  report.Duration * pos,
	}

	// convert video
	err = ffmpeg.Convert(source, temp, opts)
	if err != nil {
		return xo.W(err)
	}

	// convert image
	err = p.ConvertImageFile(temp, sink, sizer)
	if err != nil {
		return err
	}

	return nil
}

func (p *Processor) buffer(source io.Reader, sink io.Writer, temp bool, fn func(source, temp, sink *os.File) error) error {
	// prepare paths
	id := uuid.New().String()
	sourcePath := filepath.Join(p.config.Directory, id+"-source")
	tempPath := filepath.Join(p.config.Directory, id+"-temp")
	sinkPath := filepath.Join(p.config.Directory, id+"-sink")

	// create source file
	sourceFile, err := os.Create(sourcePath)
	if err != nil {
		return xo.W(err)
	}
	defer func() {
		_ = sourceFile.Close()
		_ = os.Remove(sourcePath)
	}()

	// copy source
	_, err = io.Copy(sourceFile, source)
	if err != nil {
		return xo.W(err)
	}
	err = sourceFile.Sync()
	if err != nil {
		return xo.W(err)
	}

	// rewind source
	_, err = sourceFile.Seek(0, io.SeekStart)
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

	// create sink file
	sinkFile, err := os.Create(sinkPath)
	if err != nil {
		return xo.W(err)
	}
	defer func() {
		_ = sinkFile.Close()
		_ = os.Remove(sinkPath)
	}()

	// yield
	err = fn(sourceFile, tempFile, sinkFile)
	if err != nil {
		return err
	}

	// sync sink
	err = sinkFile.Sync()
	if err != nil {
		return xo.W(err)
	}

	// rewind sink
	_, err = sinkFile.Seek(0, io.SeekStart)
	if err != nil {
		return xo.W(err)
	}

	// copy sink
	_, err = io.Copy(sink, sinkFile)
	if err != nil {
		return xo.W(err)
	}

	return nil
}
