.PHONY: build install clean test lint

build:
	go build -o tatami ./cmd/tatami

install:
	go install ./cmd/tatami

clean:
	rm -f tatami

test:
	go test -v ./...

lint:
	golangci-lint run
