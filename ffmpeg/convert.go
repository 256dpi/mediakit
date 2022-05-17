package ffmpeg

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// TODO: Support segmented encoding:
//  https://video.stackexchange.com/questions/32297/resuming-a-partially-completed-encode-with-ffmpeg

// WarningsLogger is the logger used to print warnings.
var WarningsLogger *log.Logger

// Preset represents a conversion preset.
// https://handbrake.fr/docs/en/1.5.0/technical/official-presets.html
type Preset int

// The available presets.
const (
	// AudioMP3VBRStandard is a standard MP3 variable encoding preset.
	// https://trac.ffmpeg.org/wiki/Encode/MP3
	AudioMP3VBRStandard = iota

	// VideoMP4H264AACFast is a fast MP4 H.264/AAC encoding preset.
	// https://trac.ffmpeg.org/wiki/Encode/H.264
	// https://trac.ffmpeg.org/wiki/Encode/AAC
	VideoMP4H264AACFast
)

// Valid returns whether the preset is valid.
func (p Preset) Valid() bool {
	return len(p.Args()) != 0
}

// Args returns the ffmpeg args for the preset.
func (p Preset) Args() []string {
	switch p {
	case AudioMP3VBRStandard:
		return []string{
			"-f", "mp3",
			"-codec:a", "libmp3lame",
			"-q:a", "2", // 170-210 kbit/s
			"-ac", "2", // stereo
		}
	case VideoMP4H264AACFast:
		return []string{
			"-f", "mp4",
			"-codec:v", "libx264",
			"-preset:v", "fast",
			"-movflags", "+faststart",
			"-movflags", "frag_keyframe",
			"-codec:a", "aac",
			"-q:a", "4", // 64-72 kbit/s/ch
			"-ac", "2", // stereo
		}
	default:
		return nil
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

	// Limit duration of the output.
	Duration float64

	// Apply scaling, set one component to -1 to keep aspect ratio.
	// https://trac.ffmpeg.org/wiki/Scaling
	Width, Height int

	// Force frames per second.
	FPS int

	// Receive progress updates.
	Progress func(Progress)
}

// Convert will run the ffmpeg utility to convert the specified input to the
// configured output.
func Convert(r io.Reader, w io.Writer, opts ConvertOptions) error {
	// check input and output
	rFile, _ := r.(*os.File)
	wFile, _ := w.(*os.File)
	rIsFile := rFile != nil && rFile.Name() != ""
	wIsFile := wFile != nil && wFile.Name() != ""

	// check preset
	if !opts.Preset.Valid() {
		return fmt.Errorf("invalid preset")
	}

	// prepare args
	args := []string{
		"-nostats",
		"-hide_banner",
		"-loglevel", "repeat+warning",
	}

	// add input
	if rIsFile {
		args = append(args, "-i", rFile.Name())
	} else {
		args = append(args, "-i", "pipe:")
	}

	// enable progress
	if opts.Progress != nil {
		args = append(args, "-progress", "pipe:3")
	}

	// apply preset
	args = append(args, opts.Preset.Args()...)

	// handle options
	if opts.Duration != 0 {
		args = append(args, "-t", strconv.FormatFloat(opts.Duration, 'f', -1, 64))
	}

	// handle filters
	if opts.Width != 0 || opts.Height != 0 || opts.FPS != 0 {
		args = append(args, "-filter:v")
		if opts.FPS != 0 {
			args = append(args, fmt.Sprintf("fps=%d", opts.FPS))
		}
		if opts.Width != 0 || opts.Height != 0 {
			args = append(args, fmt.Sprintf("scale=%d:%d", opts.Width, opts.Height))
		}
	}

	// finish args
	if wIsFile {
		args = append(args, wFile.Name())
	} else {
		args = append(args, "pipe:")
	}

	// prepare command
	cmd := exec.Command("ffmpeg", args...)

	// set input
	if !rIsFile {
		cmd.Stdin = r
	}

	// set outputs
	var stderr bytes.Buffer
	if !wIsFile {
		cmd.Stdout = w
	}
	cmd.Stderr = &stderr

	// handle progress
	if opts.Progress != nil {
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
					opts.Progress(progress)
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
		return err
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
