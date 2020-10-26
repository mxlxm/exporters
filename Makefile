PROJECT_NAME=$(shell basename $(PWD))
BUILD_DIR=$(PWD)/build

BUILD=go build -o $(BUILD_DIR)/bin/$(PROJECT_NAME) && echo done || exit -1

default: proxy clean create build

proxy:
	@go env -w GOPROXY=https://goproxy.io,direct
	@go env -w GOSUMDB=gosum.io+ce6e7565+AY5qEHUk/qmHc5btzW45JVoENfazw8LielDsaI+lEbq6

clean:
	@rm -rf build/bin/$(PROJECT_NAME)
	@rm -rf build/logs/*

create:
	@mkdir -p build/{bin,conf,logs};

build:
	@printf "building...\n"
	@$(BUILD)

linux:
	@printf "building...\n"
	@GOARCH=amd64 CGO_ENABLED=0 GOOS=linux $(BUILD)
	@tar -czf $(PROJECT_NAME).tar.gz build/*

.PHONY: proxy clean create build linux
