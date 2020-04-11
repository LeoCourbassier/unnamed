PROJECT_NAME := "Unnamed"
PKG_LIST := $(shell go list)
GO_FILES := $(shell find . -name '*.go' | grep -v _test.go)

all: clean build

lint: ## Lint the files
	@echo "Checking for linting errors"
	@golint -set_exit_status ${PKG_LIST}

test: ## Run unittests
	@echo "Running tests"
	@go test -short ${PKG_LIST}

race: dep ## Run data race detector
	@echo "Running race detector"
	@go test -race -short ${PKG_LIST}

msan: dep ## Run memory sanitizer
	@echo "Running memory sanitizer"
	@go test -msan -short ${PKG_LIST}

dep: ## Get the dependencies
	@echo "Getting dependencies"
	@go get -v -d ./...

build: dep ## Build the binary file
	@echo "Building the binary"
	@go build -i -v

clean: ## Remove previous build
	@echo "Cleaning the previous build"
	@rm -f $(PROJECT_NAME)
	@rm -f campaign-options_*

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
