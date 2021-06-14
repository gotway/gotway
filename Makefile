.PHONY: clean deps deps-sync fmt vet lint build install run clean-test test cover mocks

all: build
clean:
	rm -rf bin/
deps:
	go get -u -v
deps-sync:
	go mod vendor
fmt:
	gofmt -s -w .
vet:
	go vet ./...
lint: fmt vet
build: clean
	go build -o bin/gotway -v cmd/gotway/*.go
install:
	go install -v cmd/gotway/*.go
run:
	go run cmd/gotway/*.go
clean-test:
	go clean -testcache ./...
test: lint
	go test -v  -coverprofile=cover.out -timeout 10s ./...
cover: test
	go tool cover -html=cover.out -o=cover.html
mocks:
	mockery --all --keeptree