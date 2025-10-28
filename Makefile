.PHONY: build test lint docker clean

build:
	go build -o bin/rootca ./cmd/rootca

test:
	go test ./... -v

lint:
	golangci-lint run ./...

docker:
	docker build -t gigvault/rootca:local .

clean:
	rm -rf bin/ certs/
	go clean

