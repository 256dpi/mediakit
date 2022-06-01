package ffmpeg

import (
	"log"
)

func init() {
	WarningsLogger = log.Default()
}
