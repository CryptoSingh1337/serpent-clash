# Makefile

# Variables
GO_APP_NAME=server
GO_APP_PATH=./server
GO_APP_OUTPUT_PATH=../../
VUE_APP_PATH=./client

# Targets
.PHONY: all build serve clean

all: build

## Build Go and Vue applications
build: build-go build-vue

## Build Go application
build-go:
	@echo "Building Go server..."
	cd $(GO_APP_PATH)/cmd/app && go build -o $(GO_APP_OUTPUT_PATH)/$(GO_APP_NAME).exe

## Build Vue application
install-npm-dep:
	@echo "Installing npm dependencies"
	cd $(VUE_APP_PATH) && npm install

build-vue:
	@echo "Building Vue app..."
	cd $(VUE_APP_PATH) && npm run build

## Serve Go and Vue applications
serve: build-go build-vue serve-go

## Serve Go application
serve-go:
	@echo "Running Go server..."
	cd $(GO_APP_PATH) && $(GO_APP_NAME).exe

## Serve Vue application
serve-vue:
	@echo "Running Vue development server..."
	cd $(VUE_APP_PATH) && npm run serve

## Clean Go and Vue build files
clean: clean-go clean-vue

clean-go:
	@echo "Cleaning Go server build..."
	rm -f $(GO_APP_PATH)/$(GO_APP_NAME)

clean-vue:
	@echo "Cleaning Vue build..."
	rm -rf $(VUE_APP_PATH)/dist

## Help
help:
	@echo "Usage:"
	@echo "  make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  build        Build Go and Vue applications"
	@echo "  serve        Serve Go and Vue applications"
	@echo "  clean        Clean build files"
	@echo "  help         Show this help message"

