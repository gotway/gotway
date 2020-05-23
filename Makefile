.PHONY: clean deps format build install run release deploy-pre deploy-pro

all: build
clean:
	rm -rf bin/
deps:
	go get -u -v
sync-deps:
	go mod vendor
format:
	go fmt .
build: clean format
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