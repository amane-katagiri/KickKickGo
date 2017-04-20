NAME := kick-kick-go
VERSION := v0.2
REVISION := $(shell git rev-parse --short HEAD)

SRCS := $(shell find . -type f -name '*.go')
LDFLAGS := -ldflags="-s -w -X \"main.Version=$(VERSION)\" -X \"main.Revision=$(REVISION)\" -extldflags \"-static\""

DOCKER_IMAGE_NAME := amane/kick-kick-go
DOCKER_IMAGE_TAG  ?= latest
DOCKER_IMAGE      := $(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG)
DOCKER_WORK_DIR   := /go/src/github.com/amane-katagiri/kick-kick-go/

bin/$(NAME): $(SRCS)
	go build -a -tags netgo -installsuffix netgo $(LDFLAGS) -o bin/$(NAME)

.PHONY: glide
glide:
ifeq ($(shell command -v glide 2> /dev/null),)
	go get github.com/Masterminds/glide
	go install github.com/Masterminds/glide
endif

.PHONY: deps
deps: glide
	glide install

.PHONY: run
run:
	go run *.go

.PHONY: install
install:
	go install $(LDFLAGS)

.PHONY: clean
clean:
	rm -rf bin/*

.PHONY: lint
lint:
	 find . -type d -name "vendor" -prune -o -type f -name "*.go" -print | xargs -IXXX golint XXX

.PHONY: docker-build
docker-build:
	docker build -t $(DOCKER_IMAGE) .
