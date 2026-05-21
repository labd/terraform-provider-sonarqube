.PHONY: build test testrace fmt vet lint

default: build

build:
	go build

test:
	go test -v ./... -timeout=120s -parallel=4

fmt:
	gofmt -s -w .

vet:
	go vet ./...
