.PHONY: clean deps deps-sync format build install run

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
	go build -o bin/stock -v .
install:
	go install -v .
run: build
	./bin/stock