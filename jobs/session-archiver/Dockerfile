FROM golang:1.22.0-alpine as builder

WORKDIR /app

COPY jobs/session-archiver/go.mod jobs/session-archiver/go.sum ./
COPY /operator/go.mod /operator/go.sum /operator/
COPY /common/storage/go.mod /common/storage/go.sum /common/storage/


RUN go mod download
COPY jobs/session-archiver/cmd/ cmd/
COPY jobs/session-archiver/internal internal/

COPY /operator /operator/
COPY /common/storage /common/storage/

RUN --mount=type=cache,target=/root/.cache/go-build CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} go build -a -o session-archiver cmd/main.go


FROM gcr.io/distroless/static as runtime
COPY --from=builder /app/session-archiver /usr/bin/session-archiver

ENTRYPOINT ["/usr/bin/session-archiver"]