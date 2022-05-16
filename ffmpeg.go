package mediakit

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

// FFMPEGOptions define transcoding options.
type FFMPEGOptions struct {
	Format   string
	Duration int
}

// FFMPEG will run the ffmpeg utility to convert the specified input to the
// configured output.
func FFMPEG(r io.Reader, w io.Writer, opts FFMPEGOptions) error {
	// prepare args
	args := []string{
		"-i", "pipe:",
		"-nostats",
		"-hide_banner",
		"-progress", "pipe:3",
	}

	// handle options
	if opts.Format != "" {
		args = append(args, "-f", opts.Format)
	}
	if opts.Duration != 0 {
		args = append(args, "-t", strconv.Itoa(opts.Duration))
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
