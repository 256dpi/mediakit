package ffmpeg

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// WarningsLogger is the logger used to print warnings.
var WarningsLogger *log.Logger

// Preset represents a conversion preset.
// https://handbrake.fr/docs/en/1.5.0/technical/official-presets.html
type Preset int

// The available presets.
const (
	// AudioMP3VBRStandard is a standard MP3 variable encoding preset.
	// https://trac.ffmpeg.org/wiki/Encode/MP3
	AudioMP3VBRStandard Preset = iota + 1

	// VideoMP4H264AACFast is a fast MP4 H.264/AAC encoding preset.
	// https://trac.ffmpeg.org/wiki/Encode/H.264
	// https://trac.ffmpeg.org/wiki/Encode/AAC
	VideoMP4H264AACFast

	// ImageJPEG is a basic JPEG encoding preset.
	ImageJPEG

	// ImagePNG is a basic PNG encoding preset.
	ImagePNG

	// ImageWebP is a basic WebP encoding preset.
	ImageWebP

	// AnimGIF is a basic GIF encoding preset.
	AnimGIF

	// AnimWebP is a basic WebP animation encoding preset.
	AnimWebP
)

// Valid returns whether the preset is valid.
func (p Preset) Valid() bool {
	return len(p.Args(false)) != 0
}

// Args returns the ffmpeg args for the preset.
func (p Preset) Args(isFile bool) []string {
	switch p {
	case AudioMP3VBRStandard:
		return []string{
			"-f", "mp3",
			"-codec:a", "libmp3lame",
			"-q:a", "2", // 170-210 kbit/s
			"-ac", "2", // stereo
		}
	case VideoMP4H264AACFast:
		args := []string{
			"-f", "mp4",
			"-codec:v", "libx264",
			"-preset:v", "fast",
			"-colorspace:v", "bt709",
			"-color_primaries:v", "bt709",
			"-color_trc:v", "bt709",
			"-color_range:v", "tv",
			"-movflags", "+faststart",
			"-codec:a", "aac",
			"-q:a", "4", // 64-72 kbit/s/ch
			"-ac", "2", // stereo
		}
		if !isFile {
			args = append(args, "-movflags", "frag_keyframe")
		}
		return args
	case ImageJPEG:
		return []string{
			"-f", "image2",
			"-update", "1",
			"-frames:v", "1",
			"-codec:v", "mjpeg",
			"-q:v", "3",
		}
	case ImagePNG:
		return []string{
			"-f", "image2",
			"-update", "1",
			"-frames:v", "1",
			"-codec:v", "png",
		}
	case ImageWebP:
		return []string{
			"-f", "image2",
			"-update", "1",
			"-frames:v", "1",
			"-codec:v", "libwebp",
			"-q:v", "90",
		}
	case AnimGIF:
		return []string{
			"-f", "gif",
			"-codec:v", "gif",
			"-q:v", "3",
			"-loop", "0",
		}
	case AnimWebP:
		return []string{
			"-f", "webp",
			"-codec:v", "libwebp",
			"-preset:v", "default",
			"-q:v", "90",
			"-compression_level", "6",
			"-loop", "0",
		}
	default:
		return nil
	}
}

// Filters returns the ffmpeg filters for the preset.
func (p Preset) Filters() []string {
	switch p {
	case VideoMP4H264AACFast:
		return []string{
			// h264 requires even height
			`pad=ceil(iw/2)*2:ceil(ih/2)*2`,
			// pixel format and color space conversion
			"format=yuv420p",
			"scale=in_color_matrix=auto:in_range=auto:out_color_matrix=bt709:out_range=tv",
		}
	case ImageJPEG:
		return []string{"format=yuvj444p"}
	case ImagePNG:
		return []string{"format=rgb24"}
	default:
		return nil
	}
}

// ScaleFlags returns additional ffmpeg scale flags for the preset.
func (p Preset) ScaleFlags() string {
	switch p {
	case AnimGIF:
		return ":flags=lanczos"
	default:
		return ""
	}
}

// Progress is emitted during conversion.
type Progress struct {
	Duration float64
	Size     int64
}

// ConvertOptions defines conversion options.
type ConvertOptions struct {
	// Select the desired preset.
	Preset Preset

	// Set the start of the output.
	Start float64

	// Limit duration of the output.
	Duration float64

	// Apply scaling, set one part to -1 to keep the aspect ratio.
	// https://trac.ffmpeg.org/wiki/Scaling
	Width, Height int

	// Force a frame rate.
	FrameRate float64

	// Force a sample rate.
	SampleRate int

	// Receive progress updates.
	ProgressFunc func(Progress)
	ProgressRate time.Duration
}

