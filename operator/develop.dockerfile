FROM golang:1.22.0

ARG UID=65532
ARG GID=65532
ARG TARGETOS
ARG TARGETARCH

ENV CGO_ENABLED=0
ENV GOOS=${TARGETOS:-linux}
ENV GOARCH=${TARGETARCH}

RUN apt-get update && apt-get install -y \
      ca-certificates tzdata && \
      rm -rf /var/lib/apt/lists/* && \
    curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b $(go env GOPATH)/bin

RUN addgroup --gid $GID nonroot && \
    adduser --uid $UID --gid $GID --disabled-password --gecos "" nonroot && \
    echo 'nonroot ALL=(ALL) NOPASSWD: ALL' >> /etc/sudoers

USER nonroot
WORKDIR /home/nonroot/app

# Copy the Go Modules manifests
COPY --chown=nonroot:nonroot go.mod go.sum .air.toml ./
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the go source
COPY --chown=nonroot:nonroot cmd/main.go cmd/main.go
COPY --chown=nonroot:nonroot api/ api/
COPY --chown=nonroot:nonroot internal/controller/ internal/controller/

## this is just to populate build cache
RUN go build -a -o ./bin/manager cmd/main.go

ENTRYPOINT ["air", "-c", ".air.toml", "--"]
