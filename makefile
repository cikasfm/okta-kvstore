clean:
	rm main >> /dev/null || true

setup:
	go mod download

build: setup
	go build cmd/store/*.go

run: clean setup build
	./main

test: setup
	GIN_MODE=test go test -v ./...
