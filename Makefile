# --- PREAMBLE
SHELL := bash
.ONESHELL:
.SHELLFLAGS := -eu -o pipefail -c
MAKEFLAGS += --warn-undefined-variables

BUILD_DATE := $(shell date +'%Y%m%d')
VERSION := $(shell support/get-version.sh)

build:
	go build -o target/gomicsv \
    -ldflags="-X main.buildDate=$(BUILD_DATE) -X main.versionString=$(VERSION)" \
    cmd/gomicsv/main.go
.PHONY: build

clean:
	rm -rf target/
.PHONY: clean

dist: build
	rm -rf target/dist
	rm -rf target/*.zip
	mkdir -p target/dist/gomicsv
	cp target/gomicsv target/dist/gomicsv/gomicsv
	cp README.md target/dist/gomicsv
	cp COPYING.md target/dist/gomicsv
	cp CHANGELOG.md target/dist/gomicsv
	pushd target/dist
	zip -r ../Gomics-v_$(VERSION).zip gomicsv
	rm -rf gomicsv
.PHONY: dist
