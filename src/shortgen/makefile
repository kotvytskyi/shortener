-include ./cmd/shortgen/.env
export 

test:
	go test ./...

lint:
	golangci-lint run

ci:
	make test

dev:
	make test
	make lint