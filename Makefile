.PHONY: build run docker-build docker-up migrate lint fmt test


build:
	go build -o bin/pr-reviewer ./cmd/server


run: build
	./bin/pr-reviewer


docker-build:
	docker build -t pr-reviewer:local -f build/Dockerfile .


docker-up:
	docker-compose up --build


migrate:
	docker-compose run --rm migrate


lint:
	golangci-lint run


fmt:
	gofmt -w .


test:
	go test ./...