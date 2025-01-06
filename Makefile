include .env
help:
	cat Makefile|ag "[a-zA-Z-0-9]+:"

SHELL := /bin/bash

# The name of the executable (default is current directory name)
TARGET := $(shell echo "$${PWD##*/}")
.DEFAULT_GOAL: $(TARGET)

# These will be provided to the target
VERSION := 1.0.0
BUILD := `git rev-parse HEAD`

# Use linker flags to provide version/build settings to the target
LDFLAGS=-ldflags "-X=main.Version=$(VERSION) -X=main.Build=$(BUILD)"
# go source files, ignore vendor directory
SRC = $(shell find . -type f -name '*.go' -not -path "./vendor/*")

.PHONY: all build clean install uninstall fmt simplify check run

$(TARGET): $(SRC)
	@GOOS=linux GOARCH=arm GOARM=5 go build $(LDFLAGS) -o $(TARGET)
x86: $(SRC)
	go build $(LDFLAGS) -o $(TARGET)
build: $(TARGET)
	@true
clean:
	@rm -f $(TARGET)
upload: clean build test
	scp $(TARGET) ${NAS_UPLOAD}
test:
	go test ./... -v
