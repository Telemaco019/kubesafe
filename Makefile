# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

# Setting SHELL to bash allows bash commands to be executed by recipes.
# Options are set to exit when a recipe line exits non-zero or a piped command fails.
SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec

.PHONY: all
all: build

##@ General

# The help target prints out all targets with their descriptions organized
# beneath their categories. The categories are represented by '##@' and the
# target descriptions by '##'. The awk commands is responsible for reading the
# entire set of makefiles included in this invocation, looking for lines of the
# file as xyz: ## something, and then pretty-format the target and help. Then,
# if there's a line with ##@ something, that gets pretty-printed as a category.
# More info on the usage of ANSI control characters for terminal formatting:
# https://en.wikipedia.org/wiki/ANSI_escape_code#SGR_parameters
# More info on the awk command:
# http://linuxcommand.org/lc3_adv_awk.php

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development

.PHONY: fmt
fmt: ## Run go fmt against code.
	go fmt ./...

.PHONY: vet
vet: ## Run go vet against code.
	go vet ./...

.PHONY: lint
lint: vet golangci-lint ## Run Go linter.
	$(GOLANGCI_LINT) run ./... -v

.PHONY: license-check
license-check: license-eye ## Check all files have the license header
	$(LICENSE_EYE) header check

.PHONY: license-fix
license-fix: license-eye ## Add license header to files that still don't have it
	$(LICENSE_EYE) header fix

.PHONY: test
test: ## Run tests.
	go test -v ./...

.PHONY: check
check: fmt vet lint test license-check ## Check the code


##@ Build

.PHONY: build
build: ## Build binary.
	go build -o bin/kubesafe kubesafe/kubesafe.go

##@ Build Dependencies

## Location to install dependencies to
LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)

## Tool Binaries
GOLANGCI_LINT ?= $(LOCALBIN)/golangci-lint
GORELEASER ?= $(LOCALBIN)/goreleaser
LICENSE_EYE ?= $(LOCALBIN)/license-eye
VHS ?= $(LOCALBIN)/vhs

## Tool Versions
GOLANGCI_LINT_VERSION ?= 2.5.0
GORELEASER_VERSION ?= 1.26.1

.PHONY: golangci-lint ## Download golanci-lint if necessary
golangci-lint: $(GOLANGCI_LINT)
$(GOLANGCI_LINT): $(LOCALBIN)
	test -s $(LOCALBIN)/golanci-lint || GOBIN=$(LOCALBIN) go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v${GOLANGCI_LINT_VERSION}


.PHONY: goreleaser ## Download goreleaser if necessary
goreleaser: $(GORELEASER)
$(GORELEASER): $(LOCALBIN)
	test -s $(LOCALBIN)/goreleaser || GOBIN=$(LOCALBIN) go install github.com/goreleaser/goreleaser@v${GORELEASER_VERSION}

.PHONY: license-eye ## Download license-eye if necessary
license-eye: $(LICENSE_EYE)
$(LICENSE_EYE): $(LOCALBIN)
	test -s $(LOCALBIN)/license-eye || GOBIN=$(LOCALBIN) go install github.com/apache/skywalking-eyes/cmd/license-eye@latest

.PHONY: vhs
vhs: $(VHS)
$(VSH): $(LOCALBIN)
	test -s $(LOCALBIN)/vhs || GOBIN=$(LOCALBIN) go install github.com/charmbracelet/vhs@latest
