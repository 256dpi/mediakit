package ffmpeg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// Format is a ffprobe format.
type Format struct {
	Name     string  `json:"format_name"`
	Duration float64 `json:"duration,string"`
}

// FrameRate is a video frame rate.
type FrameRate float64

// UnmarshalJSON implements the json.Unmarshaler interface.
func (r *FrameRate) UnmarshalJSON(bytes []byte) error {
	str := strings.Trim(string(bytes), `"`)
	parts := strings.Split(str, "/")
	if len(parts) == 0 || len(parts) > 2 {
		return nil
	} else if len(parts) == 1 {
		f, err := strconv.ParseFloat(parts[0], 64)
		if err != nil {
			return err
		}
		*r = FrameRate(f)
	} else {
		f1, err := strconv.ParseFloat(parts[0], 64)
		if err != nil {
			return err
		}
		f2, err := strconv.ParseFloat(parts[1], 64)
		if err != nil {
			return err
		}
		f := f1 / f2
		if !math.IsNaN(f) {
			*r = FrameRate(f)
		}
	}
	return nil
}

// Stream is a ffprobe stream.
type Stream struct {
	// codec
	Type  string `json:"codec_type"`
	Codec string `json:"codec_name"`

	// generic
	BitRate  int     `json:"bit_rate,string"`
	Duration float64 `json:"duration,string"`

	// audio
	SampleRate int `json:"sample_rate,string"`
	Channels   int `json:"channels"`

	// video
	Width     int       `json:"width"`
	Height    int       `json:"height"`
	FrameRate FrameRate `json:"r_frame_rate"`
}

// Report is a ffprobe report.
type Report struct {
	Duration float64
	Format   Format   `json:"format"`
	Streams  []Stream `json:"streams"`
	DidParse bool
}

// Size returns the maximum stream width and height.
func (r Report) Size() (int, int) {
	// get size
	var width, height int
	for _, stream := range r.Streams {
		if stream.Width > 0 && stream.Height > 0 {
			if stream.Width > width {
				width = stream.Width
			}
			if stream.Height > height {
				height = stream.Height
			}
		}
	}

	return width, height
}

// FrameRate returns the maximum stream frame rate.
func (r Report) FrameRate() float64 {
	// get frame rate
	var frameRate float64
	for _, stream := range r.Streams {
		frameRate = math.Max(frameRate, float64(stream.FrameRate))
	}

	return frameRate
}

// Analyze will run the ffprobe and ffmpeg utilities on the specified input and
// return the parsed report. If the input is an *os.File and has a name it will
// be mapped via the filesystem. Otherwise, a pipe is created to connect the
// input. Using a file is recommended to allow ffprobe to seek within the file.
func Analyze(r io.Reader) (*Report, error) {
	// check input
	file, _ := r.(*os.File)
	isFile := file != nil && file.Name() != ""

	// prepare args
	args := []string{
		"-print_format", "json",
		"-show_format",
		"-show_streams",
		"-show_error",
	}

	// add input
	if isFile {
		args = append(args, file.Name())
	} else {
		args = append(args, "pipe:")
	}

	// prepare command
	cmd := exec.Command("ffprobe", args...)

	// set input
	if !isFile {
		cmd.Stdin = r
	}

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

	// get seeker
	seeker, _ := r.(io.Seeker)

	// decode full file to get duration if still missing
	if report.Duration == 0 && seeker != nil {
		// set flag
		report.DidParse = true

		// seek start
		_, err = seeker.Seek(0, io.SeekStart)
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
