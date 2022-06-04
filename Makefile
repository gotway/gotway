ROOT := $(shell git rev-parse --show-toplevel)
OS := $(shell uname -s | awk '{print tolower($$0)}')
ARCH := amd64
PROJECT := gotway
VERSION := $(git describe --abbrev=0 --tags)
LD_FLAGS := -X main.version=$(VERSION) -s -w
SOURCE_FILES ?= ./internal/... ./pkg/... ./cmd/...

# go
export CGO_ENABLED := 0
export GO111MODULE := on

# gotway
export PORT ?= 11000
export ENV ?= local
export LOG_LEVEL ?= debug
export REDIS_URL ?= redis://localhost:6379/11
export KUBECONFIG ?= $(HOME)/.kube/config
export METRICS ?= true
export METRICS_PATH ?= /metrics
export METRICS_PORT ?= 2112
export PPROF ?= false
export PPROF_PORT ?= 6060

# bin
BIN := $(ROOT)/bin
KIND := $(BIN)/kind
KIND_VERSION := v0.12.0
MOCKERY := $(BIN)/mockery
MOCKERY_VERSION := 2.12.3

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

$(BIN):
	@mkdir -p $(BIN)

$(KIND): $(BIN)
	@curl -Lo $(KIND) https://kind.sigs.k8s.io/dl/$(KIND_VERSION)/kind-$(OS)-$(ARCH)
	@chmod +x $(KIND)

$(MOCKERY): $(BIN)
	@curl -Lo mockery.tar.gz https://github.com/vektra/mockery/releases/download/v$(MOCKERY_VERSION)/mockery_$(MOCKERY_VERSION)_$(shell uname -s)_$(shell uname -m).tar.gz
	@tar -C /tmp -zxvf mockery.tar.gz 
	@mv /tmp/mockery $(MOCKERY)
	@chmod +x $(MOCKERY)

.PHONY: generate
generate: vendor ### Generate code
	@bash ./hack/hack.sh

.PHONY: install
install: ## Install CRDs
	@kubectl apply -f manifests/crds

.PHONY: cluster
cluster: $(KIND) ### Create a KIND cluster
	$(KIND) create cluster --name $(PROJECT)

.PHONY: docker
docker: ### Spin up docker dependencies
	@docker compose up -d

clean: ### Clean build files
	@rm -rf ./bin
	@go clean

.PHONY: build
build: generate clean ### Build binary
	@go build -tags netgo -a -v -ldflags "${LD_FLAGS}" -o ./bin/gotway ./cmd/gotway/*.go
	@chmod +x ./bin/*

.PHONY: run
run: docker install ### Quick run
	@CGO_ENABLED=1 go run -race cmd/gotway/*.go

.PHONY: deps
deps: ### Optimize dependencies
	@go mod tidy

.PHONY: vendor
vendor: ### Vendor dependencies
	@go mod vendor

.PHONY: fmt
fmt: ### Format
	@gofmt -s -w .

.PHONY: vet
vet: ### Vet
	@go vet ./...

### Lint
.PHONY: lint
lint: fmt vet

### Clean test 
.PHONY: test-clean
test-clean: ### Clean test cache
	@go clean -testcache ./...

.PHONY: test
test: lint ### Run tests
	@go test -v  -coverprofile=cover.out -timeout 10s ./...

.PHONY: cover
cover: test ### Run tests and generate coverage
	@go tool cover -html=cover.out -o=cover.html

.PHONY: mocks
mocks: $(MOCKERY) ### Generate mocks
	$(MOCKERY) --all --dir internal --output internal/mocks