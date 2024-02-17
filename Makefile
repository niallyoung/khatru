SHELL := /bin/bash

all: build test cover
.PHONY: all

build:
	go build
.PHONY: build

test: build
	go test ./...
.PHONY: test

cover: build
	@go test \
		-coverprofile=coverage.out \
		-coverpkg $(go list github.com/fiatjaf/khatru/...) \
		./... 2>/dev/null 1>&2
	@./cover.sh
.PHONY: cover