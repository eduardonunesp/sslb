.phony: build

build:
	@echo "Building..."
	@go build -o bin/sslb cli/*.go