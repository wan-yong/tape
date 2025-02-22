# -------------------------------------------------------------
# This makefile defines the following targets
#   - tape - builds a binary program
#   - docker - build the tape image
#   - unit-test - runs the go-test based unit tests
#   - integration-test - runs the integration tests
#   - install - installs a binary program to GOBIN path

FABRIC_VERSION = latest
INTERGATION_CASE = ANDLogic

BASE_VERSION = 0.2.0
PREV_VERSION = 0.1.2

PROJECT_NAME = tape
DOCKERIMAGE = ghcr.io/hyperledger-twgc/tape
export DOCKERIMAGE
EXTRA_VERSION ?= $(shell git rev-parse --short HEAD)
BuiltTime ?= $(shell date)
PROJECT_VERSION=$(BASE_APISERVER_VERSION)-snapshot-$(EXTRA_VERSION)

PKGNAME = github.com/hyperledger-twgc/$(PROJECT_NAME)
CGO_FLAGS = CGO_CFLAGS=" "
ARCH=$(shell go env GOARCH)
MARCH=$(shell go env GOOS)-$(shell go env GOARCH)

# defined in pkg/infra/version.go
METADATA_VAR = Version=$(BASE_VERSION)
METADATA_VAR += CommitSHA=$(EXTRA_VERSION)

GO_LDFLAGS = $(patsubst %,-X $(PKGNAME)/pkg/infra/cmdImpl.%,$(METADATA_VAR))
GO_LDFLAGS += -X '$(PKGNAME)/pkg/infra/cmdImpl.BuiltTime=$(BuiltTime)'

GO_TAGS ?=

export GO_LDFLAGS GO_TAGS FABRIC_VERSION INTERGATION_CASE

tape:
	@echo "Building tape program......"
	go build -tags "$(GO_TAGS)" -ldflags "$(GO_LDFLAGS)" ./cmd/tape

escapes:
	@echo "go build check for escapes"
	go build -gcflags="-m -l" ./... | grep "escapes to heap" || true

.PHONY: docker
docker:
	@echo "Building tape docker......"
	docker build . --tag=ghcr.io/hyperledger-twgc/tape

.PHONY: unit-test
unit-test:
	@echo "Run unit test......"
	go test -v ./... --bench=. -cover --count=1

.PHONY: integration-test
integration-test:
	@echo "Run integration test......"
	./test/integration-test.sh $(FABRIC_VERSION) $(INTERGATION_CASE)

.PHONY: install
install:
	@echo "Install tape......"
	go install -tags "$(GO_TAGS)" -ldflags "$(GO_LDFLAGS)" ./cmd/tape

include gotools.mk

.PHONY: basic-checks
basic-checks: gotools-install linter

.PHONY: linter
linter:
	@echo "LINT: Running code checks......"
	./scripts/golinter.sh
