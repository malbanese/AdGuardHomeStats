export GO111MODULE=on

#################
## ENVIRONMENT ##
#################

APP=adguardhomestats
OUT_ROOT=./out
BIN_ROOT=${OUT_ROOT}/bin
REPORT_ROOT=${OUT_ROOT}/report
SHELL := /bin/bash

#############
## Linting ##
#############
.PHONY: check-all check-format check-lint check-vet

# Runs all quality checks
check: check-format check-lint check-vet

# Uses the go formatter
check-format:
	go fmt ./...

# Uses the golangci linter
check-lint:
	golangci-lint run

# Uses the go vet static analysis tool
check-vet:
	go vet ./... 

############
## Chores ##
############
.PHONY: chore-deps

# Chore to update the go dependency files
chore-deps:
	go mod tidy
	go mod vendor

#############
## Testing ##
#############
.PHONY: test test-coverage

# Run all tests
test: chore-deps
	mkdir -p "${REPORT_ROOT}"
	go test -v -timeout 10m ./... -coverprofile="${REPORT_ROOT}/coverage.out" -json > "${REPORT_ROOT}/report.json"

# Run tests with coverage
test-coverage: 
	make test
	go tool cover -html="${REPORT_ROOT}/coverage.out"

##############
## Building ##
##############
.PHONY: build-linux build-darwin build-windows build-linux-armv6 clean

build: build-linux build-linux-armv6 build-windows

build-linux:
	GOOS=linux GOARCH=amd64 go build -o "${BIN_ROOT}/${APP}-linux-amd64" ./cmd/main.go

build-linux-armv6:
	GOOS=linux GOARCH=arm GOARM=6 go build -o "${BIN_ROOT}/${APP}-linux-armv6" ./cmd/main.go

build-windows:
	GOOS=windows GOARCH=amd64 go build -o "${BIN_ROOT}/${APP}-amd64.exe" ./cmd/main.go

clean:
	go clean
	rm -rf out/

#############
## General ##
#############
.PHONY: all

all: check test build

