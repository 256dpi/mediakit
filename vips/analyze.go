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

// Report is an analysis report.
type Report struct {
	Width  int
	Height int
	Bands  int
	Color  string
	Format string
}

// Analyze will run the vipsheader utility on the specified input and
// return the parsed report.
func Analyze(ctx context.Context, r io.Reader) (*Report, error) {
	// ensure context
	if ctx == nil {
		ctx = context.Background()
	}

	// prepare command
	cmd := exec.CommandContext(ctx, "vipsheader", "stdin")

	// set input
	cmd.Stdin = r

	// set outputs
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// run command
	err := cmd.Run()
	if err != nil {
		if stderr.Len() > 0 {
			return nil, fmt.Errorf(strings.ToLower(strings.TrimSpace(stderr.String())))
		}
		return nil, err
	}

	// split string
	parts := strings.Split(strings.TrimSpace(stdout.String()), ", ")

	// parse size
	var width, height int
	size := strings.Split(parts[0], "x")
	if len(size) == 2 {
		width, _ = strconv.Atoi(size[0])
		height, _ = strconv.Atoi(size[1][:strings.Index(size[1], " ")])
	}

	// parse bands
	bands, _ := strconv.Atoi(parts[1][:strings.Index(parts[1], " ")])

	// parse format
	format := strings.TrimSuffix(parts[3], "load_source")

	// prepare report
	report := Report{
		Width:  width,
		Height: height,
		Bands:  bands,
		Color:  parts[2],
		Format: format,
	}

	return &report, nil
}
