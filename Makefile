SHELL := /bin/bash

all: build test cover
.PHONY: all

build:
	go build
.PHONY: build

test: build
	go test ./...
.PHONY: test

cover:
	@$(MAKE) build >/dev/null
	@go test \
		-coverprofile=coverage.out \
		-coverpkg $(go list github.com/fiatjaf/khatru/...) \
		./... 1>/dev/null 2>&1
	@./cover.sh
.PHONY: cover