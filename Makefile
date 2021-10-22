all: tidy mod-download build
.PHONY: all

export GO111MODULE := on
export GOPROXY := direct
#export GOOS := linux
#export GOARCH := 386

tidy:
	@echo "==> Tidying module"
	@go mod tidy
.PHONY: tidy

mod-download:
	@echo "==> Downloading Go module"
	@go mod download
.PHONY: mod-download

build:
	@echo "==> Building"
	@env GOOS=linux GOARCH=386 go build -o bin/tw
.PHONY: build