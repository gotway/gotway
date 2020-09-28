.PHONY: clean deps deps-sync format build install run test cover

all: build
clean:
	rm -rf bin/
deps:
	go get -u -v
deps-sync:
	go mod vendor
format:
	go fmt .
build: clean format
	go build -o bin/microgateway -v .
install:
	go install -v .
run: build
	./bin/microgateway
test:
	go test -v ./... -coverprofile=cover.out
cover: test
	go tool cover -html=cover.out