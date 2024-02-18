.phony: build clean dockerize

build:
	@echo "Building..."
	@go build -o bin/sslb cli/*.go

clean:
	@echo "Cleaning..."
	@rm -rf bin

dockerize:
	@echo "Building Docker image..."
	@docker build -t sslb .