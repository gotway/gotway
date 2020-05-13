.PHONY: clean deps format build run
all: build
clean:
	rm -rf bin/
deps:
	dep ensure -v
format:
	go fmt .
build: clean deps format
	go build -o bin/microgateway -v .
install: deps format
	go install -v .
run: build
	./bin/microgateway
release:
	./scripts/release.sh
deploy:
	./scripts/deploy.sh