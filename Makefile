export GO111MODULE=on
GOOS := $(shell go env GOOS)
APP_PATH := scraper.go

# determins the variables based on GOOS 
ifeq ($(GOOS), windows)
    RM = del /Q
	HOME = $(shell echo %USERPROFILE%)
	CONFIG_PATH = $(subst  ,,$(HOME)\.golangci.yaml)
	OUTPUT_PATH = C:\cli-tools
	DATE=$(shell powershell -Command "(Get-Date).ToString('yyyy-MM-ddTHH:mm:sszzz')")
	DIRTY=$(shell powershell -Command "if (git status --porcelain) { echo true } else { echo false }")
else
    RM = rm -f
	HOME = $(shell echo $$HOME)
	CONFIG_PATH = $(HOME)/.golangci.yaml
	OUTPUT_PATH = /usr/local/bin
	DATE=$(shell date +"%Y-%m-%dT%H:%M:%S%z")
	DIRTY=$(shell git status --porcelain | grep . > /dev/null && echo true || echo false)
endif

GO_BUILD_LDFLAGS=\
  -X go.szostok.io/version.version=$(shell git describe --tags --always) \
  -X go.szostok.io/version.commit=$(shell git rev-parse --short HEAD) \
  -X go.szostok.io/version.buildDate=$(DATE) \
  -X go.szostok.io/version.commitDate=$(shell git log -1 --date=format:"%Y-%m-%dT%H:%M:%S%z" --format=%cd) \
  -X go.szostok.io/version.dirtyBuild=$(DIRTY)

check-quality:
	@make tidy
	@make fmt
	@make vet
#@make lint

lint:
	@make fmt
	@make vet
	golangci-lint run --config="$(CONFIG_PATH)" ./...

vet:
	go vet ./...

fmt:
	go fmt ./...

tidy:
	go mod tidy

update-packages:
	go get -u all

update-common:
	go get github.com/ondrovic/common@latest

test:
	make tidy
	go test -v -timeout 10m ./... -coverprofile=unit.coverage.out || (echo "Tests failed. See report.json for details." && exit 1)

coverage:
	make test
	go tool cover -html=coverage.out -o coverage.html

build: 
	go build -ldflags="$(GO_BUILD_LDFLAGS)" -o $(OUTPUT_PATH) $(APP_PATH)
.PHONY: build

all:
	make check-quality
	make test
	make build

clean:
	go clean
	$(RM) *coverage*
	$(RM) *report*
	$(RM) *lint*

vendor:
	go mod vendor