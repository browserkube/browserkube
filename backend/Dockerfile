# syntax = docker/dockerfile:1-experimental

## This one is based on Debian
FROM golang:1.22.0-alpine3.19 AS builder
WORKDIR /app

# install dependencies
RUN apk add --no-cache build-base tzdata


# copy deps files
COPY backend/Makefile backend/go.mod backend/go.sum ./
COPY operator/go.mod operator/go.sum /operator/
COPY common/storage/go.mod common/storage/go.sum /common/storage/
# cache dependencies
RUN go mod download

ARG COMPONENT

COPY backend/pkg ./pkg
COPY backend/$COMPONENT ./$COMPONENT
COPY operator /operator
COPY common/storage /common/storage

# build an application
RUN --mount=type=cache,target=/root/.cache/go-build \
     make build component=$COMPONENT
#RUN make build component=$COMPONENT

# select image
#FROM gcr.io/distroless/static
FROM alpine:latest
WORKDIR /app
ARG COMPONENT
ARG USER_NAME=user
ARG USER_HOME=/home/user
ARG USER_UID=4096

RUN adduser --system \
    --home $USER_HOME \
    --uid $USER_UID \
    --ingroup root \
    --disabled-password \
    --shell /bin/bash \
    $USER_NAME

USER $USER_NAME

COPY --from=builder /app/bin/$COMPONENT/ /app/

CMD ["/app/app.bin"]
