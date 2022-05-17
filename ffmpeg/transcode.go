package ffmpeg

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// TODO: Support segmented encoding:
//  https://video.stackexchange.com/questions/32297/resuming-a-partially-completed-encode-with-ffmpeg

// Preset represents a transcoding preset.
type Preset string

// The available presets.
const (
	// AudioMP3VBRStandard is a standard MP3 variable encoding preset.
	// https://trac.ffmpeg.org/wiki/Encode/MP3
	AudioMP3VBRStandard = "audio-mp3-vbr-2"

	// VideoMP4H264Fast is a fast MP4 H.264 encoding preset.
	// https://trac.ffmpeg.org/wiki/Encode/H.264
	VideoMP4H264Fast = "video-mp4-h264-fast"
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
	case VideoMP4H264Fast:
		return []string{
			"-f", "mp4",
			"-codec", "libx264",
			"-preset", "fast",
			"-movflags", "+faststart",
			"-filter:v", "fps=30",
			"-movflags", "frag_keyframe",
		}
	default:
		return nil
	}
}

// TranscodeOptions define transcoding options.
type TranscodeOptions struct {
	Preset   Preset
	Duration float64
}

// Transcode will run the ffmpeg utility to transcode the specified input to the
// configured output.
func Transcode(r io.Reader, w io.Writer, opts TranscodeOptions) error {
	// check preset
	if !opts.Preset.Valid() {
		return fmt.Errorf("invalid preset")
	}

	// prepare args
	args := []string{
		"-i", "pipe:",
		"-nostats",
		"-hide_banner",
		"-progress", "pipe:3",
	}

	// apply preset
	args = append(args, opts.Preset.Args()...)

	// handle options
	if opts.Duration != 0 {
		args = append(args, "-t", strconv.FormatFloat(opts.Duration, 'f', -1, 64))
	}

	// finish args
	args = append(args, "pipe:4")

	// prepare output pipe
	or, ow, err := os.Pipe()
	if err != nil {
		return err
	}
	defer or.Close()
	defer ow.Close()

	// prepare progress pipe
	pr, pw, err := os.Pipe()
	if err != nil {
		return err
	}

	// prepare command
	cmd := exec.Command("ffmpeg", args...)

	// set input
	cmd.Stdin = r

	// set outputs
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// set output
	cmd.ExtraFiles = append(cmd.ExtraFiles, pw, ow)

	// handle progress
	go func() {
		scanner := bufio.NewScanner(pr)
		for scanner.Scan() {
			println(scanner.Text())
		}
	}()

	// copy output
	go func() {
		_, _ = io.Copy(w, or)
	}()

	// run command
	err = cmd.Run()
	if err != nil {
		if stderr.Len() > 0 {
			return fmt.Errorf(strings.ToLower(strings.TrimSpace(stderr.String())))
		}
		return err
	}

	return nil
}
