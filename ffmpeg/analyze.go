package ffmpeg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"os/exec"
	"strings"
	"time"
)

// Report is a ffprobe report.
type Report struct {
	Duration float64
	Format   Format   `json:"format"`
	Streams  []Stream `json:"streams"`
}

// Format is ffprobe format.
type Format struct {
	Name       string  `json:"format_name"`
	LongName   string  `json:"format_long_name"`
	ProbeScore float64 `json:"probe_score"`
	Duration   float64 `json:"duration,string"`
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

// Analyze will run the ffprobe utility on the specified input and return the
// parsed report.
func Analyze(r io.Reader, reset func() error) (*Report, error) {
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

	// find duration
	report.Duration = report.Format.Duration
	for _, stream := range report.Streams {
		report.Duration = math.Max(report.Duration, stream.Duration)
	}

	// TODO: If file is too big we may just count the read bytes from the source
	//  to estimate the progress of the operation.

	// decode full file to get duration if still missing
	if report.Duration == 0 && reset != nil {
		// reset reader
		err = reset()
		if err != nil {
			return nil, err
		}

		// prepare command
		cmd = exec.Command("ffmpeg", "-nostats", "-hide_banner", "-i", "pipe:", "-f", "null", "-")

		// set input
		cmd.Stdin = r

		// run command
		out, err := cmd.CombinedOutput()
		if err != nil {
			if stderr.Len() > 0 {
				return nil, fmt.Errorf(strings.ToLower(strings.TrimSpace(stderr.String())))
			}
			return nil, err
		}

		// find duration string
		var durStr string
		lines := strings.Split(strings.TrimSpace(string(out)), "\n")
		for i := len(lines) - 1; i >= 0; i-- {
			if strings.Contains(lines[i], " time=") {
				parts := strings.Split(lines[i], " ")
				for _, part := range parts {
					if strings.HasPrefix(part, "time=") {
						durStr = part[5:]
						break
					}
				}
				if durStr != "" {
					break
				}
			}
		}

		// parse duration
		duration, err := parseDuration(durStr)
		if err != nil {
			return nil, err
		}

		// set duration
		report.Duration = duration.Seconds()
	}

	return &report, nil
}

func parseDuration(str string) (time.Duration, error) {
	// parse string
	ts, err := time.Parse("15:04:05.999999999", str)
	if err != nil {
		return 0, err
	}

	// add year
	ts = ts.AddDate(1970, 0, 0)

	// get duration
	dur := time.Duration(ts.UnixNano())

	return dur, nil
}
