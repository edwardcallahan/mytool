VERSION ?= dev

pkgs := $(shell go list ./... | grep -v tck)
files := $(shell find . -name '*.go' -print)

.PHONY: all
all: format test build-binaries

VERSION:=$(shell ./scripts/version.sh)
BUILDTIME:=$(shell date +%FT%T%z)
binary=mytool

.PHONY: install
install:
	go install -ldflags "-X main.version=${VERSION}" cmd/

.PHONY: format
format:
	goimports -w $(files)

.PHONY: test
test: checkformat vet lint gotest

.PHONY: gotest
gotest:
	go test -race $(pkgs)

.PHONY: gotestnocache
gotestnocache:
	go clean -testcache
	go test -race $(pkgs)

.PHONY: build-binaries
build-binaries:
	go build -o "./bin/${binary}"

.PHONY: docker-build
docker-build:
	docker image build -t mytool .