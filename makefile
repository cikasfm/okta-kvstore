clean:
	rm main >> /dev/null || true

setup:
	go mod download

build: setup
	go build cmd/store/*.go

run: clean setup build
	NODE_ID=app1 RAFT_BIND_ADDR=127.0.0.1:7000 ./main

test: setup
	GIN_MODE=test go test -v ./...

docker:
	docker compose up --build