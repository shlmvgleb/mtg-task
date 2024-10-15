all: build

bin:
	mkdir -p bin

build: bin
	go build -o bin ./cmd/...

run: build
	go run cmd/app/main.go

test:
	go test -v ./...
