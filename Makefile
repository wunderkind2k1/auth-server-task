.PHONY: all build test clean lint

all: build

build:
	go build -o bin/auth-server ./cmd/auth-server

test:
	go test -v ./...

clean:
	rm -rf bin/

lint:
	cd server && golangci-lint run ./...
	cd keytool && golangci-lint run ./...