// Convert will run the ffmpeg utility to convert the specified input to the
// configured output. If the input or output is an *os.File and has a name, it
// will be mapped via the filesystem. Otherwise, pipes are created to connect
// the input or output. Using files is recommended to allow ffmpeg to seek
// within the files.
func Convert(ctx context.Context, r io.Reader, w io.Writer, opts ConvertOptions) error {
	// ensure context
	if ctx == nil {
		ctx = context.Background()
	}

	// check input and output
	rFile, _ := r.(*os.File)
	wFile, _ := w.(*os.File)
	rIsFile := rFile != nil && rFile.Name() != ""
	wIsFile := wFile != nil && wFile.Name() != ""

	// check preset
	if !opts.Preset.Valid() {
		return fmt.Errorf("invalid preset")
	}

	// generate palette for GIF images$
	var palette *os.File
	if opts.Preset == AnimGIF {
		// check input
		if !rIsFile {
			return fmt.Errorf("GIF requires file input")
		}

		// prepare command
		cmd := exec.CommandContext(ctx, "ffmpeg", []string{
			"-nostats", "-hide_banner", "-loglevel", "repeat+warning", "-y", "-i", rFile.Name(),
			"-vf", "palettegen", "-f", "image2pipe", "-vcodec", "png", "pipe:",
		}...)

		// run command
		out, err := cmd.Output()
		if err != nil {
			return fmt.Errorf("ffmpeg: %s", err.Error())
		}

		// setup palette pipe
		var pw *os.File
		palette, pw, err = os.Pipe()
		if err != nil {
			return err
		}
		go func() {
			defer pw.Close()
			_, _ = pw.Write(out)
		}()
	}

	// prepare args
	args := []string{
		"-nostats",
		"-hide_banner",
		"-loglevel", "repeat+warning",
		"-y", // overwrite
	}

	// handle early options
	if opts.Start != 0 {
		args = append(args, "-ss", strconv.FormatFloat(opts.Start, 'f', -1, 64))
	}

	// add input(s)
	if rIsFile {
		args = append(args, "-i", rFile.Name())
	} else {
		args = append(args, "-i", "pipe:")
	}
	if palette != nil {
		args = append(args, "-i", "pipe:3")
	}

	// enable progress
	if opts.ProgressFunc != nil && opts.ProgressRate > 0 {
		args = append(args,
			"-progress",
			"pipe:3",
			"-stats_period",
			strconv.FormatFloat(opts.ProgressRate.Seconds(), 'f', -1, 64),
		)
	}

	// prepare filters
	var filters []string

	// add scale filter
	if opts.Width != 0 || opts.Height != 0 {
		filters = append(filters, fmt.Sprintf("scale=%d:%d%s", opts.Width, opts.Height, opts.Preset.ScaleFlags()))
	}

	// apply preset filters
	filters = append(filters, opts.Preset.Filters()...)

	// append filter arg
	if len(filters) > 0 && palette == nil {
		args = append(args, "-filter:v", strings.Join(filters, ", "))
	} else if len(filters) > 0 && palette != nil {
		scale := strings.Join(filters, ",")
		args = append(args, "-filter_complex", fmt.Sprintf("[0:v]%s[x];[x][1:v]paletteuse", scale))
	} else if palette != nil {
		args = append(args, "-filter_complex", "[0:v][1:v]paletteuse")
	}

	// append preset args (output)
	args = append(args, opts.Preset.Args(wIsFile)...)

	// handle options
	if opts.Duration != 0 {
		args = append(args, "-t", strconv.FormatFloat(opts.Duration, 'f', -1, 64))
	}
	if opts.FrameRate != 0 {
		args = append(args, "-r", strconv.FormatFloat(opts.FrameRate, 'f', -1, 64))
	}
	if opts.SampleRate != 0 {
		args = append(args, "-ar", strconv.Itoa(opts.SampleRate))
	}

	// finish args
	if wIsFile {
		args = append(args, wFile.Name())
	} else {
		args = append(args, "pipe:")
	}

	// prepare command
	cmd := exec.CommandContext(ctx, "ffmpeg", args...)

	// set input
	if !rIsFile {
		cmd.Stdin = r
	}

	// set palette input pipe if needed
	if palette != nil {
		cmd.ExtraFiles = append(cmd.ExtraFiles, palette)
	}

	// set outputs
	var stderr bytes.Buffer
	if !wIsFile {
		cmd.Stdout = w
	}
	cmd.Stderr = &stderr

	// handle progress
	if opts.ProgressFunc != nil && opts.ProgressRate > 0 {
		// prepare progress pipe
		pr, pw, err := os.Pipe()
		if err != nil {
			return err
		}

		// set output
		cmd.ExtraFiles = append(cmd.ExtraFiles, pw)

		go func() {
			// prepare variables
			var progress Progress

			// scan output
			scanner := bufio.NewScanner(pr)
			for scanner.Scan() {
				line := scanner.Text()
				if strings.HasPrefix(line, "out_time=") {
					// parse duration
					durStr := strings.TrimPrefix(line, "out_time=")
					duration, _ := parseDuration(durStr)
					progress.Duration = duration.Seconds()
				} else if strings.HasPrefix(line, "total_size=") {
					// parse size
					sizeStr := strings.TrimPrefix(line, "total_size=")
					progress.Size, _ = strconv.ParseInt(sizeStr, 10, 64)
				} else if strings.HasPrefix(line, "progress=") {
					// emit and clear progress
					opts.ProgressFunc(progress)
					progress = Progress{}
				}
			}
		}()
	}

	// run command
	err := cmd.Run()
	if err != nil {
		if stderr.Len() > 0 {
			return fmt.Errorf(strings.ToLower(strings.TrimSpace(stderr.String())))
		}
		return fmt.Errorf("ffmpeg: %s", err.Error())
	}

	// print warnings
	if WarningsLogger != nil {
		scanner := bufio.NewScanner(&stderr)
		for scanner.Scan() {
			WarningsLogger.Print(scanner.Text())
		}
	}

	return nil
}
