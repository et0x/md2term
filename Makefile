.PHONY: build install test clean

build:
	go build -o bin/md2term cmd/md2term/main.go

install:
	go install ./cmd/md2term

test:
	go test ./...

clean:
	rm -rf bin