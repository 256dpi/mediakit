FROM golang:alpine
RUN apk add build-base ffmpeg vips-tools vips-poppler
RUN ffmpeg -version; vips --version
