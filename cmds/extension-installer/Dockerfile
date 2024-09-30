FROM golang:1.23.1 as builder

RUN apt update && apt install --no-install-recommends

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY cmd cmd/
COPY preinstaller.go utils.go browsers.go ./

RUN --mount=type=cache,target=/root/.cache/go-build CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} go build -a -o extension-installer cmd/preinstaller/preinstaller.go


FROM alpine:3.18 as runtime
COPY --from=builder /app/extension-installer /usr/bin/extension-installer

ENTRYPOINT ["/usr/bin/extension-installer"]