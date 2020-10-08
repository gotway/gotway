.PHONY: clean deps deps-sync fmt vet lint build install run test cover mocks

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
build: lint clean
	go build -o bin/gotway -v .
install:
	go install -v .
run: build
	./bin/gotway
test: lint
	go test -v ./... -coverprofile=cover.out
cover: test
	go tool cover -html=cover.out
mocks:
	mockery --all --keeptree