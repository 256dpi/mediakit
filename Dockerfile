FROM golang:alpine3.19
RUN apk add build-base ffmpeg vips-tools vips-poppler vips-heif chromium
RUN ffmpeg -version; vips --version; chromium-browser --version
