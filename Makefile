ROOT := $(shell git rev-parse --show-toplevel)
OS := $(shell uname -s | awk '{print tolower($$0)}')
ARCH := amd64
PROJECT := gotway
VERSION := $(shell git describe --abbrev=0 --tags)
LD_FLAGS := -X main.version=$(VERSION) -s -w
SOURCE_FILES ?= ./internal/... ./pkg/... ./cmd/...

export CGO_ENABLED := 0
export GO111MODULE := on
export GOBIN := $(shell pwd)/bin
export PORT ?= 9111
export ENV ?= local
export LOG_LEVEL ?= debug
export REDIS_URL ?= redis://localhost:6379/11
export KUBECONFIG ?= $(HOME)/.kube/config
export HA ?= true
export HEALTH_CHECK ?= true
export CACHE ?= true
export METRICS ?= true
export PPROF ?= false

.PHONY: all
all: help

.PHONY: help
help:	### Show targets documentation
ifeq ($(OS), linux)
	@grep -P '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
else
	@awk -F ':.*###' '$$0 ~ FS {printf "%15s%s\n", $$1 ":", $$2}' \
		$(MAKEFILE_LIST) | grep -v '@awk' | sort
endif

KIND := $(GOBIN)/kind
KIND_VERSION := v0.14.0
kind:
	$(call go-install,sigs.k8s.io/kind@$(KIND_VERSION))

GOLANGCI_LINT := $(GOBIN)/golangci-lint
GOLANGCI_LINT_VERSION := v1.46.2
golangci-lint:
	$(call go-install,github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION))

GORELEASER := $(GOBIN)/goreleaser
GORELEASER_VERSION := v1.9.2
goreleaser:
	$(call go-install,github.com/goreleaser/goreleaser@$(GORELEASER_VERSION))

MOCKERY := $(GOBIN)/mockery
MOCKERY_VERSION := v2.12.3
mockery:
	$(call go-install,github.com/vektra/mockery/v2@$(MOCKERY_VERSION))

.PHONY: generate
generate: vendor ### Generate code
	@bash ./hack/hack.sh

.PHONY: install-crds
install-crds: ## Install CRDs
	@kubectl apply -f manifests/crds

.PHONY: cluster
cluster: kind ### Create a KIND cluster
	$(KIND) create cluster --name $(PROJECT)

.PHONY: docker
docker: ### Spin up docker dependencies
	@docker compose up -d

clean: ### Clean build files
	@rm -rf ./bin
	@go clean

.PHONY: build
build: clean ### Build binary
	@go build -tags netgo -a -v -ldflags "${LD_FLAGS}" -o ./bin/gotway ./cmd/gotway/*.go
	@chmod +x ./bin/*

.PHONY: run
run: lint docker install-crds ### Quick run
	@CGO_ENABLED=1 go run -race cmd/gotway/*.go

.PHONY: deps
deps: ### Optimize dependencies
	@go mod tidy

.PHONY: vendor
vendor: ### Vendor dependencies
	@go mod vendor

.PHONY: lint
lint: golangci-lint ### Lint
	$(GOLANGCI_LINT) run

.PHONY: release
release: goreleaser ### Dry-run release
	$(GORELEASER) release --snapshot --rm-dist

.PHONY: test-clean
test-clean: ### Clean test cache
	@go clean -testcache ./...

.PHONY: test
test: ### Run tests
	@go test -v -coverprofile=cover.out -timeout 10s ./...

.PHONY: cover
cover: test ### Run tests and generate coverage
	@go tool cover -html=cover.out -o=cover.html

.PHONY: mocks
mocks: mockery ### Generate mocks
	$(MOCKERY) --all --dir internal --output internal/mocks

# go-get-tool will 'go get' any package $2 and install it to $1.
PROJECT_DIR := $(shell dirname $(abspath $(lastword $(MAKEFILE_LIST))))
define go-install
@[ -f $(1) ] || { \
go install $(1) ; \
}
endef
