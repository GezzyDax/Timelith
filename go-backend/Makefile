.PHONY: build run test clean install

build:
	go build -o bin/server ./cmd/server

run:
	go run ./cmd/server/main.go

install:
	go mod download
	go mod tidy

test:
	go test -v ./...

clean:
	rm -rf bin/
	go clean

docker-build:
	docker build -t timelith-backend .

docker-run:
	docker run -p 8080:8080 --env-file .env timelith-backend
