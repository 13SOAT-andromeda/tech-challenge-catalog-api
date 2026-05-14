.PHONY: build run test lint clean docker-up docker-down

BINARY_NAME=server

build:
	go build -o $(BINARY_NAME) ./cmd/api

run:
	go run ./cmd/api/main.go

test:
	go test -v ./...

lint:
	golangci-lint run

clean:
	rm -f $(BINARY_NAME)
	rm -rf tmp

docker-up:
	docker-compose up --build

docker-down:
	docker-compose down
