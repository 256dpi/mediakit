package vips

import (
	"bytes"
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
	JPGWeb = iota

	// PNGWeb is a Web optimized preset for PNGs with stripped metadata,
	// optimized coding, SRGB color space and removed embedded color profiles.
	// `vips pngsave`
	PNGWeb = iota
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
func Convert(r io.Reader, w io.Writer, opts ConvertOptions) error {
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
	cmd := exec.Command("vips", args...)

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
		return err
	}

	return nil
}
