.PHONY: samples

all:
	go fmt ./...
	go vet ./...
	staticcheck ./...

samples:
	# audio (https://filesamples.com/categories/audio)
	wget -nc -O ./samples/sample.wav https://filesamples.com/samples/audio/wav/sample3.wav || true
	wget -nc -O ./samples/sample.aac https://filesamples.com/samples/audio/aac/sample3.aac || true
	wget -nc -O ./samples/sample.aiff https://filesamples.com/samples/audio/aiff/sample3.aiff || true
	wget -nc -O ./samples/sample.flac https://filesamples.com/samples/audio/flac/sample3.flac || true
	wget -nc -O ./samples/sample.m4a https://filesamples.com/samples/audio/m4a/sample3.m4a || true
	wget -nc -O ./samples/sample.mp2 https://filesamples.com/samples/audio/mp2/sample3.mp2 || true
	wget -nc -O ./samples/sample.mp3 https://filesamples.com/samples/audio/mp3/sample3.mp3 || true
	wget -nc -O ./samples/sample.ogg https://filesamples.com/samples/audio/ogg/sample3.ogg || true
	wget -nc -O ./samples/sample.wma https://filesamples.com/samples/audio/wma/sample3.wma || true
	# video (https://filesamples.com/categories/video)
	wget -nc -O ./samples/sample.hevc https://filesamples.com/samples/video/hevc/sample_1280x720.hevc || true
	wget -nc -O ./samples/sample.avi https://filesamples.com/samples/video/avi/sample_1280x720.avi || true
	wget -nc -O ./samples/sample.mov https://filesamples.com/samples/video/mov/sample_1280x720.mov || true
	wget -nc -O ./samples/sample.mp4 https://filesamples.com/samples/video/mp4/sample_1280x720.mp4 || true
	wget -nc -O ./samples/sample.mpeg https://filesamples.com/samples/video/mpeg/sample_1280x720.mpeg || true
	wget -nc -O ./samples/sample.mpg https://filesamples.com/samples/video/mpg/sample_1280x720.mpg || true
	wget -nc -O ./samples/sample.webm https://filesamples.com/samples/video/webm/sample_1280x720.webm || true
	wget -nc -O ./samples/sample.wmv https://filesamples.com/samples/video/wmv/sample_1280x720.wmv || true
	# combined (https://www.elecard.com/videos)
	wget -nc -O ./samples/bunny.ts https://www.elecard.com/storage/video/bbb_360p_c.ts || true
