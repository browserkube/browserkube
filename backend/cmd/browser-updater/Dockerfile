FROM golang:1.23.1-alpine as builder

WORKDIR /app

COPY jobs/browser-updater/go.mod jobs/browser-updater/go.sum ./
COPY /operator/go.mod /operator/go.sum /operator/


RUN go mod download
COPY jobs/browser-updater/main.go main.go
COPY jobs/browser-updater/internal internal/
COPY jobs/browser-updater/utils utils/

COPY /operator /operator/

RUN --mount=type=cache,target=/root/.cache/go-build CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} go build -a -o browser-updater main.go


FROM gcr.io/distroless/static as runtime
COPY --from=builder /app/browser-updater /usr/bin/browser-updater

ENTRYPOINT ["/usr/bin/browser-updater"]