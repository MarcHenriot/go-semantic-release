GIT_BASE_URL=github.com/MarcHenriot
BINARY_NAME=go-semantic-release

export PKG_CONFIG_PATH=/usr/lib/x86_64-linux-gnu/pkgconfig

.PHONY: debug deps sec fmt test

all: test fmt vet sec debug help

debug: bin/$(BINARY_NAME)
	@bin/$(BINARY_NAME) release -u https://github.com/argoproj/argo-cd.git

help: bin/$(BINARY_NAME)
	@bin/$(BINARY_NAME) -h

vet:
	@go vet ./...

fmt: go.mod go.sum
	@go fmt ./...

test:
	go test ./... -cover

go.mod:
	@go mod init $(GIT_BASE_URL)/$(BINARY_NAME)

go.sum: go.mod
	@go mod tidy

bin/$(BINARY_NAME): go.mod go.sum cmd/* pkg/**/*
	@go build -o bin/$(BINARY_NAME) cmd/cmd.go