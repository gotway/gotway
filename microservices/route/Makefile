.PHONY: clean deps deps-sync fmt-srv fmt-cli fmt proto build srv run cli cli.mgw

MGW_LOCAL_ARGS = -server localhost:8000 -tls true
MGW_PRE_ARGS = -server microgateway.duckdns.org:30000 -tls true

all: build
clean:
	rm -rf bin/
deps:
	go get -u -v
deps-sync:
	go mod vendor
fmt-srv: 
	go fmt server/*.go 
fmt-cli: 
	go fmt client/*.go
fmt: fmt-srv fmt-cli
proto:
	protoc -I pb/ --go_opt=paths=source_relative --go_out=plugins=grpc:pb pb/*.proto
build: fmt-srv proto clean
	go build -o bin/route -v .
srv: build
	./bin/route
run: srv
cli: fmt-cli
	go run client/*.go $(MGW_LOCAL_ARGS)
cli-pre: fmt-cli
	go run client/*.go $(MGW_PRE_ARGS)
