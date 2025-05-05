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
	Pages  int
	Delay  []int
}

// Analyze will run the vipsheader utility on the specified input and
// return the parsed report.
func Analyze(ctx context.Context, r io.Reader) (*Report, error) {
	// ensure context
	if ctx == nil {
		ctx = context.Background()
	}

	// prepare command
	cmd := exec.CommandContext(ctx, "vipsheader", "-a", "stdin")

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
		return nil, fmt.Errorf("vipsheader: %s", err.Error())
	}

	// get output
	lines := strings.Split(stdout.String(), "\n")
	line := lines[0]

	// strip filename if present
	if i := strings.Index(line, ":"); i >= 0 {
		line = line[i+1:]
	}

	// split string
	parts := strings.Split(strings.TrimSpace(line), ", ")

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

	// parse pages and delay
	pages := 1
	var delay []int
	for _, line := range lines {
		if strings.HasPrefix(line, "n-pages:") {
			pages, _ = strconv.Atoi(strings.TrimSpace(strings.Split(line, ":")[1]))
		} else if strings.HasPrefix(line, "delay:") {
			delayParts := strings.Split(strings.TrimSpace(strings.Split(line, ":")[1]), " ")
			for _, d := range delayParts {
				dInt, _ := strconv.Atoi(d)
				delay = append(delay, dInt)
			}
		}
	}

	// prepare report
	report := Report{
		Width:  width,
		Height: height,
		Bands:  bands,
		Color:  parts[2],
		Format: format,
		Pages:  pages,
		Delay:  delay,
	}

	return &report, nil
}
