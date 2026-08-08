.PHONY: test test-race vet build fmt run

test:
	go test ./...

test-race:
	go test -race ./...

vet:
	go vet ./...

build:
	go build ./...

fmt:
	gofmt -w .

run:
	go run ./cmd/api
