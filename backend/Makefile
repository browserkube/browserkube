.DEFAULT_GOAL := build

COMMIT_HASH = `git rev-parse --short HEAD 2>/dev/null`
BUILD_DATE = `date +%FT%T%z`

GO=go
BINARY_DIR=bin
BINARY_NAME=app.bin
CURR_DIR = $(shell pwd)
GODIRS_NOVENDOR = $(shell go list ./... | grep -v /vendor/)
GOFILES_NOVENDOR = $(shell find . -type f -name '*.go' -not -path "./vendor/*")
BUILD_INFO_LDFLAGS=-ldflags "-w -extldflags '"-static"' -X buildInfo.buildDate=${BUILD_DATE} -X buildInfo.version=${COMMIT_HASH}"

.PHONY: build

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO) build ${BUILD_INFO_LDFLAGS} -o ${BINARY_DIR}/${component}/${BINARY_NAME} ./${component}

install-lint:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v1.54.2

lint:
	#docker run --rm -v $(CURR_DIR):/app -w /app golangci/golangci-lint:v1.32.2 golangci-lint run --enable-all --deadline 10m ./...
	bin/golangci-lint run --timeout=10m -v ./...

fmt:
	gofumpt -extra -l -w ${GOFILES_NOVENDOR}
	goimports -local github.com/browserkube/browserkube -w ${GOFILES_NOVENDOR}
	gci write --skip-generated --section Standard --section Default --section "Prefix(github.com/browserkube/browserkube)" ${GOFILES_NOVENDOR}

test:
	$(GO) test ${GODIRS_NOVENDOR}

build-image: build
	DOCKER_BUILDKIT=1 docker build -t browserkube .

release-image: build-image
	DOCKER_BUILDKIT=1 docker tag browserkube:latest quay.io/browserkube/browserkube
	DOCKER_BUILDKIT=1 docker push quay.io/browserkube/browserkube

release-helm: release-image
	echo '--- helm delete ---'
	helm delete -n ci browserkube
	echo '--- kustomize build ---'
	cd operator && kustomize build config/crd | kubectl apply -f -
	echo '--- install---'
	helm install --namespace ci browserkube ./helm
	
run:
	air -c .air.toml

swagger:
	swag init --parseDependency --parseInternal --generalInfo browserkube/main.go --output browserkube/docs/

generate-mocks:
	mockery --name Indexer --dir $GOPATH/pkg/mod/k8s.io/client-go@v0.29.3/tools/cache --output internal/provision/k8s/mocks