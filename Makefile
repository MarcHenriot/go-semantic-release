GIT_BASE_URL=github.com/MarcHenriot
BINARY_NAME=go-semantic-release

export PKG_CONFIG_PATH=/usr/lib/x86_64-linux-gnu/pkgconfig

.PHONY: run deps sec fmt test

all: test fmt vet sec run

run: go.mod go.sum
	@go run cmd/main.go

vet:
	@go vet ./...

fmt: go.mod go.sum
	@go fmt ./...

test:
	go test ./...

go.mod:
	@go mod init $(GIT_BASE_URL)/$(BINARY_NAME)

go.sum: go.mod
	@go mod tidy

bin/$(BINARY_NAME): go.mod go.sum cmd/*
	@go build -o bin/$(BINARY_NAME) cmd/main.go