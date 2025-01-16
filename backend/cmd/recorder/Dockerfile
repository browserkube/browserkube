# syntax = docker/dockerfile:1-experimental
FROM golang:1.23.1 as builder

RUN apt update && apt install --no-install-recommends -y ffmpeg x11-xserver-utils

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY cmd cmd/
COPY internal internal/

RUN --mount=type=cache,target=/root/.cache/go-build CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} go build -a -o recorder cmd/recorder/recorder.go


# TODO use distroless image
FROM builder as runtime
COPY --from=builder /usr/bin/ffmpeg /usr/bin/ffmpeg
COPY --from=builder /app/recorder /usr/bin/recorder
USER 65532:65532

ENTRYPOINT ["/usr/bin/recorder"]