package mediakit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
	"strings"
)

// Report is a ffprobe report.
type Report struct {
	Format  Format   `json:"format"`
	Streams []Stream `json:"streams"`
}

// Format is ffprobe format.
type Format struct {
	Name       string  `json:"format_name"`
	LongName   string  `json:"format_long_name"`
	ProbeScore float64 `json:"probe_score"`
}

// Stream is a ffprobe stream.
type Stream struct {
	// codec
	CodecName     string `json:"codec_name"`
	CodecLongName string `json:"codec_long_name"`
	CodecType     string `json:"codec_type"`

	// generic
	BitRate  int     `json:"bit_rate,string"`
	Duration float64 `json:"duration,string"`

	// audio
	SampleRate int `json:"sample_rate,string"`
	Channels   int `json:"channels"`

	// video
	Width  int `json:"width"`
	Height int `json:"height"`
}

// FFProbe will run the ffprobe utility on the specified input and return the
// parsed report.
func FFProbe(r io.Reader) (*Report, error) {
	// prepare args
	args := []string{
		"-print_format", "json",
		"-show_format",
		"-show_streams",
		"-show_error",
		"-",
	}

	// prepare command
	cmd := exec.Command("ffprobe", args...)

	// set input
	cmd.Stdin = r

	// set outputs
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// run command
	err := cmd.Run()
	if err != nil {
		// decode report
		var report struct {
			Error struct {
				String string `json:"string"`
			} `json:"error"`
		}
		err = json.Unmarshal(stdout.Bytes(), &report)
		if err != nil || report.Error.String == "" {
			return nil, fmt.Errorf("unkown error")
		}

		return nil, fmt.Errorf(strings.ToLower(report.Error.String))
	}

	// decode report
	var report Report
	err = json.Unmarshal(stdout.Bytes(), &report)
	if err != nil {
		return nil, err
	}

	return &report, nil
}
