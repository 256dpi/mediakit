package samples

import (
	"embed"
	"io"
	"os"
	"path/filepath"
)

//go:embed files
var files embed.FS

// The available image samples.
const (
	ImageGIF    = "image.gif"
	ImageHEIF   = "image.heif"
	ImageJPEG   = "image.jpg"
	ImageJPEG2K = "image.jpf"
	ImagePDF    = "image.pdf"
	ImagePNG    = "image.png"
	ImageTIFF   = "image.tiff"
	ImageWebP   = "image.webp"
)

// The available audio samples.
const (
	AudioAAC   = "audio.aac"
	AudioAIFF  = "audio.aif"
	AudioFLAC  = "audio.flac"
	AudioMPEG2 = "audio.mp2"
	AudioMPEG3 = "audio.mp3"
	AudioMPEG4 = "audio.m4a" // aac
	AudioOGG   = "audio.ogg" // vorbis
	AudioWAV   = "audio.wav"
	AudioWMA   = "audio.wma"
)

// The available video samples.
const (
	VideoAVI   = "video.avi" // h264/aac
	VideoFLV   = "video.flv" // h263/mp3
	VideoGIF   = "video.gif"
	VideoMKV   = "video.mkv" // h265/ac3
	VideoMOV   = "video.mov" // h264/aac
	VideoMPEG  = "video.mpeg"
	VideoMPEG2 = "video.mpg"
	VideoMPEG4 = "video.mp4"  // h264/aac
	VideoOGG   = "video.ogv"  // theora/flac
	VideoWebM  = "video.webm" // vp9/vorbis
	VideoWMV   = "video.wmv"
)

// Images returns all image samples.
func Images() []string {
	return []string{
		ImageGIF,
		ImageHEIF,
		ImageJPEG,
		ImageJPEG2K,
		ImagePDF,
		ImagePNG,
		ImageTIFF,
		ImageWebP,
	}
}

// Audio returns all audio samples.
func Audio() []string {
	return []string{
		AudioAAC,
		AudioAIFF,
		AudioFLAC,
		AudioMPEG2,
		AudioMPEG3,
		AudioMPEG4,
		AudioOGG,
		AudioWAV,
		AudioWMA,
	}
}

// Video returns all video samples.
func Video() []string {
	return []string{
		VideoAVI,
		VideoFLV,
		VideoGIF,
		VideoMKV,
		VideoMOV,
		VideoMPEG,
		VideoMPEG2,
		VideoMPEG4,
		VideoOGG,
		VideoWebM,
		VideoWMV,
	}
}

// Load will load the specified sample.
func Load(sample string) io.ReadCloser {
	// open file
	stream, err := files.Open(filepath.Join("files", sample))
	if err != nil {
		panic(err)
	}

	return stream
}

// Buffer will load and buffer a sample.
func Buffer(sample string) *os.File {
	// create file
	file, err := os.Create(filepath.Join(os.TempDir(), "mediakit-samples-"+sample))
	if err != nil {
		panic(err)
	}

	// copy data
	_, err = io.Copy(file, Load(sample))
	if err != nil {
		panic(err)
	}

	// rewind file
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		panic(err)
	}

	return file
}
