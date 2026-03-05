.PHONY: build install uninstall clean test lint release

VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS := -ldflags "-s -w -X main.version=$(VERSION)"

build:
	go build $(LDFLAGS) -o tatami ./cmd/tatami

install: build
	go install $(LDFLAGS) ./cmd/tatami

uninstall:
	rm -f $(shell go env GOPATH)/bin/tatami

clean:
	rm -f tatami

test:
	go test -v ./...

lint:
	golangci-lint run

# Create a new release (run: make release VERSION=v0.1.0)
release:
	@if [ -z "$(VERSION)" ]; then echo "VERSION required"; exit 1; fi
	git tag -a $(VERSION) -m "Release $(VERSION)"
	git push origin $(VERSION)
	@echo "Release $(VERSION) created. GitHub Actions will build and publish."
