PROJECT := github.com/anmolbabu/contact-book
GITCOMMIT := $(shell git rev-parse --short HEAD 2>/dev/null)
PKGS := $(shell go list  ./... | grep -v $(PROJECT)/vendor)
BUILD_FLAGS := -ldflags="-w -X $(PROJECT)/cmd.GITCOMMIT=$(GITCOMMIT)"
FILES := contact-book dist
DEFAULT_DB_DIR := /tmp/plivo/contact-book

default: bin

.PHONY: bin
bin:
	go build ${BUILD_FLAGS} -o contact-book main.go

.PHONY: install
install:
	go install ${BUILD_FLAGS} .

.PHONY: validate-vendor-licenses
validate-vendor-licenses:
	go get github.com/frapposelli/wwhrd
	wwhrd check

.PHONY: validate
validate: gofmt check-vendor vet validate-vendor-licenses lint

.PHONY: gofmt
gofmt:
	./scripts/check-gofmt.sh

.PHONY: check-vendor
check-vendor:
	./scripts/check-vendor.sh

.PHONY: lint
lint:
	golint $(PKGS)

.PHONY: run
run: install
ifdef $(PLIVO_DB_DIR)
	mkdir -p $(PLIVO_DB_DIR)
else
	mkdir -p $(DEFAULT_DB_DIR)
endif
	./contact-book

.PHONY: test
test: install
	go test -v $(PKGS)
