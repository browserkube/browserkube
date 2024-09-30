FROM golang:1.23.1-alpine as builder

WORKDIR /app/cmds/clipboard/

COPY cmds/clipboard/go.mod cmds/clipboard/go.sum ./
COPY /backend/go.mod /backend/go.sum /app/backend/

RUN go mod download
COPY cmds/clipboard/main.go main.go

COPY /backend /app/backend/

RUN --mount=type=cache,target=/root/.cache/go-build CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} go build -a -o clipboard main.go

FROM ubuntu:23.10

RUN apt-get update && \
    apt-get install -y xsel && \
    apt-get clean

COPY --from=builder /app/cmds/clipboard/clipboard /usr/bin/clipboard

ENTRYPOINT ["/usr/bin/clipboard"]