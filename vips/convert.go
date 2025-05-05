package vips

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os/exec"
	"strconv"
	"strings"
)

// Preset represents a conversion preset.
type Preset int

// The available preset.
const (
	// JPGWeb is a Web optimized preset for JPGs with stripped metadata,
	// optimized coding, SRGB color space and removed embedded color profiles.
	// `vips jpegsave`
	JPGWeb Preset = iota + 1

	// PNGWeb is a Web optimized preset for PNGs with stripped metadata, sRGB
	// color space and removed embedded color profiles.
	// `vips pngsave`
	PNGWeb

	// WebP is a Web optimized preset for WebP with stripped metadata,
	// optimized coding, SRGB color space and removed embedded color profiles.
	// `vips webpsave`
	WebP
)

// Valid returns whether the preset is valid.
func (p Preset) Valid() bool {
	return p.Arg() != ""
}

// Arg returns the vips argument for the preset.
func (p Preset) Arg() string {
	switch p {
	case JPGWeb:
		return ".jpg[Q=90,strip,optimize_coding]"
	case PNGWeb:
		return ".png[Q=90,strip]"
	case WebP:
		return ".webp[Q=90,strip,smart_subsample]"
	default:
		return ""
	}
}

// ConvertOptions defines conversion options.
type ConvertOptions struct {
	// Select the desired preset.
	Preset Preset

	// The mandatory output width.
	Width int

	// The optional output height.
	Height int

	// Whether to fill area and crop.
	Crop bool

	// Whether to keep embedded color profile (no sRGB conversion).
	KeepProfile bool

	// Whether to skip metadata rotation.
	NoRotate bool
}

// Convert will run the vips utility to convert the specified input to the
// configured output.
func Convert(ctx context.Context, r io.Reader, w io.Writer, opts ConvertOptions) error {
	// ensure context
	if ctx == nil {
		ctx = context.Background()
	}

	// check preset
	if !opts.Preset.Valid() {
		return fmt.Errorf("invalid preset")
	}

	// prepare args
	args := []string{
		"thumbnail_source",
		"[descriptor=0]",
		opts.Preset.Arg(),
		strconv.Itoa(opts.Width),
	}

	// handle height
	if opts.Height != 0 {
		args = append(args, "--height", strconv.Itoa(opts.Height))
	}

	// handle crop
	if opts.Crop {
		args = append(args, "--crop", "centre")
	}

	// handle profile
	if !opts.KeepProfile {
		args = append(args, "--export-profile", "srgb")
	}

	// handle no rotate
	if opts.NoRotate {
		args = append(args, "--no-rotate")
	}

	// prepare command
	cmd := exec.CommandContext(ctx, "vips", args...)

	// set input {
	cmd.Stdin = r

	// set outputs
	var stderr bytes.Buffer
	cmd.Stdout = w
	cmd.Stderr = &stderr

	// run command
	err := cmd.Run()
	if err != nil {
		if stderr.Len() > 0 {
			return fmt.Errorf(strings.ToLower(strings.TrimSpace(stderr.String())))
		}
		return fmt.Errorf("vips: %s", err.Error())
	}

	return nil
}

// Pipeline will run the vips utility multiple times to convert the specified
// input to the configured output using a pipeline of operations. Operations
// are standard vips CLI operations with the command name and "stdin" input
// argument omitted.
func Pipeline(ops [][]string, r io.Reader, w io.Writer) error {
	// prepare stderr
	var stderr bytes.Buffer

	// prepare commands
	var list []*exec.Cmd
	for i, args := range ops {
		// check args
		if len(args) == 0 {
			return fmt.Errorf("empty args")
		}

		// create command
		cmd := exec.Command("vips", append([]string{args[0], "stdin"}, args[1:]...)...)
		cmd.Stdin = r

		// set up stdout pipe unless it's the last command
		if i == len(ops)-1 {
			cmd.Stdout = w
		} else {
			pipeR, pipeW, err := os.Pipe()
			if err != nil {
				return fmt.Errorf("pipe error: %w", err)
			}
			cmd.Stdout = pipeW
			r = pipeR
			defer pipeW.Close()
		}

		// set up stderr
		cmd.Stderr = &stderr

		// add process
		list = append(list, cmd)
	}

	// start all commands
	for i, cmd := range list {
		if err := cmd.Start(); err != nil {
			if stderr.Len() > 0 {
				err = fmt.Errorf(strings.ToLower(strings.TrimSpace(stderr.String())))
			}
			return fmt.Errorf("vips: %s: %s", ops[i][0], err.Error())
		}
	}

	// wait for all commands
	for i, cmd := range list {
		if err := cmd.Wait(); err != nil {
			if stderr.Len() > 0 {
				err = fmt.Errorf(strings.ToLower(strings.TrimSpace(stderr.String())))
			}
			return fmt.Errorf("vips: %s: %s", ops[i][0], err.Error())
		}
	}

	return nil
}
