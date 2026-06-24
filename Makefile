.PHONY: all
all: test

VERSION := $(shell git describe --tags --abbrev=0 2>/dev/null || echo "unknown")
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS := -ldflags "-X main.Version=$(VERSION) -X main.GitCommit=$(GIT_COMMIT) -X main.BuildTime=$(BUILD_TIME)"

.PHONY: build
build:
	go build $(LDFLAGS) ./cmd/

.PHONY: clean
clean:
	rm ../golang-github-k0swe-wsjtx-go*

.PHONY: test
test:
	go test ./...
	go vet ./...
