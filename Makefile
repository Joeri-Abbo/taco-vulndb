.PHONY: build test lint clean docker docker-vulndb

build:
	go build -o bin/taco-vulndb ./cmd/taco-vulndb

test:
	go test -race -coverprofile=coverage.out ./...

lint:
	golangci-lint run ./...

vet:
	go vet ./...

clean:
	rm -rf bin/ coverage.out

docker:
	docker build -f taco-vulndb/Dockerfile -t taco-vulndb:latest ..

docker-vulndb:
	docker build -f taco-vulndb/Dockerfile.vulndb -t taco-vulndb-image:latest ..
