# Requirements: git, go, vgo
NAME     := gitpanda
VERSION  := $(shell cat VERSION)
REVISION := $(shell git rev-parse --short HEAD)

SRCS    := $(shell find . -type f -name '*.go')
LDFLAGS := "-s -w -X \"main.Version=$(VERSION)\" -X \"main.Revision=$(REVISION)\" -extldflags \"-static\""

.DEFAULT_GOAL := bin/$(NAME)

bin/$(NAME): $(SRCS)
	GO111MODULE=on go build -ldflags=$(LDFLAGS) -o bin/$(NAME)

.PHONY: gox
gox:
	gox -ldflags=$(LDFLAGS) -output="bin/gitpanda_{{.OS}}_{{.Arch}}"

.PHONY: clean
clean:
	rm -rf bin/*

.PHONY: tag
tag:
	git tag -a $(VERSION) -m "Release v$(VERSION)"
	git push --tags

.PHONY: release
release: tag
	git push origin master