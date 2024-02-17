SHELL := /bin/bash

test:
	go test \
		-coverprofile=coverage.out \
		-coverpkg $(go list github.com/fiatjaf/khatru/...) \
		./...
.PHONY: test

cover:
	@$(MAKE) test 2>/dev/null 1>&2
	@./cover.sh
.PHONY: cover