FROM golang:1.23.1-alpine3.20

WORKDIR /app

RUN apk add --update --no-cache \
      ca-certificates \
      make git curl tzdata && \
    curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b $(go env GOPATH)/bin

# copy deps files
COPY backend/.air.toml backend/Makefile backend/go.mod backend/go.sum /app/backend/
COPY operator/go.mod operator/go.sum /app/operator/

WORKDIR /app/backend

# cache dependencies
RUN go mod download

COPY backend/pkg /app/backend/pkg
COPY backend/browserkube /app/backend/browserkube
COPY operator /app/operator

## this is just to populate build cache
RUN go build -o ./bin/browserkube ./browserkube

CMD ["make", "run"]
