package ffmpeg

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/samber/lo"
)

var imageCodecs = []string{
	"png",
	"mjpeg",
	"jpeg2000",
	"tiff",
	"webp",
}

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

// SideData defines stream side data.
type SideData struct {
	Rotation int `json:"rotation"`
}

// Stream is a ffprobe stream.
type Stream struct {
	// generic
	Type     string  `json:"codec_type"`
	Codec    string  `json:"codec_name"`
	Duration float64 `json:"duration,string"`

	// audio
	Channels   int `json:"channels"`
	SampleRate int `json:"sample_rate,string"`

	// video
	Width       int       `json:"width"`
	Height      int       `json:"height"`
	FrameRate   FrameRate `json:"r_frame_rate"`
	PixelFormat string    `json:"pix_fmt"`
	ColorSpace  string    `json:"color_space"`

	// other
	SideData []SideData `json:"side_data_list"`
}

// Report is a ffprobe report.
type Report struct {
	Duration float64
	Format   Format   `json:"format"`
	Streams  []Stream `json:"streams"`
	DidScan  bool
}

// Has returns whether as stream of the specified type is available.
func (r Report) Has(typ string) bool {
	for _, stream := range r.Streams {
		if stream.Type == typ {
			return true
		}
	}
	return false
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

// SampleRate returns the maximum stream sample rate.
func (r Report) SampleRate() int {
	// get sample rate
	var sampleRate int
	for _, stream := range r.Streams {
		if stream.SampleRate > sampleRate {
			sampleRate = stream.SampleRate
		}
	}

	return sampleRate
}

// Analyze will run the ffprobe and ffmpeg utilities on the specified input and
// return the parsed report. If the input is an *os.File and has a name it will
// be mapped via the filesystem. Otherwise, a pipe is created to connect the
// input. Using a file is recommended to allow ffprobe to seek within the file.
func Analyze(ctx context.Context, r io.Reader) (*Report, error) {
	// ensure context
	if ctx == nil {
		ctx = context.Background()
	}

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
	cmd := exec.CommandContext(ctx, "ffprobe", args...)

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
		_ = json.Unmarshal(stdout.Bytes(), &report)
		if report.Error.String != "" {
			return nil, fmt.Errorf(strings.ToLower(report.Error.String))
		}

		return nil, fmt.Errorf("ffprobe: %s", err.Error())
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

	// determine if image
	image := len(report.Streams) == 1 && lo.Contains(imageCodecs, report.Streams[0].Codec)

	// handle side data
	for i, stream := range report.Streams {
		for _, sd := range stream.SideData {
			if sd.Rotation == 90 || sd.Rotation == -90 {
				report.Streams[i].Width = stream.Height
				report.Streams[i].Height = stream.Width
			}
		}
		report.Streams[i].SideData = nil
	}

	// get seeker
	seeker, _ := r.(io.Seeker)

	// decode full file to get duration if still missing
	if !image && report.Duration == 0 && seeker != nil {
		// set flag
		report.DidScan = true

		// seek start
		_, err = seeker.Seek(0, io.SeekStart)
		if err != nil {
			return nil, err
		}

		// prepare command
		cmd = exec.CommandContext(ctx, "ffmpeg", "-nostats", "-hide_banner", "-i", "pipe:", "-f", "null", "-")

		// set input
		cmd.Stdin = r

		// run command
		out, err := cmd.CombinedOutput()
		if err != nil {
			if stderr.Len() > 0 {
				return nil, fmt.Errorf(strings.ToLower(strings.TrimSpace(stderr.String())))
			}
			return nil, fmt.Errorf("ffmpeg: %s", err.Error())
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
