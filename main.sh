#!/usr/bin/env sh

# If you need a special setup to run your tests, modify this file.
# Alternatively you can use the terminal window below, to run any commands, like:
# go mod download
# go test ./...
# go run cmd/store/.go

echo "==> Running go mod download && go test -v ./..."
go mod download && go test -v ./...
