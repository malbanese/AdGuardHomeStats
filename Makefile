export GO111MODULE=on

#################
## ENVIRONMENT ##
#################

APP=AdGuardHome
OUT_ROOT="./out"
BIN_ROOT="${OUT_ROOT}/bin"
REPORT_ROOT="${OUT_ROOT}/report"
SHELL := /bin/bash

#############
## Linting ##
#############

# Runs all quality checks
check-all: check-format check-lint check-vet

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

# Chore to update the go dependency files
chore-deps:
	go mod tidy
	go mod vendor

#############
## Testing ##
#############

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

build-linux:
	GOOS=linux GOARCH=amd64 go build -o "${BIN_ROOT}/adguardhome-linux-amd64" ./cmd/proxy/main.go

build-darwin:
	GOOS=darwin GOARCH=amd64 go build -o "${BIN_ROOT}/adguardhome-darwin-amd64" ./cmd/proxy/main.go

build-windows:
	GOOS=windows GOARCH=amd64 go build -o "${BIN_ROOT}/adguardhome-windows-amd64.exe" ./cmd/proxy/main.go

build-all: build-linux build-darwin build-windows

clean:
	go clean
	rm -rf out/

#############
## General ##
#############

.PHONY: all check check-format check-lint check-vet chore-deps test coverage build-linux build-darwin build-windows build-all clean
all: check test build

