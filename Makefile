MAIN_PACKAGE_PATH := ./cmd/corg
BINARY_NAME := corg

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

## tidy: format code and tidy modfile
.PHONY: tidy
tidy:
	go fmt ./...
	go mod tidy -v

## audit: run quality control checks
.PHONY: audit
audit:
	go mod verify
	go vet ./...
	go run honnef.co/go/tools/cmd/staticcheck@latest -checks=all,-ST1000,-U1000 ./...
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...
	go test -race -buildvcs -vet=off ./...

## build: build the application
.PHONY: build
build:
	go build -o=./bin/${BINARY_NAME} ${MAIN_PACKAGE_PATH}
	chmod +x ./bin/${BINARY_NAME}

## run: run the  application
.PHONY: run
run:
	go run ${MAIN_PACKAGE_PATH}

## test: run all tests
.PHONY: test
test:
	go test -v -race -buildvcs ./...

## test/cover: run all tests and display coverage
.PHONY: test/cover
test/cover:
	go test -v -race -buildvcs -coverprofile=./bin/coverage.out ./...
	go tool cover -html=./bin/coverage.out