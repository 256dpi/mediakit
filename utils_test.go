package mediakit

import (
	"log"

	"github.com/256dpi/mediakit/ffmpeg"
)

func init() {
	ffmpeg.WarningsLogger = log.Default()
}
