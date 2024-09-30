run:
	go run .
build:
	go build -ldflags="-s -w" -o bin/ .

test:
	go test -v ./...