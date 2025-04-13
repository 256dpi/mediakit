FROM golang:alpine3.21
RUN apk add build-base ffmpeg vips-tools vips-poppler vips-heif chromium
RUN ffmpeg -version; vips --version; chromium-browser --version
