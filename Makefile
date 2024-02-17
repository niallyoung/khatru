SHELL := /bin/bash

test:
	go test ./...
.PHONY: test

cover:
	@go test \
		-coverprofile=coverage.out \
		-coverpkg $(go list github.com/fiatjaf/khatru/...) \
		./... 2>/dev/null 1>&2
	@./cover.sh
.PHONY: cover