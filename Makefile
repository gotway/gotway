.PHONY: clean deps format build install run release deploy-pre deploy-pro

all: build
clean:
	rm -rf bin/
deps:
	go get -u -v
format:
	go fmt .
build: clean deps format
	go build -o bin/microgateway -v .
install: deps
	go install -v .
run: build
	./bin/microgateway
release:
	./scripts/release.sh
deploy-pre:
	./scripts/deploy.sh "preproduction"
deploy-pro:
	./scripts/deploy.sh "production"