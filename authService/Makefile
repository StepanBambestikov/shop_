OSNAME				:=
BINARY_NAME_FILE	:=
ifeq ($(OS),Windows_NT)
	OSNAME=windows
else
	UNAME_S :=$(shell uname -s)
	ifeq ($(UNAME_S),Linux)
		OSNAME=linux
	endif
	ifeq ($(UNAME_S),Darwin)
		OSNAME=drawin
	endif
endif

# Env
CGO_ENABLED=1
GOCMD=go
GOARCH=amd64
BINARY_NAME=server
BINARY_NAME_FILE =./dist/$(OSNAME)/
BINARY_NAME_LINUX=./dist/linux/
BINARY_NAME_MACOS=./dist/darwin/
BINARY_NAME_WIN=./dist/windows/
GIT_COMMIT=$(shell git rev-list -1 HEAD)
VERSION=$(shell date "+%Y.%m.%d.%H:%M:%S")
GIT_TAG=$(shell git describe --all)
BUILD_FLAGS=-v -ldflags "-X main.GitCommit=$(GIT_COMMIT) -X main.Version=$(GIT_TAG) -X main.BuildDate=$(VERSION)"
prebuild:
	mkdir -p ./dist/$(OSNAME)/
prebuild-all:
	mkdir -p $(BINARY_NAME_LINUX)
	mkdir -p $(BINARY_NAME_MACOS)
	mkdir -p $(BINARY_NAME_WIN)
_dist_os:
	$(GOCMD) build $(BUILD_FLAGS) -o $(BINARY_NAME_FILE) ./cmd/...
build: prebuild _dist_os
build-linux:
	GOOS=linux CGO_ENABLED=0 $(GOCMD) build $(BUILD_FLAGS) -o $(BINARY_NAME_LINUX) ./cmd/...
build-mac:
	GOOS=darwin CGO_ENABLED=0 $(GOCMD) build $(BUILD_FLAGS) -o $(BINARY_NAME_MACOS) ./cmd/...
test:
	$(GOCMD) test -v ./...
clean:
	$(GOCMD) clean ./...
	rm -rf ./dist/
doc:
	$(GOCMD) run github.com/swaggo/swag/cmd/swag@v1.8.12 init --parseDependency --parseInternal --output ./docs/gen -g root.go --dir ./internal/app
download:
	go mod tidy
	go mod download
build-all: build-mac build-linux
all: doc test prebuild-all build-all
run:
	./dist/drawin/goauth start
