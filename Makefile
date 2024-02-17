SHELL:=/bin/sh

NAME:=khatru
HASH:=$(shell git rev-parse --short HEAD)

all: test cover
.PHONY: all

lint:
	golangci-lint run --timeout=5m

test:
	go test ./...
.PHONY: test

cover:
	@go test \
		-coverprofile=coverage.out \
#		-covermode=atomic \
		-coverpkg $(go list github.com/fiatjaf/khatru/...) \
		./... 1>/dev/null 2>&1
	@./cover.sh
.PHONY: cover

docker.build:
	docker build . \
		-f Dockerfile -t $(NAME):$(HASH) \
		--build-arg BUILD_REVISION=$(HASH)
	docker tag $(NAME):$(HASH) $(NAME):latest

docker.lint: docker.build
	docker run --rm -v $(PWD):/app $(NAME):$(HASH) make lint

docker.test: docker.build
	docker run --rm -v $(PWD):/app $(NAME):$(HASH) make test

docker.cover: docker.build
	docker run --rm -v $(PWD):/app $(NAME):$(HASH) make cover

docker.shell: docker.build
	docker run --rm -it -v $(PWD):/app $(NAME):$(HASH) /bin/sh
